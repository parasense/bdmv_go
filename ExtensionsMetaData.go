package main

import (
	"encoding/binary"
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
		return nil, err
	}

	if err := binary.Read(file, binary.BigEndian, &metaData.EntryDataStartAddr); err != nil {
		return nil, err
	}

	// Skip 3-byte reserve space
	if _, err := file.Seek(3, io.SeekCurrent); err != nil {
		return nil, err
	}

	if err := binary.Read(file, binary.BigEndian, &metaData.EntryDataCount); err != nil {
		return nil, err
	}

	return metaData, err
}

func (metaData *ExtensionsMetaData) Print() {
	PadPrintf(4, "ExtensionsMetaData.Length: %d\n", metaData.Length)
	PadPrintf(4, "ExtensionsMetaData.EntryDataStartAddr: %d\n", metaData.EntryDataStartAddr)
	PadPrintf(4, "ExtensionsMetaData.EntryDataCount: %d\n", metaData.EntryDataCount)
}
