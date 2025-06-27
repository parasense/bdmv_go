package mpls

import (
	"encoding/binary"
	"fmt"
	"io"
)

// ExtensionStaticMetaData implements the ExtensionEntryData interface.
type ExtensionStaticMetaData struct {
	Length  uint32
	Count   uint8
	Entries []*StaticMetaDataEntry
}

// StaticMetaDataEntry represents a single entry in the static metadata extension.
type StaticMetaDataEntry struct {
	DynamicRangeType             uint8     // 4-bits high
	DisplayPrimariesX            [3]uint16 // 48-bits total
	DisplayPrimariesY            [3]uint16 // 48-bits total
	WhitePointX                  uint16
	WhitePointY                  uint16
	MaxDisplayMasteringLuminance uint16
	MinDisplayMasteringLuminance uint16
	MaxCLL                       uint16
	MaxFALL                      uint16
}

// Read reads the ExtensionStaticMetaData from the provided file.
func (staticMetaData *ExtensionStaticMetaData) Read(file io.ReadSeeker, offsets *OffsetsUint32, entryMeta *ExtensionEntryMetaData) (err error) {

	// Jump to the start offset
	if _, err := file.Seek(int64(entryMeta.ExtDataStartAddress), io.SeekStart); err != nil {
		return fmt.Errorf("failed to seek Entry Offset Start: (%d); error: %w", entryMeta.ExtDataStartAddress, err)
	}

	if err := binary.Read(file, binary.BigEndian, &staticMetaData.Length); err != nil {
		return fmt.Errorf("failed to read ExtensionStaticMetaData.Length: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &staticMetaData.Count); err != nil {
		return fmt.Errorf("failed to read ExtensionStaticMetaData.Count: %w", err)
	}

	// skip 3-bytes reserve
	if _, err := file.Seek(3, io.SeekCurrent); err != nil {
		return fmt.Errorf("failed to read seek past reserve space: %w", err)
	}

	if staticMetaData.Count > 0 {
		staticMetaData.Entries = make([]*StaticMetaDataEntry, staticMetaData.Count)
	}

	return nil
}

// Read reads the StaticMetaDataEntry from the provided file.
// It expects the file (io.ReadSeeker) to be positioned at the start of the
// StaticMetaDataEntry structure
func (smEntry *StaticMetaDataEntry) Read(file io.ReadSeeker) (err error) {

	var flagBuffer uint8
	if err := binary.Read(file, binary.BigEndian, &flagBuffer); err != nil {
		return fmt.Errorf("failed to read flagBuffer: %w", err)
	}
	smEntry.DynamicRangeType = (flagBuffer & 0xF0) >> 4

	// skip 3-bytes reserve
	if _, err := file.Seek(3, io.SeekCurrent); err != nil {
		return fmt.Errorf("failed to read seek past reserve space: %w", err)
	}

	for i := range 3 {
		if err := binary.Read(file, binary.BigEndian, &smEntry.DisplayPrimariesX[i]); err != nil {
			return fmt.Errorf("failed to read StaticMetaDataEntry.DisplayPrimariesX: %w", err)
		}

		if err := binary.Read(file, binary.BigEndian, &smEntry.DisplayPrimariesY[i]); err != nil {
			return fmt.Errorf("failed to read StaticMetaDataEntry.DisplayPrimariesY: %w", err)
		}
	}

	if err := binary.Read(file, binary.BigEndian, &smEntry.WhitePointX); err != nil {
		return fmt.Errorf("failed to read StaticMetaDataEntry.WhitePointX: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &smEntry.WhitePointY); err != nil {
		return fmt.Errorf("failed to read StaticMetaDataEntry.WhitePointY: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &smEntry.MaxDisplayMasteringLuminance); err != nil {
		return fmt.Errorf("failed to read StaticMetaDataEntry.MaxDisplayMasteringLuminance: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &smEntry.MinDisplayMasteringLuminance); err != nil {
		return fmt.Errorf("failed to read StaticMetaDataEntry.MinDisplayMasteringLuminance: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &smEntry.MaxCLL); err != nil {
		return fmt.Errorf("failed to read StaticMetaDataEntry.MaxCLL: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &smEntry.MaxFALL); err != nil {
		return fmt.Errorf("failed to read StaticMetaDataEntry.MaxFALL: %w", err)
	}

	return nil
}
