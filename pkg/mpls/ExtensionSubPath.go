package mpls

import (
	"encoding/binary"
	"fmt"
	"io"
)

// ExtensionSubPath implements the ExtensionEntryData interface.
type ExtensionSubPath struct {
	Length   uint32
	Count    uint16
	SubPaths []*SubPath
}

// Read reads the ExtensionSubPath data from the provided file at the specified offsets
// and entry metadata. It expects the file (io.ReadSeeker) to be positioned at the start of the
// ExtensionSubPath structure. The function reads the length and count of sub paths,
// followed by the sub paths themselves. It returns an error if any occurs during reading.
func (extensionSubPath *ExtensionSubPath) Read(file io.ReadSeeker, offsets *OffsetsUint32, entryMeta *ExtensionEntryMetaData) (err error) {

	// Jump to the start offset
	if _, err := file.Seek(int64(entryMeta.ExtDataStartAddress), io.SeekStart); err != nil {
		return fmt.Errorf("failed to seek Entry Offset Start: (%d); error: %w", entryMeta.ExtDataStartAddress, err)
	}

	if err := binary.Read(file, binary.BigEndian, &extensionSubPath.Length); err != nil {
		return fmt.Errorf("failed to read extensionSubPath.Length: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &extensionSubPath.Count); err != nil {
		return fmt.Errorf("failed to read extensionSubPath.Count: %w", err)
	}

	extensionSubPath.SubPaths = make([]*SubPath, extensionSubPath.Count)
	for i := range extensionSubPath.SubPaths {
		if extensionSubPath.SubPaths[i], err = ReadSubPath(file); err != nil {
			return fmt.Errorf("failed calling ReadSubPath() in ExtensionSubPath.Read(): %w", err)
		}
	}

	return nil
}
