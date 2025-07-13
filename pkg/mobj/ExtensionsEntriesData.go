package mobj

import (
	"io"
	"log"
)

// Because there are many different types of extension data...
// Put them behind an interface and make the different extensions implement the interface.
// Each extension has access to it's own metadata, and the overall extensions start/stop offsets.
// That way each extension can calculate boundaries.
type ExtensionEntryData interface {
	Read(io.ReadSeeker, *OffsetsUint32, *ExtensionEntryMetaData) error
}

// ReadExtensionEntryData reads the extension entry data from the provided io.ReadSeeker.
// It expects the file (io.ReadSeeker) to be positioned at the start of the
// ExtensionEntryData structure.
// The function reads the extension entry data based on the metadata provided in
// the ExtensionsMetaData and the entriesMetaData slices.
// It returns a slice of ExtensionEntryData and an error if any occurs during reading.
func ReadExtensionEntryData(file io.ReadSeeker, offsets *OffsetsUint32, metaData *ExtensionsMetaData, entriesMetaData *[]*ExtensionEntryMetaData) (entriesData []ExtensionEntryData, err error) {

	// Create N slices for the number of extension entries.
	entriesData = make([]ExtensionEntryData, metaData.EntryDataCount)

	for i, entryMeta := range *entriesMetaData {

		switch {
		case // PiP metadata extension
			entryMeta.ExtDataType == 1 && entryMeta.ExtDataVersion == 1:
			//entriesData[i] = &ExtensionPIP{}

		case // MVC (3D) STNs extension
			entryMeta.ExtDataType == 2 && entryMeta.ExtDataVersion == 1:
			//entriesData[i] = &ExtensionMVCStream{}

		case // SubPath entries extension
			entryMeta.ExtDataType == 2 && entryMeta.ExtDataVersion == 2:
			//entriesData[i] = &ExtensionSubPath{}

		case // Static metadata extension
			entryMeta.ExtDataType == 3 && entryMeta.ExtDataVersion == 5:
			//entriesData[i] = &ExtensionStaticMetaData{}
		}

		// Skip any unimplemented extensions.
		if entriesData[i] != nil {

			if err = entriesData[i].Read(file, offsets, entryMeta); err != nil {
				// XXX - if the extension fails, it's not fatal.
				// xxx - because some extensions might have errors (MVC mostly)
				//return nil, fmt.Errorf("failed to read ExtensionEntryData: %w", err)
				//fmt.Printf("\n\n WARNING! EXTENSION ERROR WHILE PARSING:\t[%+v:%+v]\n\n", entryMeta.ExtDataType, entryMeta.ExtDataVersion)

				// XXX - temporary fatal error to help debugging.
				// This should be removed when all extensions are implemented.
				log.Fatalf("failed to read ExtensionEntryData: %v\n", err)

			}

		} else {
			log.Fatal("Extension not implemented: ", entryMeta.ExtDataType, entryMeta.ExtDataVersion)
		}

	}

	return entriesData, nil
}
