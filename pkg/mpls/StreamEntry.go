package mpls

import (
	"encoding/binary"
	"fmt"
	"io"
)

// BasicStreamEntry is the base structure for all StreamEntry types.
type BasicStreamEntry struct {
	Length     uint8
	StreamType uint8
}

// SetLength sets the Length for the BasicStreamEntry.
func (streamEntry *BasicStreamEntry) SetLength(length uint8) {
	streamEntry.Length = length
}

// SetStreamType sets the StreamType for the BasicStreamEntry.
func (streamEntry *BasicStreamEntry) SetStreamType(streamType uint8) {
	streamEntry.StreamType = streamType
}

// StreamEntryTypeI is used for StreamType 1.
type StreamEntryTypeI struct {
	BasicStreamEntry
	RefToStreamPID uint16
}

// StreamEntryTypeII is used for StreamType 2.
type StreamEntryTypeII struct {
	BasicStreamEntry
	RefToSubPathID uint8
	RefToSubClipID uint8
	RefToStreamPID uint16
}

// StreamEntryTypeIII is used for StreamType 3 and 4.
type StreamEntryTypeIII struct {
	BasicStreamEntry
	RefToSubPathID uint8
	RefToStreamPID uint16
}

// StreamEntry is an interface that defines the methods for reading and setting
type StreamEntry interface {
	Read(io.ReadSeeker) error
	SetLength(uint8)
	SetStreamType(uint8)
}

// ReadStreamEntry reads a StreamEntry from the provided io.ReadSeeker.
// It expects the file (io.ReadSeeker) to be positioned at the start of the
// StreamEntry structure.
// It determines the type of StreamEntry based on the StreamType byte and
// reads the corresponding structure accordingly.
// It returns a StreamEntry interface and an error if any occurs during reading.
func ReadStreamEntry(file io.ReadSeeker) (entry StreamEntry, err error) {
	var buffer byte

	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return nil, fmt.Errorf("Failed to read Stream Entry Length: %w", err)
	}
	var length uint8 = buffer

	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return nil, fmt.Errorf("Failed to read Stream Entry StreamType: %w", err)
	}
	var streamType uint8 = buffer

	switch streamType {
	case 1:
		entry = &StreamEntryTypeI{}
	case 2:
		entry = &StreamEntryTypeII{}
	case 3, 4:
		entry = &StreamEntryTypeIII{}
	default:
		return nil, fmt.Errorf("ReadStreamEntry(): Unknown Stream Entry type: %d", streamType)
	}

	if entry != nil {
		entry.SetLength(length)
		entry.SetStreamType(streamType)
		if err := entry.Read(file); err != nil {
			return nil, fmt.Errorf("Failed calling entry.Read() on StreamEntry (type %d): %w", streamType, err)
		}
	} else {
		// xxx
		fmt.Printf("ERROR: unkown StreamEntry is nil\n")
	}

	return entry, nil
}

// Read reads the StreamEntryTypeI structure from the provided io.ReadSeeker.
func (entry *StreamEntryTypeI) Read(file io.ReadSeeker) (err error) {

	if err := binary.Read(file, binary.BigEndian, &entry.RefToStreamPID); err != nil {
		return fmt.Errorf("failed to read Stream Entry RefToStreamPID: %w", err)
	}

	// 6 tail padding bytes
	if _, err := file.Seek(6, io.SeekCurrent); err != nil {
		return fmt.Errorf("failed to seek past reserve space: %w", err)
	}

	return nil
}

// Read reads the StreamEntryTypeII structure from the provided io.ReadSeeker.
func (entry *StreamEntryTypeII) Read(file io.ReadSeeker) (err error) {

	if err := binary.Read(file, binary.BigEndian, &entry.RefToSubPathID); err != nil {
		return fmt.Errorf("failed to read stream Entry RefToSubPathID: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &entry.RefToSubClipID); err != nil {
		return fmt.Errorf("failed to read stream Entry RefToSubClipID: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &entry.RefToStreamPID); err != nil {
		return fmt.Errorf("failed to read stream Entry RefToStreamPID: %w", err)
	}

	// 4 tail padding bytes
	if _, err := file.Seek(4, io.SeekCurrent); err != nil {
		return fmt.Errorf("failed to seek past reserve space: %w", err)
	}

	return nil
}

// Read reads the StreamEntryTypeIII structure from the provided io.ReadSeeker.
func (entry *StreamEntryTypeIII) Read(file io.ReadSeeker) (err error) {

	if err := binary.Read(file, binary.BigEndian, &entry.RefToSubPathID); err != nil {
		return fmt.Errorf("failed to read stream Entry RefToSubPathID: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &entry.RefToStreamPID); err != nil {
		return fmt.Errorf("failed to read stream Entry RefToStreamPID: %w", err)
	}

	// 5 tail padding bytes
	if _, err := file.Seek(5, io.SeekCurrent); err != nil {
		return fmt.Errorf("failed to seek past reserve space: %w", err)
	}

	return nil
}
