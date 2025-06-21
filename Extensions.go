package main

import (
	"fmt"
	"io"
)

// This how extensions are organized...
// There are three distinct regions.
// * Metadata for the entire extensions.
// * Metadata for the individual extension entries.
// * Data for the Extension entrty.

type Extensions struct {
	MetaData        *ExtensionsMetaData
	EntriesMetaData []*ExtensionEntryMetaData
	EntriesData     []ExtensionEntryData
}

func ReadExtensions(file io.ReadSeeker, offsets *OffsetsUint32) (extensions *Extensions, err error) {
	extensions = &Extensions{}

	fmt.Printf("\n\nEXTENSIONS: START OFFSET: %+v\n\n\n", offsets.Start)

	// Jump to start address
	if _, err := file.Seek(offsets.Start, io.SeekStart); err != nil {
		return nil, fmt.Errorf("failed to seek to start address: %w\n", err)
	}

	// Sanity check
	// We know the metadata is 12-bytes.
	if offsets.Start > (offsets.Stop - 12) {
		return nil, fmt.Errorf("Not enough byte remain to parse extension metadata")
	}

	// 12-bytes total
	extensions.MetaData, err = ReadMetaData(file)

	//// XXX
	PadPrintln(0, "Extension (overall) Metadata:")
	PadPrintf(2, "MetaData.Length == %d\n", extensions.MetaData.Length)
	PadPrintf(2, "MetaData.EntryDataStartAddr == %d\n", extensions.MetaData.EntryDataStartAddr)
	PadPrintf(2, "MetaData.EntryDataCount == %d\n", extensions.MetaData.EntryDataCount)
	PadPrintln(0, "---")
	PadPrintf(2, "Caclulated END: %d  (offsets.Start + MetaData.Length + 4)\n", offsets.Start+4+int64(extensions.MetaData.Length))
	PadPrintln(0, "---")

	// Sanity check
	// These should sum together to equal the EOF.
	if offsets.Start+4+int64(extensions.MetaData.Length) != offsets.Stop {
		return nil, fmt.Errorf("Not enough byte remain to parse extension metadata")
	}

	// Read the extension entries metadata
	// 12-bytes for-each metadata entry
	if extensions.EntriesMetaData, err = ReadExtensionsEntriesMetaData(file, offsets, extensions.MetaData); err != nil {
		return nil, fmt.Errorf("failed calling ReadExtensionsEntriesMetaData(): %w", err)
	}

	// Read the actual extension data entries
	extensions.EntriesData, err = ReadExtensionEntryData(file, offsets, extensions.MetaData, &extensions.EntriesMetaData)

	return extensions, err
}

func (extensions *Extensions) Print() {
	PadPrintln(0)
	PadPrintln(0, "Extensions:")
	PadPrintln(2, "Extensions MetaData:")
	extensions.MetaData.Print()
	PadPrintln(2, "---")
	for i, metaData := range extensions.EntriesMetaData {
		PadPrintf(2, "Extension Entry MetaData [%d]:\n", i+1)
		metaData.Print()
		PadPrintln(2, "---")
		if extensions.EntriesData[i] == nil {
			PadPrintln(4, "[Empty Extension payload]")
			continue
		} else {
			PadPrintf(2, "Extension Entry payload [%d]:\n", i+1)
			extensions.EntriesData[i].Print()
		}
		PadPrintln(2, "---")
		PadPrintln(2, "---")
	}
}
