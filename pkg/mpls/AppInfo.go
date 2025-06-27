package mpls

import (
	"encoding/binary"
	"fmt"
	"io"
)

// AppInfo holds application-specific information in an MPLS file.
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

// ReadAppInfo reads the AppInfo structure from the provided io.ReadSeeker.
func ReadAppInfo(file io.ReadSeeker, offsets *OffsetsUint32) (appinfo *AppInfo, err error) {
	appinfo = &AppInfo{}

	// Jump to start address
	if _, err := file.Seek(offsets.Start, io.SeekStart); err != nil {
		return nil, fmt.Errorf("failed to seek to start address: %w\n", err)
	}

	if err := binary.Read(file, binary.BigEndian, &appinfo.Length); err != nil {
		return nil, fmt.Errorf("failed to read appinfo.Length: %w", err)
	}

	// Reserve space 1 byte.
	if _, err := file.Seek(1, io.SeekCurrent); err != nil {
		return nil, fmt.Errorf("failed to seek past reserve byte: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &appinfo.PlaybackType); err != nil {
		return nil, fmt.Errorf("failed to read appinfo.PlaybackType: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &appinfo.PlaybackCount); err != nil {
		return nil, fmt.Errorf("failed to read appinfo.PlaybackCount: %w", err)
	}

	if appinfo.UserOptions, err = ReadUserOptions(file); err != nil {
		return nil, fmt.Errorf("failed to read UserOptions: %w", err)
	}

	// flags 5 bits of 1 byte
	var flagBuffer uint8
	if err := binary.Read(file, binary.BigEndian, &flagBuffer); err != nil {
		return nil, fmt.Errorf("failed to read flagBuffer: %w", err)
	}
	appinfo.RandomAccessFlag = flagBuffer&0x80 != 0
	appinfo.AudioMixFlag = flagBuffer&0x40 != 0
	appinfo.LosslessBypassFlag = flagBuffer&0x20 != 0
	appinfo.MVCBaseViewRFlag = flagBuffer&0x10 != 0
	appinfo.SDRConversionNotificationFlag = flagBuffer&0x08 != 0

	// Reserve space 1 byte.
	if _, err := file.Seek(1, io.SeekCurrent); err != nil {
		return nil, fmt.Errorf("failed to seek past reserve byte: %w", err)
	}
	return appinfo, nil
}

func (appinfo *AppInfo) String() string {
	return fmt.Sprintf("AppInfo:\n"+
		"  Length: %d\n"+
		"  PlaybackType: %d\n"+
		"  PlaybackCount: %d\n"+
		"  UserOptions: %s\n"+
		"  RandomAccessFlag: %t\n"+
		"  AudioMixFlag: %t\n"+
		"  LosslessBypassFlag: %t\n"+
		"  MVCBaseViewRFlag: %t\n"+
		"  SDRConversionNotificationFlag: %t",
		appinfo.Length,
		appinfo.PlaybackType,
		appinfo.PlaybackCount,
		appinfo.UserOptions,
		appinfo.RandomAccessFlag,
		appinfo.AudioMixFlag,
		appinfo.LosslessBypassFlag,
		appinfo.MVCBaseViewRFlag,
		appinfo.SDRConversionNotificationFlag)
}
