package main

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

type ExtensionPIP struct {
	ClipRef           uint16
	SecondaryVideoRef uint8
	TimelineType      uint8 // 4-bits 0b11110000
	LumaKeyFlag       bool  // 1-bit  0b00001000
	TrickPlayFlag     bool  // 1-bit  0b00000100
	UpperLimitLumaKey uint8
	DataAddress       uint32
	NumberOfEntries   uint16
	PIPEntries        []*PIPEntry
}

type PIPEntry struct {
	Time        uint32
	Xpos        uint16 // 12-bits high
	Ypos        uint16 // 12-bits low  (this & prev consume 3-bytes)
	ScaleFactor uint8  // 4-bits (high)
}

func (pip *ExtensionPIP) Read(file io.ReadSeeker) (err error) {
	pip = &ExtensionPIP{}

	if err := binary.Read(file, binary.BigEndian, &pip.ClipRef); err != nil {
		return fmt.Errorf("failed reading ExtensionPIP.ClipRef: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &pip.SecondaryVideoRef); err != nil {
		return fmt.Errorf("failed reading ExtensionPIP.SecondaryVideoRef: %w", err)
	}

	if _, err := file.Seek(1, io.SeekCurrent); err != nil {
		return fmt.Errorf("failed to seek past reserve space: %w", err)
	}

	var buffer uint8
	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return fmt.Errorf("failed reading buffer: %w", err)
	}
	pip.TimelineType = (buffer & 0xF0) >> 4
	pip.LumaKeyFlag = buffer&0x08 != 0
	pip.TrickPlayFlag = buffer&0x04 != 0

	if _, err := file.Seek(1, io.SeekCurrent); err != nil {
		return fmt.Errorf("failed to seek past reserve space: %w", err)
	}

	if pip.LumaKeyFlag {

		if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
			return fmt.Errorf("failed reading buffer: %w", err)
		}

		if err := binary.Read(file, binary.BigEndian, &pip.UpperLimitLumaKey); err != nil {
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

	if err := binary.Read(file, binary.BigEndian, &pip.DataAddress); err != nil {
		return fmt.Errorf("failed reading ExtensionPIP.NumberOfEntries: %w", err)
	}
	//fmt.Printf("ExtensionPIP: DataAddress: %d \n\n", pip.DataAddress)

	// Sanity checks here (?)
	// _parse_pip_data

	// WARNING: there could be a jump/seek to pip.DataAddress

	if err := binary.Read(file, binary.BigEndian, &pip.NumberOfEntries); err != nil {
		return fmt.Errorf("failed reading ExtensionPIP.NumberOfEntries: %w", err)
	}

	if pip.NumberOfEntries > 0 {
		pip.PIPEntries = make([]*PIPEntry, pip.NumberOfEntries)
	}

	for i := range pip.PIPEntries {
		if err := binary.Read(file, binary.BigEndian, &pip.PIPEntries[i].Time); err != nil {
			return fmt.Errorf("failed reading PIPEntry.Time: %w", err)
		}

		if err := binary.Read(file, binary.BigEndian, &pip.PIPEntries[i].Xpos); err != nil {
			return fmt.Errorf("failed reading PIPEntry.Xpos: %w", err)
		}

		if err := binary.Read(file, binary.BigEndian, &pip.PIPEntries[i].Ypos); err != nil {
			return fmt.Errorf("failed reading PIPEntry.Ypos: %w", err)
		}

		if err := binary.Read(file, binary.BigEndian, &pip.PIPEntries[i].ScaleFactor); err != nil {
			return fmt.Errorf("failed reading PIPEntry.ScaleFactor: %w", err)
		}
	}

	return nil
}

func (pip *ExtensionPIP) Print() {
	PadPrintln(4, "PIP Extension:")
	PadPrintf(6, "ClipRef: %d\n", pip.ClipRef)
	PadPrintf(6, "SecondaryVideoRef: %d\n", pip.SecondaryVideoRef)
	PadPrintf(6, "TimelineType: %d\n", pip.TimelineType)
	PadPrintf(6, "LumaKeyFlag: %v\n", pip.LumaKeyFlag)
	PadPrintf(6, "TrickPlayFlag: %v\n", pip.TrickPlayFlag)
	PadPrintf(6, "UpperLimitLumaKey: %d\n", pip.UpperLimitLumaKey)
	PadPrintf(6, "DataAddress: %d\n", pip.DataAddress)
}
