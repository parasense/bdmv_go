package indx

import (
	"encoding/binary"
	"fmt"
	"io"
)

type AppInfo struct {
	Length                      uint32
	_                           bool     //  1-bit  0b10000000 & 0x80
	InitialOutputModePreference bool     //  1-bit  0b01000000 & 0x40
	SSContentExistFlag          bool     //  1-bit  0b00100000 & 0x20
	_                           bool     //  1-bit  0b00010000 & 0x10
	InitialDynamicRangeType     uint8    //  4-bits 0b00001111 & 0x0F
	VideoFormat                 uint8    //  4-bits 0b11110000 & 0xF0 >> 4
	FrameRate                   uint8    //  4-bits 0b00001111 & 0x0F
	UserData                    [32]byte // 32-bytes
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

	var flagBuffer uint8
	if err := binary.Read(file, binary.BigEndian, &flagBuffer); err != nil {
		return nil, fmt.Errorf("failed to read flagBuffer: %w", err)
	}
	appinfo.InitialOutputModePreference = flagBuffer&0x40 != 0
	appinfo.SSContentExistFlag = flagBuffer&0x20 != 0
	appinfo.InitialDynamicRangeType = flagBuffer & 0x0F

	if err := binary.Read(file, binary.BigEndian, &flagBuffer); err != nil {
		return nil, fmt.Errorf("failed to read flagBuffer: %w", err)
	}
	appinfo.VideoFormat = flagBuffer & 0xF0 >> 4
	appinfo.FrameRate = flagBuffer & 0x0F

	if err := binary.Read(file, binary.BigEndian, &appinfo.UserData); err != nil {
		return nil, fmt.Errorf("failed to read UserData: %w", err)
	}

	return appinfo, nil
}

func (appinfo *AppInfo) String() string {
	return fmt.Sprintf(
		"{Length: %d, InitialOutputModePreference: %t, SSContentExistFlag: %t, InitialDynamicRangeType: %d, VideoFormat: %d, FrameRate: %d, UserData: %x}",
		appinfo.Length,
		appinfo.InitialOutputModePreference,
		appinfo.SSContentExistFlag,
		appinfo.InitialDynamicRangeType,
		appinfo.VideoFormat,
		appinfo.FrameRate,
		appinfo.UserData,
	)
}
