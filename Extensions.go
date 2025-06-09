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
	MetaData        *ExtensionsMetaData
	EntriesMetaData []*ExtensionEntryMetaData
	EntriesData     []ExtensionEntryData
}

func ReadExtensions(file io.ReadSeeker, offsets *OffsetsUint32) (extensions *Extensions, err error) {
	extensions = &Extensions{}

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

	fmt.Printf("DEBUG: MetaData: MetaData.Length == %d\n", extensions.MetaData.Length)
	fmt.Printf("DEBUG: MetaData: MetaData.EntryDataStartAddr == %d\n", extensions.MetaData.EntryDataStartAddr)
	fmt.Printf("DEBUG: MetaData: MetaData.EntryDataCount == %d\n\n", extensions.MetaData.EntryDataCount)
	fmt.Printf("DEBUG: MetaData: offsets.Start + MetaData.Length + 4 == %d\n\n", offsets.Start+4+int64(extensions.MetaData.Length))

	// Sanity check
	// These should sum together to equal the EOF.
	if offsets.Start+4+int64(extensions.MetaData.Length) != offsets.Stop {
		return nil, fmt.Errorf("Not enough byte remain to parse extension metadata")
	}

	// Read the extension entries metadata
	// 12-bytes for-each metadata entry
	if extensions.EntriesMetaData, err = ReadExtensionsEntriesMetaData(file, offsets, extensions.MetaData); err != nil {
		return nil, err
	}

	// Read the actual extension data entries
	extensions.EntriesData, err = NEWReadExtensionEntryData(file, offsets, extensions.MetaData, &extensions.EntriesMetaData)

	return extensions, err
}

func (extensions *Extensions) Print() {
	PadPrintln(0, "Extensions:")
	PadPrintln(2, "Extensions MetaData:")
	extensions.MetaData.Print()
	PadPrintln(2, "---")
	for i, metaData := range extensions.EntriesMetaData {
		PadPrintf(2, "Extension Entry MetaData [%d]:\n", i)
		metaData.Print()
		fmt.Println()
	}
	PadPrintln(2, "---")
	for i, entryData := range extensions.EntriesData {
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
