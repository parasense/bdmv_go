package mpls

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Testing:
// Valerian 3D is used for validation.

// XXX - "Alice In Wonderland 3d" has broken MVC extension data
//
//	ALICE_IN_WONDERLAND_3D/BDMV/PLAYLIST/00002.mpls
//	This file has a uint16 length field set to 65288 (0xff00).
//	The length is factually correct, and yet the whole length is all zeros.
//
//	The same file has tail padding filled with "255 0 255 0 255 0 255 0 ..."
//

// ExtensionMVCStream implements the ExtensionEntryData interface.
type ExtensionMVCStream struct {
	MVCStreams []*MVCStream
}

// MVCStream represents a single MVC (Multiview Video Coding) stream entry in an MPLS file.
type MVCStream struct {
	Length                  uint16
	FixedOffsetPopUpFlag    bool // 0b10000000
	Entry                   StreamEntry
	Attr                    StreamAttributes
	NumberOfOffsetSequences uint8
	//remainder               []byte
}

// ReadMVC reads a single MVCStream entry from the provided io.ReadSeeker.
// It expects the file (io.ReadSeeker) to be positioned at the start of the
// MVCStream structure.
func (mvcStream *MVCStream) ReadMVC(file io.ReadSeeker, length uint16) (err error) {

	pos, _ := ftell(file)
	PadPrintf(2, "ReadMVC pos: %d\n", pos)
	mvcStream.Length = length

	var buffer byte
	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return fmt.Errorf("failed to read MVCStream.buffer: %w", err)
	}
	mvcStream.FixedOffsetPopUpFlag = buffer&0x80 != 0 // 0b10000000
	PadPrintf(4, "MVCStream.FixedOffsetPopUpFlag: %+v\n", mvcStream.FixedOffsetPopUpFlag)

	// 1-byte reserve space
	if _, err := file.Seek(1, io.SeekCurrent); err != nil {
		return fmt.Errorf("failed to seek past reserve space: %w", err)
	}

	mvcStream.Entry, err = ReadStreamEntry(file)
	if err != nil {
		return fmt.Errorf("failed to read StreamEntry: %w", err)
	}
	PadPrintf(4, "MVCStream.Entry: %+v\n", mvcStream.Entry)

	mvcStream.Attr, err = ReadStreamAttributes(file, STREAM_TYPE_PRIMARY_VIDEO)
	if err != nil {
		return fmt.Errorf("failed to read StreamAttributes: %w", err)
	}
	PadPrintf(4, "MVCStream.Attr: %+v\n", mvcStream.Attr)

	// 1-byte reserve space
	if _, err := file.Seek(1, io.SeekCurrent); err != nil {
		return fmt.Errorf("failed to seek past reserve space: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &mvcStream.NumberOfOffsetSequences); err != nil {
		return fmt.Errorf("failed to read MVCStream.NumberOfOffsetSequences: %w", err)
	}
	PadPrintf(4, "MVCStream.NumberOfOffsetSequences: %d\n", mvcStream.NumberOfOffsetSequences)
	return nil
}

