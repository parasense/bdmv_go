package indx

import (
	"encoding/binary"
	"fmt"
	"io"
)

// ExtensionsMetaData holds metadata about extensions Entries in an MPLS file.
type ExtensionEntryMetaData struct {
	ExtDataType         uint16 // 2-bytes
	ExtDataVersion      uint16 // 2-bytes
	ExtDataStartAddress uint32 // 4-bytes
	ExtDataLength       uint32 // 4-bytes
}

// ReadExtensionsEntriesMetaData reads the extension entries metadata from the provided io.ReadSeeker.
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

func (e *ExtensionEntryMetaData) String() string {
	return fmt.Sprintf("ExtensionEntryMetaData{ExtDataType: %d, ExtDataVersion: %d, ExtDataStartAddress: %d, ExtDataLength: %d}",
		e.ExtDataType,
		e.ExtDataVersion,
		e.ExtDataStartAddress,
		e.ExtDataLength)
}
