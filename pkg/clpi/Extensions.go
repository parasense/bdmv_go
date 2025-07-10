package clpi

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

// ReadExtensions reads the extensions from the provided io.ReadSeeker.
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

func (e *Extensions) String() string {
	return fmt.Sprintf(
		"Extensions{"+
			"MetaData: %s, "+
			"EntriesMetaData: %d, "+
			"EntriesData: %d, "+
			"}",
		e.MetaData,
		len(e.EntriesMetaData),
		len(e.EntriesData),
	)
}
