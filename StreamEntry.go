package main

import (
	"encoding/binary"
	"fmt"
	"io"
)

// These seem to be 10 byte fixed length structures.
// Depending on the StreamType different tail padding sizes happen.

type BasicStreamEntry struct {
	Length     uint8
	StreamType uint8
}

type StreamEntryTypeI struct {
	BasicStreamEntry
	RefToStreamPID uint16
}

type StreamEntryTypeII struct {
	BasicStreamEntry
	RefToSubPathID uint8
	RefToSubClipID uint8
	RefToStreamPID uint16
}

type StreamEntryTypeIII struct {
	BasicStreamEntry
	RefToSubPathID uint8
	RefToStreamPID uint16
}

type StreamEntry interface {
	Print()
	Read(io.ReadSeeker) error
}

func ReadStreamEntry(file io.ReadSeeker) (entry StreamEntry, err error) {
	var buffer byte

	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return nil, fmt.Errorf("failed to read Stream Entry Length: %v\n", err)
	}
	var length uint8 = buffer

	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return nil, fmt.Errorf("failed to read Stream Entry StreamType: %v\n", err)
	}
	var streamType uint8 = buffer

	switch streamType {
	case 1:
		entry := &StreamEntryTypeI{}
		entry.Length = length
		entry.StreamType = streamType
		if err := entry.Read(file); err != nil {
			return nil, fmt.Errorf("failed calling entry.Read() on StreamEntry: %v\n", err)
		}
		return entry, nil

	case 2:
		entry := &StreamEntryTypeII{}
		entry.Length = length
		entry.StreamType = streamType
		if err := entry.Read(file); err != nil {
			return nil, fmt.Errorf("failed calling entry.Read() on StreamEntry: %v\n", err)
		}
		return entry, nil

	case 3, 4:
		entry := &StreamEntryTypeIII{}
		entry.Length = length
		entry.StreamType = streamType
		if err := entry.Read(file); err != nil {
			return nil, fmt.Errorf("failed calling entry.Read() on StreamEntry: %v\n", err)
		}
		return entry, nil

	default:
		return nil, fmt.Errorf("Unknown Stream Entry type\n")
	}

}

// Print implements NewPrimaryStreamEntry.
func (entry *StreamEntryTypeI) Print() {
	PadPrintln(6, "StreamEntry:")
	PadPrintf(8, "Length: %d\n", entry.Length)
	PadPrintf(8, "StreamType: %d\n", entry.StreamType)
	PadPrintf(8, "RefToStreamPID: %d\n", entry.RefToStreamPID)
}

func (entry *StreamEntryTypeII) Print() {
	PadPrintln(6, "StreamEntry:")
	PadPrintf(8, "Length: %d\n", entry.Length)
	PadPrintf(8, "StreamType: %d\n", entry.StreamType)
	PadPrintf(8, "RefToSubPathID: %d\n", entry.RefToSubPathID)
	PadPrintf(8, "RefToSubClipID: %d\n", entry.RefToSubClipID)
	PadPrintf(8, "RefToStreamPID: %d\n", entry.RefToStreamPID)
}

func (entry *StreamEntryTypeIII) Print() {
	PadPrintln(6, "StreamEntry:")
	PadPrintf(8, "Length: %d\n", entry.Length)
	PadPrintf(8, "StreamType: %d\n", entry.StreamType)
	PadPrintf(8, "RefToSubPathID: %d\n", entry.RefToSubPathID)
	PadPrintf(8, "RefToStreamPID: %d\n", entry.RefToStreamPID)
}

func (entry *StreamEntryTypeI) Read(file io.ReadSeeker) (err error) {

	if err := binary.Read(file, binary.BigEndian, &entry.RefToStreamPID); err != nil {
		return fmt.Errorf("failed to read Stream Entry RefToStreamPID: %v\n", err)
	}

	// 6 tail padding bytes
	if _, err := file.Seek(6, io.SeekCurrent); err != nil {
		return fmt.Errorf("failed to seek past reserve space: %v\n", err)
	}

	return nil
}

func (entry *StreamEntryTypeII) Read(file io.ReadSeeker) (err error) {

	if err := binary.Read(file, binary.BigEndian, &entry.RefToSubPathID); err != nil {
		return fmt.Errorf("failed to read stream Entry RefToSubPathID: %v\n", err)
	}

	if err := binary.Read(file, binary.BigEndian, &entry.RefToSubClipID); err != nil {
		return fmt.Errorf("failed to read stream Entry RefToSubClipID: %v\n", err)
	}

	if err := binary.Read(file, binary.BigEndian, &entry.RefToStreamPID); err != nil {
		return fmt.Errorf("failed to read stream Entry RefToStreamPID: %v\n", err)
	}

	// 4 tail padding bytes
	if _, err := file.Seek(4, io.SeekCurrent); err != nil {
		return fmt.Errorf("failed to seek past reserve space: %v\n", err)
	}

	return nil
}

func (entry *StreamEntryTypeIII) Read(file io.ReadSeeker) (err error) {

	if err := binary.Read(file, binary.BigEndian, &entry.RefToSubPathID); err != nil {
		return fmt.Errorf("failed to read stream Entry RefToSubPathID: %v\n", err)
	}

	if err := binary.Read(file, binary.BigEndian, &entry.RefToStreamPID); err != nil {
		return fmt.Errorf("failed to read stream Entry RefToStreamPID: %v\n", err)
	}

	// 5 tail padding bytes
	if _, err := file.Seek(5, io.SeekCurrent); err != nil {
		return fmt.Errorf("failed to seek past reserve space: %v\n", err)
	}

	return nil
}
