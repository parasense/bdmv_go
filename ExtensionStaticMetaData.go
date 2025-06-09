package main

import (
	"encoding/binary"
	"io"
)

type ExtensionStaticMetaData struct {
	Length  uint32
	Count   uint8
	Entries []*StaticMetaDataEntry
}

type StaticMetaDataEntry struct {
	DynamicRangeType             uint8     // 4-bits high
	DisplayPrimariesX            [3]uint16 // 48-bits
	DisplayPrimariesY            [3]uint16 // 48-bits
	WhitePointX                  uint16
	WhitePointY                  uint16
	MaxDisplayMasteringLuminance uint16
	MinDisplayMasteringLuminance uint16
	MaxCLL                       uint16
	MaxFALL                      uint16
}

func (staticMetaData *ExtensionStaticMetaData) Read(file io.ReadSeeker) (err error) {

	if err := binary.Read(file, binary.BigEndian, &staticMetaData.Length); err != nil {
		return err
	}

	if err := binary.Read(file, binary.BigEndian, &staticMetaData.Count); err != nil {
		return err
	}

	if _, err := file.Seek(3, io.SeekCurrent); err != nil {
		return err
	}

	if staticMetaData.Count > 0 {
		staticMetaData.Entries = make([]*StaticMetaDataEntry, staticMetaData.Count)
	}

	return nil
}

func (smEntry *StaticMetaDataEntry) Read(file io.ReadSeeker) (err error) {

	var flagBuffer uint8
	if err := binary.Read(file, binary.BigEndian, &flagBuffer); err != nil {
		return err
	}
	smEntry.DynamicRangeType = (flagBuffer & 0xF0) >> 4

	// skip 3-bytes reserve
	if _, err := file.Seek(3, io.SeekCurrent); err != nil {
		return err
	}

	for i := range 3 {
		if err := binary.Read(file, binary.BigEndian, &smEntry.DisplayPrimariesX[i]); err != nil {
			return err
		}

		if err := binary.Read(file, binary.BigEndian, &smEntry.DisplayPrimariesY[i]); err != nil {
			return err
		}
	}

	if err := binary.Read(file, binary.BigEndian, &smEntry.WhitePointX); err != nil {
		return err
	}

	if err := binary.Read(file, binary.BigEndian, &smEntry.WhitePointY); err != nil {
		return err
	}

	if err := binary.Read(file, binary.BigEndian, &smEntry.MaxDisplayMasteringLuminance); err != nil {
		return err
	}

	if err := binary.Read(file, binary.BigEndian, &smEntry.MinDisplayMasteringLuminance); err != nil {
		return err
	}

	if err := binary.Read(file, binary.BigEndian, &smEntry.MaxCLL); err != nil {
		return err
	}

	if err := binary.Read(file, binary.BigEndian, &smEntry.MaxFALL); err != nil {
		return err
	}

	return nil
}

func (smExtension *ExtensionStaticMetaData) Print() {
	PadPrintln(4, "Static MetaData Extension")
	PadPrintf(6, "Length: %d\n", smExtension.Length)
	PadPrintf(6, "Count: %d\n", smExtension.Count)
	for i, entry := range smExtension.Entries {
		PadPrintf(6, "Static MEtaData Entry [%d]:\n", i)
		entry.Print()
	}
}

func (smEntry *StaticMetaDataEntry) Print() {
	PadPrintln(4, "Static MetaData Entry")
	PadPrintf(6, "DynamicRangeType: %d\n", smEntry.DynamicRangeType)
	PadPrintf(6, "DisplayPrimariesX: %v\n", smEntry.DisplayPrimariesX)
	PadPrintf(6, "DisplayPrimariesY: %v\n", smEntry.DisplayPrimariesX)
	PadPrintf(6, "WhitePointX: %v\n", smEntry.WhitePointX)
	PadPrintf(6, "WhitePointY: %v\n", smEntry.WhitePointY)
	PadPrintf(6, "MaxDisplayMasteringLuminance: %v\n", smEntry.MaxDisplayMasteringLuminance)
	PadPrintf(6, "MinDisplayMasteringLuminance: %v\n", smEntry.MinDisplayMasteringLuminance)
	PadPrintf(6, "MaxCLL: %v\n", smEntry.MaxCLL)
	PadPrintf(6, "MaxFALL: %v\n", smEntry.MaxFALL)
}
