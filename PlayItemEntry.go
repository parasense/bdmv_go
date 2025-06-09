package main

import (
	"encoding/binary"
	"fmt"
	"io"
)

type PlayItemEntry struct {
	FileName   [5]byte
	Codec      [4]byte
	RefToSTCID uint8
}

func ReadPlayItemEntry(file io.ReadSeeker) (*PlayItemEntry, error) {
	playItemEntry := &PlayItemEntry{}

	// The 5 bytes clip name
	if err := binary.Read(file, binary.BigEndian, &playItemEntry.FileName); err != nil {
		return nil, fmt.Errorf("failed to read clip info filename: %w", err)
	}

	// The 4 byte codec should be something like "M2TS"
	if err := binary.Read(file, binary.BigEndian, &playItemEntry.Codec); err != nil {
		return nil, fmt.Errorf("failed to read clip codec identifier: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &playItemEntry.RefToSTCID); err != nil {
		return nil, fmt.Errorf("failed to read play item length: %w", err)
	}

	return playItemEntry, nil
}

func (playItemEntry *PlayItemEntry) Print() {
	PadPrintf(8, "FileName: %s\n", playItemEntry.FileName)
	PadPrintf(8, "Codec: %s\n", playItemEntry.Codec)
	PadPrintf(8, "RefToSTCID: %d\n", playItemEntry.RefToSTCID)
}
