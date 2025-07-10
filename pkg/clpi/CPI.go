package clpi

import (
	"encoding/binary"
	"fmt"
	"io"
)

type CPI struct {
	Length                   uint32
	CPIType                  uint8 // 4-bits 0b00001111
	NumberOfStreamPIDEntries uint8
	StreamPIDEntries         []*StreamPIDEntry
}

type StreamPIDEntry struct {
	StreamPID               uint16
	EPStreamType            uint8  // 4-bits 0b11110000
	NumberOfEPCoarseEntries uint16 // 16-bits
	NumberOfEPFineEntries   uint32 // 18-bits
	EPMapStreamStartAddr    uint32
	EPFineTableStartAddress uint32
	CourseEntries           []*CourseEntry
	FineEntries             []*FineEntry
}

type CourseEntry struct {
	RefToEPFineID uint32 // 18-bits 0b11111111_11111111_11000000_00000000
	PTSEPCoarse   uint16 // 14-bits 0b00000000_00000000_00111111_11111111
	SPNEPCoarse   uint32 // 32-bits
}

type FineEntry struct {
	IsAngleChangePoint bool   //  1-bit  0b10000000_00000000_00000000_00000000
	IEndPositionOffset uint8  //  3-bits 0b01110000_00000000_00000000_00000000
	PTSEPFine          uint16 // 11-bits 0b00001111_11111110_00000000_00000000
	SPNEPFine          uint32 // 17-bits 0b00000000_00000001_11111111_11111111
}

func ReadCPI(file io.ReadSeeker, offsets *OffsetsUint32) (cpi *CPI, err error) {

	// Avoid allocating the struct instance if the offsets are zero.
	if offsets.Start == 0 && offsets.Stop == 0 {
		return nil, nil
	}

	cpi = &CPI{}

	// Jump to start address
	if _, err := file.Seek(offsets.Start, io.SeekStart); err != nil {
		return nil, err
	}

	if err := binary.Read(file, binary.BigEndian, &cpi.Length); err != nil {
		return nil, err
	}

	// Testing on real CLPI files has show that sometimes the length is zero!
	// That means parsing might need to stop here.
	if cpi.Length == 0 {
		return cpi, nil
	}

	// Reserve space 1-bytes.
	if _, err := file.Seek(1, io.SeekCurrent); err != nil {
		return nil, fmt.Errorf("failed to seek to beyond reserved space: %w", err)
	}

	var buffer uint8
	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return nil, fmt.Errorf("failed to read buffer: %w", err)
	}
	cpi.CPIType = (buffer & 0x0F)

	// Reserve space 1-bytes.
	if _, err := file.Seek(1, io.SeekCurrent); err != nil {
		return nil, fmt.Errorf("failed to seek to beyond reserved space: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &cpi.NumberOfStreamPIDEntries); err != nil {
		return nil, err
	}

	// Capture StreamPID entries metadata here
	cpi.StreamPIDEntries = make([]*StreamPIDEntry, cpi.NumberOfStreamPIDEntries)
	for i := range cpi.StreamPIDEntries {
		cpi.StreamPIDEntries[i], err = ReadStreamPIDEntry(file)
	}

	for i, streamPID := range cpi.StreamPIDEntries {

		// This is where the jump to the "EPMapStreamStartAddr" happens.
		EPMapForOneStreamPIDStartAddress := offsets.Start + 6 + int64(streamPID.EPMapStreamStartAddr)

		if _, err := file.Seek(EPMapForOneStreamPIDStartAddress, io.SeekStart); err != nil {
			return nil, err
		}

		// This is where the FineEntry StartAddr is parsed.
		if err := binary.Read(file, binary.BigEndian, &cpi.StreamPIDEntries[i].EPFineTableStartAddress); err != nil {
			return nil, err
		}

		cpi.StreamPIDEntries[i].CourseEntries = make([]*CourseEntry, streamPID.NumberOfEPCoarseEntries)
		for j := range streamPID.CourseEntries {
			cpi.StreamPIDEntries[i].CourseEntries[j], err = ReadStreamPIDCourseEntry(file)
			if err != nil {
				return nil, err
			}
		}

		// Jump to the start address of the fine entries.
		if _, err := file.Seek(
			EPMapForOneStreamPIDStartAddress+int64(cpi.StreamPIDEntries[i].EPFineTableStartAddress),
			io.SeekStart); err != nil {
			return nil, err
		}

		cpi.StreamPIDEntries[i].FineEntries = make([]*FineEntry, streamPID.NumberOfEPFineEntries)
		for j := range streamPID.FineEntries {
			cpi.StreamPIDEntries[i].FineEntries[j], err = ReadStreamPIDFineEntry(file)
			if err != nil {
				return nil, err
			}
		}
	}

	return cpi, nil
}

