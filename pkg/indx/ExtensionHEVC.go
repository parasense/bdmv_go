package indx

import (
	"encoding/binary"
	"fmt"
	"io"
)

// ExtensionHEVC implements the ExtensionEntryData interface.
type ExtensionHEVC struct {
	Length    uint32
	HEVCEntry *HEVCEntry
}

// HEVCEntry represents a single entry in the HEVC extension.
type HEVCEntry struct {
	DiscType        uint8 // 4-bits 0b11110000 & 0xF0 >> 4
	Exists4KFlag    bool  // 1-bit  0b00000001 & 0x01
	HDRPlusFlag     bool  // 1-bit  0b00010000 & 0x10
	DolbyVisionFlag bool  // 1-bit  0b00000100 & 0x04
	HDRFlag         uint8 // 2-bits 0b00000011 & 0x03
}

// Read reads the ExtensionHEVC from the provided file.
func (hevc *ExtensionHEVC) Read(file io.ReadSeeker, offsets *OffsetsUint32, entryMeta *ExtensionEntryMetaData) (err error) {

	// Jump to the start offset
	if _, err := file.Seek(int64(entryMeta.ExtDataStartAddress), io.SeekStart); err != nil {
		return fmt.Errorf("failed to seek Entry Offset Start: (%d); error: %w", entryMeta.ExtDataStartAddress, err)
	}

	if err := binary.Read(file, binary.BigEndian, &hevc.Length); err != nil {
		return fmt.Errorf("failed to read ExtensionHEVC.Length: %w", err)
	}

	var buffer byte
	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return fmt.Errorf("failed to read HEVCEntry buffer: %w", err)
	}
	hevc.HEVCEntry.DiscType = buffer & 0xF0 >> 4   // 4-bits for DiscType
	hevc.HEVCEntry.Exists4KFlag = buffer&0x01 != 0 // 1-bit for Exists4KFlag

	// skip 1-bytes reserve
	if _, err := file.Seek(1, io.SeekCurrent); err != nil {
		return fmt.Errorf("failed to read seek past reserve space: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return fmt.Errorf("failed to read HEVCEntry buffer: %w", err)
	}
	hevc.HEVCEntry.HDRPlusFlag = buffer&0x10 != 0     // 1-bit for HDRPlusFlag
	hevc.HEVCEntry.DolbyVisionFlag = buffer&0x04 != 0 // 1-bit for DolbyVisionFlag
	hevc.HEVCEntry.HDRFlag = buffer & 0x03            // 2-bits for HDRFlag

	// skip 5-bytes reserve
	if _, err := file.Seek(5, io.SeekCurrent); err != nil {
		return fmt.Errorf("failed to read seek past reserve space: %w", err)
	}

	return nil
}

func (hevc *ExtensionHEVC) String() string {
	return fmt.Sprintf(
		"{Length: %d, HEVCEntry: {DiscType: %d, Exists4KFlag: %t, HDRPlusFlag: %t, DolbyVisionFlag: %t, HDRFlag: %d}}",
		hevc.Length,
		hevc.HEVCEntry.DiscType,
		hevc.HEVCEntry.Exists4KFlag,
		hevc.HEVCEntry.HDRPlusFlag,
		hevc.HEVCEntry.DolbyVisionFlag,
		hevc.HEVCEntry.HDRFlag,
	)
}
