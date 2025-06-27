package mpls

import (
	"encoding/binary"
	"fmt"
	"io"
)

// TODO: Fix the timestamps to something more sensible.

// PlaylistMarks represents a collection of chapter marks or other markers in a playlist.
// It contains the total length of the marks section, the number of marks,
// and a slice of MarkEntry pointers representing each individual mark.
// The marks are typically used to indicate chapter points or other significant timestamps
// within the playlist.
type PlaylistMarks struct {
	Length        uint32
	NumberOfMarks uint16
	Marks         []*MarkEntry
}

// MarkEntry represents a chapter mark or other marker in the playlist
// with its type, reference to a play item ID, timestamp, entry ES PID, and duration.
// The MarkType indicates the type of mark (e.g., chapter, cue point),
// RefToPlayItemID is a reference to the play item associated with the mark,
// MarkTimeStamp is the timestamp of the mark in 45kHz ticks,
// EntryESPID is the elementary stream PID associated with the mark,
// and Duration is the duration of the mark in 45kHz ticks.
// The MarkType is a single byte that indicates the type of mark, such as chapter or cue point.
// The RefToPlayItemID is a 2-byte reference to the play item ID
// that the mark is associated with, allowing for navigation to that play item.
type MarkEntry struct {
	MarkType        uint8
	RefToPlayItemID uint16
	MarkTimeStamp   uint32 // in 45kHz ticks
	EntryESPID      uint16
	Duration        uint32 // in 45kHz ticks
}

// ReadMarks reads the PlaylistMarks from the provided file at the specified offsets.
// It expects the file (io.ReadSeeker) to be positioned at the start of the
// PlaylistMarks structure, which consists of a length, number of marks, and the marks themselves.
// It returns a pointer to the PlaylistMarks and an error if any occurs during reading.
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

// ReadMarkEntry reads a single MarkEntry from the provided file.
// It expects the file (io.ReadSeeker) to be positioned at the start of the
// MarkEntry structure, which consists of a 1-byte reserved space, followed by
// the mark type (1 byte), reference to play item ID (2 bytes), timestamp (4 bytes),
// entry ES PID (2 bytes), and duration (4 bytes).
// It returns a pointer to the MarkEntry and an error if any occurs during reading.
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
