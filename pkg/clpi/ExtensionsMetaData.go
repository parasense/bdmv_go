package clpi

import (
	"encoding/binary"
	"fmt"
	"io"
)

// ExtensionsMetaData holds metadata about extensions in an MPLS file.
type ExtensionsMetaData struct {
	Length             uint32
	EntryDataStartAddr uint32
	EntryDataCount     uint8
}

// ReadMetaData reads the ExtensionsMetaData from the provided io.ReadSeeker.
// It expects the file (io.ReadSeeker) to be positioned at the start of the
// ExtensionsMetaData structure.
// The ExtensionsMetaData consists of a length, the entry data start addr,
// and the entry data count. This information is subsequently used to locate and read
// the extension entries MetaData in the MPLS file.
func ReadMetaData(file io.ReadSeeker) (metaData *ExtensionsMetaData, err error) {
	metaData = &ExtensionsMetaData{}

	if err := binary.Read(file, binary.BigEndian, &metaData.Length); err != nil {
		return nil, fmt.Errorf("failed to read ExtensionsMetaData.Length: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &metaData.EntryDataStartAddr); err != nil {
		return nil, fmt.Errorf("failed to read ExtensionsMetaData.EntryDataStartAddr: %w", err)
	}

	// Skip 3-byte reserve space
	if _, err := file.Seek(3, io.SeekCurrent); err != nil {
		return nil, fmt.Errorf("failed to read seek past reserve space: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &metaData.EntryDataCount); err != nil {
		return nil, fmt.Errorf("failed to read ExtensionsMetaData.EntryDataCount: %w", err)
	}

	return metaData, err
}

func (e *ExtensionsMetaData) String() string {
	return fmt.Sprintf("ExtensionsMetaData{Length: %d, EntryDataStartAddr: %d, EntryDataCount: %d}",
		e.Length,
		e.EntryDataStartAddr,
		e.EntryDataCount)
}
