package main

import (
	"fmt"
	"io"
)

// Because there are many different types of extension data
type ExtensionEntryData interface {
	Read(io.ReadSeeker) error
	Print()
}

func ReadExtensionEntryData(file io.ReadSeeker, offsets *OffsetsUint32, metaData *ExtensionsMetaData, entriesMetaData *[]*ExtensionEntryMetaData) (entriesData []ExtensionEntryData, err error) {

	entriesData = make([]ExtensionEntryData, metaData.EntryDataCount)

	for i, entryMeta := range *entriesMetaData {

		end, err := CalculateEndOffset(file, entryMeta.ExtDataLength)
		if err != nil {
			return nil, fmt.Errorf("failed calling CalculateEndOffset(): %w", err)
		}

		switch {

		case entryMeta.ExtDataType == 1 && entryMeta.ExtDataVersion == 1:
			// PiP metadata extension
			entriesData[i] = &ExtensionPIP{}

		case entryMeta.ExtDataType == 2 && entryMeta.ExtDataVersion == 1:
			// MVC (3D) STNs extension
			entriesData[i] = &ExtensionMVCStream{
				numberOfItems: entryMeta.ExtDataLength - 2,
			}

		case entryMeta.ExtDataType == 2 && entryMeta.ExtDataVersion == 2:
			// SubPath entries extension
			entriesData[i] = &ExtensionSubPath{}

		case entryMeta.ExtDataType == 3 && entryMeta.ExtDataVersion == 5:
			// Static metadata extension
			entriesData[i] = &ExtensionStaticMetaData{}
		}

		// Skip any unimplemented extensions.
		if entriesData[i] != nil {
			if err = entriesData[i].Read(file); err != nil {
				return nil, fmt.Errorf("failed to read ExtensionEntryData: %w", err)
			}
		}

		// In all cases we seek to the end boundary.
		if _, err := file.Seek(end, io.SeekStart); err != nil {
			return nil, fmt.Errorf("failed to seek end offset: %w", err)
		}

	}

	return entriesData, nil
}
