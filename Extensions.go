package main

import (
	"fmt"
	"io"
)

// This is the new extension organization
// There are three distinct regions.
// * Metadata for the entire extensions.
// * Metadata for the individual extension entries.
// * Data for the Extension entrty.

type Extensions struct {
	MetaData      *ExtensionsMetaData
	EntryMetaData []*ExtensionEntryMetaData
	EntryData     []ExtensionEntryData
}

func ReadExtensions(file io.ReadSeeker, offsets *OffsetsUint32) (extensions *Extensions, err error) {
	extensions = &Extensions{}

	// Jump to start address
	if _, err := file.Seek(offsets.Start, io.SeekStart); err != nil {
		return nil, fmt.Errorf("failed to seek to start address: %w\n", err)
	}

	// Read the metadata about the extensions
	// 12-bytes total
	extensions.MetaData, err = ReadMetaData(file)

	// Read the extension entries metadata
	// 12-bytes for-each metadata entry
	extensions.EntryMetaData = make([]*ExtensionEntryMetaData, extensions.MetaData.EntryDataCount)
	for i := range extensions.EntryMetaData {
		extensions.EntryMetaData[i], err = ReadEntryMetaData(file)
	}

	// Read the actual extension data entries
	extensions.EntryData = make([]ExtensionEntryData, extensions.MetaData.EntryDataCount)
	for i, entryMeta := range extensions.EntryMetaData {
		extensions.EntryData[i], err = ReadExtensionEntryData(file, entryMeta)
	}

	return extensions, err
}

func (extensions *Extensions) Print() {
	PadPrintln(0, "Extensions:")
	PadPrintln(2, "Extensions MetaData:")
	extensions.MetaData.Print()
	PadPrintln(2, "---")
	for i, metaData := range extensions.EntryMetaData {
		PadPrintf(2, "Extension Entry MetaData [%d]:\n", i)
		metaData.Print()
		fmt.Println()
	}
	PadPrintln(2, "---")
	for i, entryData := range extensions.EntryData {
		PadPrintf(2, "Extension Entry Data [%d]:\n", i)
		if entryData == nil {
			PadPrintln(4, "[Empty Extension payload]")
			continue
		} else {
			entryData.Print()
		}
		PadPrintln(2, "---")
	}
}
