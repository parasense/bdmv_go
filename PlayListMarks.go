package main

import (
	"encoding/binary"
	"fmt"
	"io"
)

// TODO: Fix the timestamps to something more sensible.

type PlaylistMarks struct {
	Length        uint32
	NumberOfMarks uint16
	Marks         []*MarkEntry
}

// MarkEntry represents a chapter mark or other marker in the playlist
type MarkEntry struct {
	MarkType        uint8
	RefToPlayItemID uint16
	MarkTimeStamp   uint32 // in 45kHz ticks
	EntryESPID      uint16
	Duration        uint32 // in 45kHz ticks
}

func ReadMarks(file io.ReadSeeker, offsets *OffsetsUint32) (marks *PlaylistMarks, err error) {
	marks = &PlaylistMarks{}

	// Jump to start address
	if _, err := file.Seek(offsets.Start, io.SeekStart); err != nil {
		return nil, fmt.Errorf("failed to seek to start address: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &marks.Length); err != nil {
		return nil, fmt.Errorf("failed to read mark length: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &marks.NumberOfMarks); err != nil {
		return nil, fmt.Errorf("failed to read number of marks: %w", err)
	}

	marks.Marks = make([]*MarkEntry, marks.NumberOfMarks)
	for i := uint16(0); i < marks.NumberOfMarks; i++ {
		if marks.Marks[i], err = ReadMarkEntry(file); err != nil {
			return nil, fmt.Errorf("failed to ReadMarkEntry: %w", err)
		}
	}
	return marks, nil

}

func (playlistMarks *PlaylistMarks) Print() {
	PadPrintf(0, "Chapter Marks: [%d]\n", len(playlistMarks.Marks))
	for i, mark := range playlistMarks.Marks {
		if mark.MarkType == 1 { // Chapter mark
			timestamp := Parse45KhzTimestamp(mark.MarkTimeStamp)
			PadPrintf(2, "Chapter [%d]: at [%v] (PlayItem: %d)\n", i+1, timestamp, mark.RefToPlayItemID)
			mark.Print()
			PadPrintln(2, "---")
		}
	}
}

func ReadMarkEntry(file io.ReadSeeker) (markEntry *MarkEntry, err error) {
	markEntry = &MarkEntry{}

	// Skip 1-byte reserve space
	if _, err := file.Seek(1, io.SeekCurrent); err != nil {
		return nil, fmt.Errorf("failed to seek past reserve space: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &markEntry.MarkType); err != nil {
		return nil, fmt.Errorf("failed to read markEntry type: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &markEntry.RefToPlayItemID); err != nil {
		return nil, fmt.Errorf("failed to read markEntry play item ref ID: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &markEntry.MarkTimeStamp); err != nil {
		return nil, fmt.Errorf("failed to read markEntry timestamp: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &markEntry.EntryESPID); err != nil {
		return nil, fmt.Errorf("failed to read markEntry ES PID: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &markEntry.Duration); err != nil {
		return nil, fmt.Errorf("failed to read markEntry duration: %w", err)
	}

	return markEntry, nil
}

func (markEntry *MarkEntry) Print() {
	PadPrintf(4, "MarkType: %d\n", markEntry.MarkType)
	PadPrintf(4, "RefToPlayItemID: %d\n", markEntry.RefToPlayItemID)
	PadPrintf(4, "MarkTimeStamp: %d\n", markEntry.MarkTimeStamp)
	PadPrintf(4, "EntryESPID: %d\n", markEntry.EntryESPID)
	PadPrintf(4, "Duration: %d\n", markEntry.Duration)
}
