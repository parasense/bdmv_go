package clpi

import (
	"encoding/binary"
	"fmt"
	"io"
)

// ExtensionExtentStartPoints implements the ExtensionEntryData interface.
type ExtensionLPCMDownMixCoefficient struct {
	Length         uint32
	NumberOfPoints uint32
	PointEntries   []*PointEntry
}

// Read reads the ExtensionHEVC from the provided file.
func (dmc *ExtensionLPCMDownMixCoefficient) Read(file io.ReadSeeker, offsets *OffsetsUint32, entryMeta *ExtensionEntryMetaData) (err error) {

	// Jump to the start offset
	if _, err := file.Seek(offsets.Start+int64(entryMeta.ExtDataStartAddress), io.SeekStart); err != nil {
		return fmt.Errorf("failed to seek Entry Offset Start: (%d); error: %w", entryMeta.ExtDataStartAddress, err)
	}

	if err := binary.Read(file, binary.BigEndian, &dmc.Length); err != nil {
		return fmt.Errorf("failed to read ExtensionHEVC.Length: %w", err)
	}

	dmc.PointEntries = make([]*PointEntry, dmc.NumberOfPoints)
	for i := range dmc.PointEntries {
		if err := binary.Read(file, binary.BigEndian, &dmc.PointEntries[i].Point); err != nil {
			return fmt.Errorf("failed to read [%d]*Point entry: %w", i, err)
		}
	}

	return nil
}

func (dmc *ExtensionLPCMDownMixCoefficient) String() string {
	return fmt.Sprintf(
		"ExtensionLPCMDownMixCoefficient{"+
			"Length: %d, "+
			"NumberOfPoints: %d, "+
			"PointEntries: %s, "+
			"}",
		dmc.Length,
		dmc.NumberOfPoints,
		dmc.PointEntries,
	)
}
