package main

import (
	"encoding/binary"
	"fmt"
	"io"
)

type ExtensionsMetaData struct {
	Length             uint32
	EntryDataStartAddr uint32
	EntryDataCount     uint8
}

func ReadMetaData(file io.ReadSeeker) (metaData *ExtensionsMetaData, err error) {
	metaData = &ExtensionsMetaData{}

	if err := binary.Read(file, binary.BigEndian, &metaData.Length); err != nil {
		return nil, fmt.Errorf("failed to read metaData.Length: %v\n", err)
	}

	if err := binary.Read(file, binary.BigEndian, &metaData.EntryDataStartAddr); err != nil {
		return nil, fmt.Errorf("failed to read metaData.EntryDataStartAddr: %v\n", err)
	}

	// Skip 3-byte reserve space
	if _, err := file.Seek(3, io.SeekCurrent); err != nil {
		return nil, fmt.Errorf("failed to seek past reserve space: %v\n", err)
	}

	if err := binary.Read(file, binary.BigEndian, &metaData.EntryDataCount); err != nil {
		return nil, fmt.Errorf("failed to read metaData.EntryDataCount: %v\n", err)
	}

	return metaData, err
}

func (metaData *ExtensionsMetaData) Print() {
	PadPrintf(4, "ExtensionsMetaData.Length: %d\n", metaData.Length)
	PadPrintf(4, "ExtensionsMetaData.EntryDataStartAddr: %d\n", metaData.EntryDataStartAddr)
	PadPrintf(4, "ExtensionsMetaData.EntryDataCount: %d\n", metaData.EntryDataCount)
}
