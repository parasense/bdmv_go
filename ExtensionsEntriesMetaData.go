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

func ReadExtensionsEntriesMetaData(file io.ReadSeeker, offsets *OffsetsUint32, metaData *ExtensionsMetaData) (entriesMetaData []*ExtensionEntryMetaData, err error) {

	// Sanity check
	if int64(metaData.EntryDataCount)*12+12+offsets.Start > offsets.Stop {
		return nil, fmt.Errorf("ERROR: ReadEntryMetaData: not enough file to read Entires metadata.")
	}

	entriesMetaData = make([]*ExtensionEntryMetaData, metaData.EntryDataCount)

	for i := range entriesMetaData {
		entriesMetaData[i] = &ExtensionEntryMetaData{}

		if err := binary.Read(file, binary.BigEndian, &entriesMetaData[i].ExtDataType); err != nil {
			return nil, fmt.Errorf("failed to read ExtensionEntryMetaData.ExtDataType: %w", err)
		}

		if err := binary.Read(file, binary.BigEndian, &entriesMetaData[i].ExtDataVersion); err != nil {
			return nil, fmt.Errorf("failed to read ExtensionEntryMetaData.ExtDataVersion: %w", err)
		}

		if err := binary.Read(file, binary.BigEndian, &entriesMetaData[i].ExtDataStartAddress); err != nil {
			return nil, fmt.Errorf("failed to read ExtensionEntryMetaData.ExtDataStartAddress: %w", err)
		}

		if err := binary.Read(file, binary.BigEndian, &entriesMetaData[i].ExtDataLength); err != nil {
			return nil, fmt.Errorf("failed to read ExtensionEntryMetaData.EntryDataCount: %w", err)
		}
	}

	return entriesMetaData, err
}

func (entryMetaData *ExtensionEntryMetaData) Print() {
	PadPrintf(4, "ExtensionEntryMetaData.ExtDataType: %d\n", entryMetaData.ExtDataType)
	PadPrintf(4, "ExtensionEntryMetaData.ExtDataVersion: %d\n", entryMetaData.ExtDataVersion)
	PadPrintf(4, "ExtensionEntryMetaData.ExtDataStartAddress: %d\n", entryMetaData.ExtDataStartAddress)
	PadPrintf(4, "ExtensionEntryMetaData.ExtDataLength: %d\n", entryMetaData.ExtDataLength)
}