func ReadStreamPIDEntry(file io.ReadSeeker) (entry *StreamPIDEntry, err error) {
	entry = &StreamPIDEntry{}

	if err := binary.Read(file, binary.BigEndian, &entry.StreamPID); err != nil {
		return nil, err
	}

	// Reserve space 1-byte.
	if _, err := file.Seek(1, io.SeekCurrent); err != nil {
		return nil, fmt.Errorf("failed to seek to beyond reserved space: %w", err)
	}

	var buf32 uint32
	if err := binary.Read(file, binary.BigEndian, &buf32); err != nil {
		return nil, fmt.Errorf("failed to read buffer: %w", err)
	}

	// 0b00111100_00000000_00000000_00000000
	//     ^^^^
	entry.EPStreamType = uint8((buf32 & 0x3C000000) >> 26)

	// 0b00000011_11111111_11111100_00000000
	//         ^^ ^^^^^^^^ ^^^^^^
	entry.NumberOfEPCoarseEntries = uint16((buf32 & 0x03FFC000) >> 10)

	var buf8 uint8
	if err := binary.Read(file, binary.BigEndian, &buf8); err != nil {
		return nil, fmt.Errorf("failed to read buffer: %w", err)
	}
	// 0b00000000_00000000_00000011_11111111 + 0b11111111
	entry.NumberOfEPFineEntries = ((buf32 & 0x3FF) << 8) | uint32(buf8)

	if err := binary.Read(file, binary.BigEndian, &entry.EPMapStreamStartAddr); err != nil {
		return nil, fmt.Errorf("failed to read EPMapStreamStartAddr: %w", err)
	}

	return entry, err
}

func ReadStreamPIDCourseEntry(file io.ReadSeeker) (ce *CourseEntry, err error) {
	ce = &CourseEntry{}

	var courseBuf uint32
	if err := binary.Read(file, binary.BigEndian, &courseBuf); err != nil {
		return nil, err
	}

	// 0b11111111_11111111_11000000_00000000
	//   ^^^^^^^^ ^^^^^^^^ ^^
	ce.RefToEPFineID = (courseBuf & 0xFFFFC000) >> 14

	// 0b00000000_00000000_00111111_11111111
	//                       ^^^^^^ ^^^^^^^^
	ce.PTSEPCoarse = uint16(courseBuf & 0x00003FFF)

	// 32-bits
	if err := binary.Read(file, binary.BigEndian, &ce.SPNEPCoarse); err != nil {
		return nil, err
	}

	return ce, err
}

func ReadStreamPIDFineEntry(file io.ReadSeeker) (fe *FineEntry, err error) {
	fe = &FineEntry{}

	var fineBuf uint32
	if err := binary.Read(file, binary.BigEndian, &fineBuf); err != nil {
		return nil, err
	}

	// 0b10000000_00000000_00000000_00000000
	//   ^
	fe.IsAngleChangePoint = fineBuf&0x80000000 != 0

	// 0b01110000_00000000_00000000_00000000
	//    ^^^
	fe.IEndPositionOffset = uint8((fineBuf & 0x70000000) >> 28)

	// 0b00001111_11111110_00000000_00000000
	//       ^^^^ ^^^^^^^
	fe.PTSEPFine = uint16((fineBuf & 0x0FFE0000) >> 17)

	// 0b00000000_00000001_11111111_11111111
	//                   ^ ^^^^^^^^ ^^^^^^^^
	fe.SPNEPFine = (fineBuf & 0x0001FFFF)

	return fe, err
}

func (cpi *CPI) String() string {
	return fmt.Sprintf(
		"CPI{"+
			"Length: %d, "+
			"CPIType: %d, "+
			"NumberOfStreamPIDEntries: %d, "+
			"StreamPIDEntries: %+v, "+
			"}",
		cpi.Length,
		cpi.CPIType,
		cpi.NumberOfStreamPIDEntries,
		cpi.StreamPIDEntries,
	)
}

func (entry *StreamPIDEntry) String() string {
	return fmt.Sprintf(
		"StreamPIDEntry{"+
			"StreamPID: %d, "+
			"EPStreamType: %d, "+
			"NumberOfEPCoarseEntries: %d, "+
			"NumberOfEPFineEntries: %d, "+
			"EPMapStreamStartAddr: %d, "+
			"EPFineTableStartAddress: %d, "+
			"CourseEntries: %+v, "+
			"FineEntries: %+v, "+
			"}",
		entry.StreamPID,
		entry.EPStreamType,
		entry.NumberOfEPCoarseEntries,
		entry.NumberOfEPFineEntries,
		entry.EPMapStreamStartAddr,
		entry.EPFineTableStartAddress,
		entry.CourseEntries,
		entry.FineEntries,
	)
}

func (ce *CourseEntry) String() string {
	return fmt.Sprintf(
		"CourseEntry{"+
			"RefToEPFineID: %d, "+
			"PTSEPCoarse: %d, "+
			"SPNEPCoarse: %d, "+
			"}",
		ce.RefToEPFineID,
		ce.PTSEPCoarse,
		ce.SPNEPCoarse,
	)
}

func (fe *FineEntry) String() string {
	return fmt.Sprintf(
		"FineEntry{"+
			"IsAngleChangePoint: %t, "+
			"IEndPositionOffset: %d, "+
			"PTSEPFine: %d, "+
			"SPNEPFine: %d, "+
			"}",
		fe.IsAngleChangePoint,
		fe.IEndPositionOffset,
		fe.PTSEPFine,
		fe.SPNEPFine,
	)
}
