package main

import (
	"encoding/binary"
	"fmt"
	"io"
)

type ExtensionEntryMetaData struct {
	ExtDataType         uint16 // 2-bytes
	ExtDataVersion      uint16 // 2-bytes
	ExtDataStartAddress uint32 // 4-bytes
	ExtDataLength       uint32 // 4-bytes
}

func ReadEntryMetaData(file io.ReadSeeker) (entryMetaData *ExtensionEntryMetaData, err error) {
	entryMetaData = &ExtensionEntryMetaData{}

	if err := binary.Read(file, binary.BigEndian, &entryMetaData.ExtDataType); err != nil {
		return nil, fmt.Errorf("failed to read extDataEntry.ExtDataType: %v\n", err)
	}

	if err := binary.Read(file, binary.BigEndian, &entryMetaData.ExtDataVersion); err != nil {
		return nil, fmt.Errorf("failed to read extDataEntry.ExtDataVersion: %v\n", err)
	}

	if err := binary.Read(file, binary.BigEndian, &entryMetaData.ExtDataStartAddress); err != nil {
		return nil, fmt.Errorf("failed to read extDataEntry.ExtDataStartAddress: %v\n", err)
	}

	if err := binary.Read(file, binary.BigEndian, &entryMetaData.ExtDataLength); err != nil {
		return nil, fmt.Errorf("failed to read extDataEntry.ExtDataLength: %v\n", err)
	}

	return entryMetaData, err
}

func (entryMetaData *ExtensionEntryMetaData) Print() {
	PadPrintf(4, "ExtensionEntryMetaData.ExtDataType: %d\n", entryMetaData.ExtDataType)
	PadPrintf(4, "ExtensionEntryMetaData.ExtDataVersion: %d\n", entryMetaData.ExtDataVersion)
	PadPrintf(4, "ExtensionEntryMetaData.ExtDataStartAddress: %d\n", entryMetaData.ExtDataStartAddress)
	PadPrintf(4, "ExtensionEntryMetaData.ExtDataLength: %d\n", entryMetaData.ExtDataLength)
}