func (extensionMVCStream *ExtensionMVCStream) Read(file io.ReadSeeker, offsets *OffsetsUint32, entryMeta *ExtensionEntryMetaData) (err error) {

	PadPrintln(0, "MVC Extension:")
	PadPrintln(2, "---")

	// Calculate the Start/Stop offsets for this extension.
	offsetStart := offsets.Start + int64(entryMeta.ExtDataStartAddress)
	offsetStop := offsetStart + int64(entryMeta.ExtDataLength)
	PadPrintf(2, "offsetStart == %d\n", offsetStart)
	PadPrintf(2, "offsetStop  == %d\n", offsetStop)
	PadPrintln(2, "---")

	// Jump to the start offset
	if _, err := file.Seek(int64(entryMeta.ExtDataStartAddress), io.SeekStart); err != nil {
		return fmt.Errorf("failed to seek Entry Offset Start: (%d); error: %w", entryMeta.ExtDataStartAddress, err)
	}

	// XXX - DEBUG block
	PadPrintln(0, "Extensions Entry DEBUG:")
	PadPrintf(2, "ExtDataType == %d\n", entryMeta.ExtDataType)
	PadPrintf(2, "ExtDataVersion == %d\n", entryMeta.ExtDataVersion)
	PadPrintf(2, "ExtDataStartAddress == %d\n", entryMeta.ExtDataStartAddress)
	PadPrintf(2, "ExtDataLength == %d\n", entryMeta.ExtDataLength)
	fmt.Println("---")
	// XXX - EO DEBUG block

	var loopIterLength uint16
	var loopIterEnd int64
	for i, loopIterStart := 1, offsetStart; loopIterEnd < offsetStop; i, loopIterStart = i+1, loopIterEnd {

		PadPrintf(4, "[%d] loopIterStart: %d\n", i, loopIterStart)

		// Before reading the length field...
		// Calculate if the uint16 (2-bytes) would go out of bounds.
		if loopIterStart+2 > offsetStop {
			PadPrintf(4, "[%d] Not enough space to read Length var uint16. break!\n", i)
			PadPrintf(4, "[%d] %d+2 > %d\n", i, loopIterStart, offsetStop)

			break
		}

		// Then go ahead to take the length uint16
		if err := binary.Read(file, binary.BigEndian, &loopIterLength); err != nil {
			return fmt.Errorf("[%d] failed to read MVCStream.Length: %w", i, err)
		}
		PadPrintf(4, "[%d] loopIterLength: %+v\n", i, loopIterLength)

		loopIterEnd = loopIterStart + 2 + int64(loopIterLength)
		PadPrintf(4, "[%d] loopIterEnd: %d\n", i, loopIterEnd)

		// Before initializing an instance of MVCStream, run sanity checks on the length.
		if loopIterLength == 0 {
			break

		} else if loopIterLength == 0xFF00 {
			continue // xxx - reject that case, for now...

		} else if loopIterEnd > offsetStop {
			// Next, if the length would go out of bounds, then bail
			PadPrintf(4, "[%d] Not enough space to read Length of segment. break!\n", i)
			PadPrintf(4, "[%d] %d+2 > %d\n", i, loopIterEnd, offsetStop)
			break
		}

		// MVCStream item
		MVCStream := &MVCStream{}

		// This reads exactly 18-bytes + 2-bytes (length)
		if err := MVCStream.ReadMVC(file, loopIterLength); err != nil {

			// XXX - the error would be an IO error.
			// Very unlikely give all thge sanity checks on boundaries.
			continue
		}

		// Append the structure to the end.
		extensionMVCStream.MVCStreams = append(extensionMVCStream.MVCStreams, MVCStream)

		// Calculate the remaining length
		remainderPos, _ := ftell(file)
		remainderLen := loopIterEnd - remainderPos
		PadPrintf(4, "[%d] remainderPos: %d\n", i, remainderPos)
		PadPrintf(4, "[%d] remainderLen: %d\n", i, remainderLen)

		var remainder []byte
		remainder = make([]byte, remainderLen)
		if err := binary.Read(file, binary.BigEndian, &remainder); err != nil {
			return fmt.Errorf("failed to read MVCStream.remainder: %w", err)
		}
		PadPrintf(4, "[%d] MVCStream.remainder: %+v\n", i, remainder)

		// Seek to the loop iteration end offset
		if _, err := file.Seek(loopIterEnd, io.SeekStart); err != nil {
			return fmt.Errorf("failed to seek past reserve space: %w", err)
		}

		PadPrintln(4, "")
		PadPrintln(4, "---")
	}

	// Jump to the extension entry stop offset
	if _, err := file.Seek(offsetStop, io.SeekStart); err != nil {
		return fmt.Errorf("failed to seek Entry Offset Stop: (%d); error: %w", offsetStart, err)
	}
	return nil
}
