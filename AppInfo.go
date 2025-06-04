package main

import (
	"encoding/binary"
	"fmt"
	"io"
)

type AppInfo struct {
	Length                        uint32
	PlaybackType                  uint8
	PlaybackCount                 uint16
	UserOptions                   *UserOptions
	RandomAccessFlag              bool
	AudioMixFlag                  bool
	LosslessBypassFlag            bool
	MVCBaseViewRFlag              bool
	SDRConversionNotificationFlag bool
}

func ReadAppInfo(file io.ReadSeeker, offsets *OffsetsUint32) (*AppInfo, error) {
	appinfo := &AppInfo{}

	// Jump to start address
	if _, err := file.Seek(int64(offsets.Start), io.SeekStart); err != nil {
		return nil, fmt.Errorf("failed to seek to start address: %w\n", err)
	}

	if err := binary.Read(file, binary.BigEndian, &appinfo.Length); err != nil {
		return nil, fmt.Errorf("failed to read appinfo.Length: %w", err)
	}

	// Reserve space 1 byte.
	if _, err := file.Seek(1, io.SeekCurrent); err != nil {
		return nil, fmt.Errorf("failed to seek past reserve byte: %v", err)
	}
	if err := binary.Read(file, binary.BigEndian, &appinfo.PlaybackType); err != nil {
		return nil, fmt.Errorf("failed to read appinfo.PlaybackType: %v", err)
	}
	if err := binary.Read(file, binary.BigEndian, &appinfo.PlaybackCount); err != nil {
		return nil, fmt.Errorf("failed to read appinfo.PlaybackCount: %v", err)
	}

	// XXX - no error handling here.
	var err error
	appinfo.UserOptions, err = ReadUserOptions(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read UserOptions: %v", err)
	}

	// flags 5 bits of 1 byte
	var flagBuffer uint8
	if err := binary.Read(file, binary.BigEndian, &flagBuffer); err != nil {
		return nil, fmt.Errorf("failed to read flagBuffer: %v", err)
	}
	appinfo.RandomAccessFlag = flagBuffer&0x80 != 0
	appinfo.AudioMixFlag = flagBuffer&0x40 != 0
	appinfo.LosslessBypassFlag = flagBuffer&0x20 != 0
	appinfo.MVCBaseViewRFlag = flagBuffer&0x10 != 0
	appinfo.SDRConversionNotificationFlag = flagBuffer&0x08 != 0

	// Reserve space 1 byte.
	if _, err := file.Seek(1, io.SeekCurrent); err != nil {
		return nil, fmt.Errorf("failed to seek past reserve byte: %v", err)
	}
	return appinfo, nil
}

func (appinfo *AppInfo) Print() {
	PadPrintln(0, "AppInfo:")
	PadPrintf(2, "Length: %d\n", appinfo.Length)
	PadPrintf(2, "PlaybackType: %d\n", appinfo.PlaybackType)
	PadPrintf(2, "PlaybackCount: %d\n", appinfo.PlaybackCount)
	appinfo.UserOptions.Print()
	PadPrintf(2, "RandomAccessFlag: %v\n", appinfo.RandomAccessFlag)
	PadPrintf(2, "AudioMixFlag: %v\n", appinfo.AudioMixFlag)
	PadPrintf(2, "LosslessBypassFlag: %v\n", appinfo.LosslessBypassFlag)
	PadPrintf(2, "MVCBaseViewRFlag: %v\n", appinfo.MVCBaseViewRFlag)
	PadPrintf(2, "SDRConversionNotificationFlag: %v\n", appinfo.SDRConversionNotificationFlag)
	PadPrintln(2, "---")
}
