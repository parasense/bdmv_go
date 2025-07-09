package clpi

import (
	"encoding/binary"
	"fmt"
	"io"
)

// ExtensionExtentStartPoints implements the ExtensionEntryData interface.
type ExtensionExtentStartPoints struct {
	Length         uint32
	NumberOfPoints uint32
	PointEntries   []*PointEntry
}

// PointEntry represents a single entry in the HEVC extension.
type PointEntry struct {
	Point uint32
}

// Read reads the ExtensionHEVC from the provided file.
func (esp *ExtensionExtentStartPoints) Read(file io.ReadSeeker, offsets *OffsetsUint32, entryMeta *ExtensionEntryMetaData) (err error) {

	// Jump to the start offset
	if _, err := file.Seek(offsets.Start+int64(entryMeta.ExtDataStartAddress), io.SeekStart); err != nil {
		return fmt.Errorf("failed to seek Entry Offset Start: (%d); error: %w", entryMeta.ExtDataStartAddress, err)
	}

	if err := binary.Read(file, binary.BigEndian, &esp.Length); err != nil {
		return fmt.Errorf("failed to read esp.Length: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &esp.NumberOfPoints); err != nil {
		return fmt.Errorf("failed to read esp.NumberOfPoints: %w", err)
	}

	esp.PointEntries = make([]*PointEntry, esp.NumberOfPoints)
	for i := range esp.PointEntries {
		esp.PointEntries[i] = &PointEntry{}
		if err := binary.Read(file, binary.BigEndian, &esp.PointEntries[i].Point); err != nil {
			return fmt.Errorf("failed to read [%d]*Point entry: %w", i, err)
		}
	}

	return nil
}

func (esp *ExtensionExtentStartPoints) String() string {
	return fmt.Sprintf(
		"{Length: %d, NumberOfPoints: %d, PointEntries: %s, }",
		esp.Length, esp.NumberOfPoints, esp.PointEntries,
	)
}

func (pnt *PointEntry) String() string {
	return fmt.Sprintf("{Point: %d}", pnt.Point)
}
