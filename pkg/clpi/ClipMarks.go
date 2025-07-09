package clpi

import (
	"encoding/binary"
	"fmt"
	"io"
)

type ClipMarks struct {
	Length            uint32
	NumberOfClipMarks uint16
	MarkEntries       []*ClipMarkEntry
}

type ClipMarkEntry struct {
	_              uint8  // 8-bit reserved (usually 0x00)
	MarkType       uint8  // 8-bit unsigned integer
	MarkPID        uint16 // 16-bit unsigned integer
	MarkTimeStamp  uint32 // 32-bit unsigned integer
	MarkEntryPoint uint32 // 32-bit unsigned integer
	MarkDuration   uint32 // 32-bit unsigned integer
}

func ReadClipMarks(file io.ReadSeeker, offsets *OffsetsUint32) (clipMarks *ClipMarks, err error) {
	// Avoid allocating the struct instance if the offsets are zero.
	if offsets.Start == 0 && offsets.Stop == 0 {
		return nil, nil
	}

	// XXX
	//fmt.Printf("ClipMarks offsets: %+v\n", offsets)

	clipMarks = &ClipMarks{}

	// XXX
	//log.Printf("ReadClipMarks\n")

	// Jump to start address
	if _, err := file.Seek(offsets.Start, io.SeekStart); err != nil {
		return nil, err
	}

	if err := binary.Read(file, binary.BigEndian, &clipMarks.Length); err != nil {
		return nil, err
	}
	fmt.Printf("clipMarks.Length: %d\n", clipMarks.Length)

	// XXX - IMPORTANT
	// Testing on real CLPI files has show that sometimes the length is zero!
	// That means parsing might need to stop here.
	if clipMarks.Length == 0 {
		return clipMarks, nil
	}

	if err := binary.Read(file, binary.BigEndian, &clipMarks.NumberOfClipMarks); err != nil {
		return nil, err
	}
	fmt.Printf("clipMarks.NumberOfClipMarks: %d\n", clipMarks.NumberOfClipMarks)

	// XXX - This is all very expeerimental bellow
	clipMarks.MarkEntries = make([]*ClipMarkEntry, clipMarks.NumberOfClipMarks)
	for i := range clipMarks.MarkEntries {
		clipMarks.MarkEntries[i] = &ClipMarkEntry{}
		fmt.Printf("DEBUG: [%d] MarkEntry: %+v\n", i, clipMarks.MarkEntries[i])
	}

	return clipMarks, nil
}

func (clipMarks *ClipMarks) String() string {

	return fmt.Sprintf(
		"ClipMarks{"+
			"Length: %d, "+
			"NumberOfClipMarks: %d, "+
			"MarkEntries: %+v, "+
			"}",
		clipMarks.Length,
		clipMarks.NumberOfClipMarks,
		clipMarks.MarkEntries,
	)
}

func (entry *ClipMarkEntry) String() string {
	return fmt.Sprintf(
		"ClipMarkEntry{MarkType: %d, MarkPID: %d, MarkTimeStamp: %d, MarkEntryPoint: %d, MarkDuration: %d}",
		entry.MarkType,
		entry.MarkPID,
		entry.MarkTimeStamp,
		entry.MarkEntryPoint,
		entry.MarkDuration,
	)
}
