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

func ReadExtensionEntryData(file io.ReadSeeker, metaData *ExtensionEntryMetaData) (entryData ExtensionEntryData, err error) {

	// Note these offsets are from the start of the file, unline the ExtDataStartAddress which is from the start of the extension.
	OffsetBegin, _ := ftell(file)
	OffsetEnd := OffsetBegin + int64(metaData.ExtDataLength)

	// Do we understand this extension?
	// If yes, process the extension.
	// If not, skip the extension.

	// DELETE ME (later when finished)
	fmt.Printf("DEBUG: ReadExtensionEntryData: metaData.ExtDataTyp == %d\n", metaData.ExtDataType)
	fmt.Printf("DEBUG: ReadExtensionEntryData: metaData.ExtDataVersion == %d\n", metaData.ExtDataVersion)

	fmt.Printf("DEBUG: ReadExtensionEntryData: metaData.ExtDataStartAddress == %d\n", metaData.ExtDataStartAddress)
	fmt.Printf("DEBUG: ReadExtensionEntryData: metaData.ExtDataLength == %d\n", metaData.ExtDataLength)

	fmt.Printf("DEBUG: ReadExtensionEntryData: OffsetBegin == %d\n", OffsetBegin)
	fmt.Printf("DEBUG: ReadExtensionEntryData: OffsetEnd == %d\n", OffsetEnd)

	switch {

	// PiP metadata extension
	case metaData.ExtDataType == 1 && metaData.ExtDataVersion == 1:
		fmt.Println("DEBUG: PiP metadata extension")

		// NOT YET, skip!
		fmt.Println("DEBUG: SKIPPING PiP metadata extension")
		file.Seek(int64(metaData.ExtDataLength), io.SeekCurrent)
		return nil, fmt.Errorf("Error: known yet unhandled extension.")

	// MVC (3D) STNs extension
	case metaData.ExtDataType == 2 && metaData.ExtDataVersion == 1:
		fmt.Println("DEBUG: MVC (3D) STNs extension")
		entryData = &MVCStreamExtension{
			numberOfItems: metaData.ExtDataLength - 2,
		}
		//entryData.SetNumberOfItems(metaData.ExtDataLength)
		entryData.Read(file)
		entryData.Print() // DEBUG
		file.Seek(OffsetEnd, io.SeekStart)
		return entryData, nil

	// SubPath entries extension
	case metaData.ExtDataType == 2 && metaData.ExtDataVersion == 2:
		entryData = &SubPathExtension{}
		entryData.Read(file)
		entryData.Print() // DEBUG
		file.Seek(OffsetEnd, io.SeekStart)
		return entryData, nil

	// Static metadata extension
	case metaData.ExtDataType == 3 && metaData.ExtDataVersion == 5:
		fmt.Println("DEBUG: Static metadata extension")

		// NOT YET, skip!
		fmt.Println("DEBUG: SKIPPING Static metadata extension")
		file.Seek(int64(metaData.ExtDataLength), io.SeekCurrent)
		return nil, fmt.Errorf("Error: known yet unhandled extension.")

	// Unknown extension.
	default:

		fmt.Printf("DEBUG: unknown extension: [id1: %d] [id2: %d].\n\n", metaData.ExtDataType, metaData.ExtDataVersion)

		// NOT YET, skip!
		fmt.Println("DEBUG: SKIPPING unknown extension")
		file.Seek(int64(metaData.ExtDataLength), io.SeekCurrent)
		//return nil, fmt.Errorf("Error: unknown extension.")
		return nil, nil
	}
}
