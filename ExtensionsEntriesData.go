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

func NEWReadExtensionEntryData(file io.ReadSeeker, offsets *OffsetsUint32, metaData *ExtensionsMetaData, entriesMetaData *[]*ExtensionEntryMetaData) (entriesData []ExtensionEntryData, err error) {

	entriesData = make([]ExtensionEntryData, metaData.EntryDataCount)

	for i, entryMeta := range *entriesMetaData {

		// Note these offsets are from the start of the file, unline the ExtDataStartAddress which is from the start of the extension.
		//entryOffsetBegin, _ := ftell(file)
		//entryOffsetEnd := entryOffsetBegin + int64(entryMeta.ExtDataLength)
		entryOffsetEnd, err := CalculateEndOffset(file, entryMeta.ExtDataLength)
		if err != nil {
			return nil, err
		}

		// DELETE ME (later when finished)
		//fmt.Printf("DEBUG: ReadExtensionEntryData: metaData.ExtDataTyp == %d\n", entryMeta.ExtDataType)
		//fmt.Printf("DEBUG: ReadExtensionEntryData: metaData.ExtDataVersion == %d\n", entryMeta.ExtDataVersion)
		//fmt.Printf("DEBUG: ReadExtensionEntryData: metaData.ExtDataStartAddress == %d\n", entryMeta.ExtDataStartAddress)
		//fmt.Printf("DEBUG: ReadExtensionEntryData: metaData.ExtDataLength == %d\n", entryMeta.ExtDataLength)
		//fmt.Printf("DEBUG: ReadExtension: DataBlockBegin == %d\n", offsets.Start)
		//fmt.Printf("DEBUG: ReadExtensionEntryData: OffsetBegin == %d\n", entryOffsetBegin)
		//fmt.Printf("DEBUG: ReadExtensionEntryData: OffsetEnd == %d\t(%d + %d)\n", entryOffsetEnd, entryOffsetBegin, int64(entryMeta.ExtDataLength))
		//fmt.Printf("DEBUG: ReadExtension: DataBlockFinish == %d\n", offsets.Stop)

		// XXX - Sanity check here for length checks of extension data block entries.
		//     - This would eliminate the need for sanity checks down in the extension parsers.

		// Do we understand this extension?
		// If yes, process the extension.
		// If not, skip the extension.
		switch {

		case entryMeta.ExtDataType == 1 && entryMeta.ExtDataVersion == 1:
			// PiP metadata extension
			// XXX - This one has not been fully tested yet!
			fmt.Println("DEBUG: PIP extension")
			entriesData[i] = &ExtensionPIP{}

		case entryMeta.ExtDataType == 2 && entryMeta.ExtDataVersion == 1:
			// MVC (3D) STNs extension
			// XXX - This one is partially implemented
			fmt.Println("DEBUG: MVC (3D) STNs extension")
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
				return nil, err
			}
		}

		// In all cases we seek to the end boundary.
		if _, err := file.Seek(entryOffsetEnd, io.SeekStart); err != nil {
			return nil, err
		}

	}

	return entriesData, nil
}
