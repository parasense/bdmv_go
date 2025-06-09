package main

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Testing NOTES:
// Blu-ray.com calls PIP "Bonus View" in their search feature.
// https://www.blu-ray.com/movies/search.php?action=search&other_bonusview=1&sortby=relevance
//

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
		return err
	}

	if err := binary.Read(file, binary.BigEndian, &pip.SecondaryVideoRef); err != nil {
		return err
	}

	if _, err := file.Seek(1, io.SeekCurrent); err != nil {
		return err
	}

	var buffer uint8
	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return err
	}
	pip.TimelineType = (buffer & 0xF0) >> 4
	pip.LumaKeyFlag = buffer&0x08 != 0
	pip.TrickPlayFlag = buffer&0x04 != 0

	if _, err := file.Seek(1, io.SeekCurrent); err != nil {
		return err
	}

	if pip.LumaKeyFlag {

		if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
			return err
		}

		if err := binary.Read(file, binary.BigEndian, &pip.UpperLimitLumaKey); err != nil {
			return err
		}

	} else {

		if _, err := file.Seek(2, io.SeekCurrent); err != nil {
			return err
		}
	}

	if _, err := file.Seek(2, io.SeekCurrent); err != nil {
		return err
	}

	if err := binary.Read(file, binary.BigEndian, &pip.DataAddress); err != nil {
		return err
	}
	fmt.Printf("ExtensionPIP: DataAddress: %d \n\n", pip.DataAddress)

	// Sanity checks here (?)
	// _parse_pip_data

	// WARNING: there could be a jump/seek to pip.DataAddress

	if err := binary.Read(file, binary.BigEndian, &pip.NumberOfEntries); err != nil {
		return err
	}

	if pip.NumberOfEntries > 0 {
		pip.PIPEntries = make([]*PIPEntry, pip.NumberOfEntries)
	}

	for i := range pip.PIPEntries {
		if err := binary.Read(file, binary.BigEndian, &pip.PIPEntries[i].Time); err != nil {
			return err
		}

		if err := binary.Read(file, binary.BigEndian, &pip.PIPEntries[i].Xpos); err != nil {
			return err
		}

		if err := binary.Read(file, binary.BigEndian, &pip.PIPEntries[i].Ypos); err != nil {
			return err
		}

		if err := binary.Read(file, binary.BigEndian, &pip.PIPEntries[i].ScaleFactor); err != nil {
			return err
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
