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

func (streamEntry *BasicStreamEntry) SetLength(length uint8) {
	streamEntry.Length = length
}
func (streamEntry *BasicStreamEntry) SetStreamType(streamType uint8) {
	streamEntry.StreamType = streamType
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
	SetLength(uint8)
	SetStreamType(uint8)
}

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

	//// XXX
	//fmt.Printf("DEBUG: StreamEntry.Length: %d\n", length)
	//fmt.Printf("DEBUG: StreamEntry.streamType: %d\n", streamType)

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
		fmt.Printf("ERROR: unkown StreamEntry is nil\n")
	}

	return entry, nil
}

// Print implements NewPrimaryStreamEntry.
func (entry *StreamEntryTypeI) Print() {
	PadPrintln(10, "Entry:")
	PadPrintf(12, "Length: %d\n", entry.Length)
	PadPrintf(12, "StreamType: %d [%s]\n", entry.StreamType, StreamType(entry.StreamType))
	PadPrintf(12, "RefToStreamPID: %d\n", entry.RefToStreamPID)
}

func (entry *StreamEntryTypeII) Print() {
	PadPrintln(10, "Entry:")
	PadPrintf(12, "Length: %d\n", entry.Length)
	PadPrintf(12, "StreamType: %d [%s]\n", entry.StreamType, StreamType(entry.StreamType))
	PadPrintf(12, "RefToSubPathID: %d\n", entry.RefToSubPathID)
	PadPrintf(12, "RefToSubClipID: %d\n", entry.RefToSubClipID)
	PadPrintf(12, "RefToStreamPID: %d\n", entry.RefToStreamPID)
}

func (entry *StreamEntryTypeIII) Print() {
	PadPrintln(10, "Entry:")
	PadPrintf(12, "Length: %d\n", entry.Length)
	PadPrintf(12, "StreamType: %d [%s]\n", entry.StreamType, StreamType(entry.StreamType))
	PadPrintf(12, "RefToSubPathID: %d\n", entry.RefToSubPathID)
	PadPrintf(12, "RefToStreamPID: %d\n", entry.RefToStreamPID)
}

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
