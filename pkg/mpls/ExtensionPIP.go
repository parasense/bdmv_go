package mpls

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Testing NOTES:
// Blu-ray.com calls PIP "Bonus View" in their search feature.
// https://www.blu-ray.com/movies/search.php?action=search&other_bonusview=1
//
// List of potential PIP test titles:
// * Never Ending Story (1984)
// * Battlestar Galactica: The Complete Series 2003 - 2009
// * Army of Darkness 1992
// * Serenity 2005
// * Finding Nemo 2003
// * Shaun of the Dead 2004
// * V for Vendetta 2005

/*
 Sellected For testing:
 	SERENITY/BDMV/PLAYLIST/00000.mpls

*/

// ExtensionPIP implements the ExtensionEntryData interface.
type ExtensionPIP struct {
	Length          uint32
	NumberOfEntries uint16
	PIPEntries      []*PIPEntry
}

// PIPEntry represents a single entry in the PIP extension.
type PIPEntry struct {
	ClipRef           uint16
	SecondaryVideoRef uint8
	TimelineType      uint8 // 0b11110000
	LumaKeyFlag       bool  // 0b00001000
	TrickPlayFlag     bool  // 0b00000100
	UpperLimitLumaKey uint8
	DataAddress       uint32
	Data              *PIPData
}

// PIPData represents the data structure for the PIP extension.
type PIPData struct {
	NumberOfEntries uint16
	Entries         []*PIPDataEntry
}

// PIPDataEntry represents a single entry in the PIP data.
type PIPDataEntry struct {
	Time        uint32
	Xpos        uint16         // 0b11111111 0b11110000 0b00000000
	Ypos        uint16         // 0b00000000 0b00001111 0b11111111
	ScaleFactor PIPScalingType // 0b11110000
}

// Read reads the PIP extension data from the provided file at the specified offsets
func (pip *ExtensionPIP) Read(file io.ReadSeeker, offsets *OffsetsUint32, entryMeta *ExtensionEntryMetaData) (err error) {

	// Jump to the start offset
	if _, err := file.Seek(int64(entryMeta.ExtDataStartAddress), io.SeekStart); err != nil {
		return fmt.Errorf("failed to seek Entry Offset Start: (%d); error: %w", entryMeta.ExtDataStartAddress, err)
	}

	startPos, _ := ftell(file)

	binary.Read(file, binary.BigEndian, &pip.Length)
	binary.Read(file, binary.BigEndian, &pip.NumberOfEntries)
	pip.PIPEntries = make([]*PIPEntry, pip.NumberOfEntries)

	// Fill the PIPEntries
	for i := range pip.PIPEntries {
		pip.PIPEntries[i] = &PIPEntry{}
		pip.PIPEntries[i].Read(file)
	}

	// Fill the PIPData
	for i, pipEntry := range pip.PIPEntries {
		pip.PIPEntries[i].Data = &PIPData{}

		// Jump to the data address.
		if _, err := file.Seek(startPos+int64(pipEntry.DataAddress), io.SeekStart); err != nil {
			return fmt.Errorf("failed to seek to PIP data address: %w", err)
		}

		if err := binary.Read(file, binary.BigEndian, &pip.PIPEntries[i].Data.NumberOfEntries); err != nil {
			return fmt.Errorf("failed reading pip.PIPEntries[i].Data.NumberOfEntries: %w", err)
		}

		pip.PIPEntries[i].Data.Entries = make([]*PIPDataEntry, pip.PIPEntries[i].Data.NumberOfEntries)

		for j := range pip.PIPEntries[i].Data.Entries {
			pip.PIPEntries[i].Data.Entries[j] = &PIPDataEntry{}

			if err := binary.Read(file, binary.BigEndian, &pip.PIPEntries[i].Data.Entries[j].Time); err != nil {
				return fmt.Errorf("failed reading PIPEntry.Time: %w", err)
			}

			var buffer [4]byte
			if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
				return fmt.Errorf("failed reading PIPEntry.Xpos: %w", err)
			}
			_tmp := uint32(buffer[0])<<16 | uint32(buffer[1])<<8 | uint32(buffer[2])
			pip.PIPEntries[i].Data.Entries[j].Xpos = uint16((_tmp & 0xFFF000) >> 12)
			pip.PIPEntries[i].Data.Entries[j].Ypos = uint16(_tmp & 0x000FFF)
			pip.PIPEntries[i].Data.Entries[j].ScaleFactor = PIPScalingType((buffer[3] & 0xF0) >> 4)
		}
	}

	return nil
}

// Read reads the PIPEntry data from the provided file.
// It expects the file (io.ReadSeeker) to be positioned at the start of the
// PIPEntry structure.
func (pipEntry *PIPEntry) Read(file io.ReadSeeker) error {

	if err := binary.Read(file, binary.BigEndian, &pipEntry.ClipRef); err != nil {
		return fmt.Errorf("failed reading ExtensionPIP.ClipRef: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &pipEntry.SecondaryVideoRef); err != nil {
		return fmt.Errorf("failed reading ExtensionPIP.SecondaryVideoRef: %w", err)
	}

	if _, err := file.Seek(1, io.SeekCurrent); err != nil {
		return fmt.Errorf("failed to seek past reserve space: %w", err)
	}

	var buffer uint8
	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return fmt.Errorf("failed reading buffer: %w", err)
	}
	pipEntry.TimelineType = (buffer & 0xF0) >> 4
	pipEntry.LumaKeyFlag = buffer&0x08 != 0
	pipEntry.TrickPlayFlag = buffer&0x04 != 0

	if _, err := file.Seek(1, io.SeekCurrent); err != nil {
		return fmt.Errorf("failed to seek past reserve space: %w", err)
	}

	if pipEntry.LumaKeyFlag {

		if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
			return fmt.Errorf("failed reading buffer: %w", err)
		}

		if err := binary.Read(file, binary.BigEndian, &pipEntry.UpperLimitLumaKey); err != nil {
			return fmt.Errorf("failed reading ExtensionPIP.UpperLimitLumaKey: %w", err)
		}

	} else {

		if _, err := file.Seek(2, io.SeekCurrent); err != nil {
			return fmt.Errorf("failed to seek past reserve space: %w", err)
		}
	}

	if _, err := file.Seek(2, io.SeekCurrent); err != nil {
		return fmt.Errorf("failed to seek past reserve space: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &pipEntry.DataAddress); err != nil {
		return fmt.Errorf("failed reading ExtensionPIP.DataAddress: %w", err)
	}

	return nil
}
