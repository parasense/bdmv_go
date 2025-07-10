package clpi

import (
	"encoding/binary"
	"fmt"
	"io"
)

type ClipInfo struct {
	Length                           uint32
	ClipStreamType                   uint8
	ApplicationType                  ClipApplicationType
	IsCC5                            bool // 1-bit 0b00000000,00000000,00000000,00000001
	TSRecordingRate                  uint32
	NumberOfSourcePackets            uint32
	TSTypeInfoBlock                  [32]byte
	FollowingClipStreamType          uint8   // 1-byte
	FollowingClipInformationFileName [5]byte // 5-byte
	FollowingClipCodecIdentifier     [4]byte // 4-byte
}

func ReadClipInfo(file io.ReadSeeker, offsets *OffsetsUint32) (clipInfo *ClipInfo, err error) {
	clipInfo = &ClipInfo{}

	// Jump to start address
	if _, err := file.Seek(offsets.Start, io.SeekStart); err != nil {
		return nil, err
	}

	if err := binary.Read(file, binary.BigEndian, &clipInfo.Length); err != nil {
		return nil, err
	}

	// Reserve space 2-bytes.
	if _, err := file.Seek(2, io.SeekCurrent); err != nil {
		return nil, fmt.Errorf("failed to seek to beyond reserved space: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &clipInfo.ClipStreamType); err != nil {
		return nil, fmt.Errorf("failed to read ClipStreamType: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &clipInfo.ApplicationType); err != nil {
		return nil, fmt.Errorf("failed to read ApplicationType: %w", err)
	}

	var buffer uint32
	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return nil, fmt.Errorf("failed to read buffer: %w", err)
	}
	clipInfo.IsCC5 = buffer&0x00000001 != 0

	if err := binary.Read(file, binary.BigEndian, &clipInfo.TSRecordingRate); err != nil {
		return nil, fmt.Errorf("failed to read TSRecordingRate: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &clipInfo.NumberOfSourcePackets); err != nil {
		return nil, fmt.Errorf("failed to read NumberOfSourcePackets: %w", err)
	}

	// Reserve space 128-bytes.
	if _, err := file.Seek(128, io.SeekCurrent); err != nil {
		return nil, fmt.Errorf("failed to seek to beyond reserved space: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &clipInfo.TSTypeInfoBlock); err != nil {
		return nil, fmt.Errorf("failed to read TSTypeInfoBlock: %w", err)
	}

	if clipInfo.IsCC5 {

		// Reserve space 1-byte.
		if _, err := file.Seek(1, io.SeekCurrent); err != nil {
			return nil, fmt.Errorf("failed to seek to beyond reserved space: %w", err)
		}

		if err := binary.Read(file, binary.BigEndian, &clipInfo.FollowingClipStreamType); err != nil {
			return nil, fmt.Errorf("failed to read FollowingClipStreamType: %w", err)
		}

		// Reserve space 4-byte.
		if _, err := file.Seek(4, io.SeekCurrent); err != nil {
			return nil, fmt.Errorf("failed to seek to beyond reserved space: %w", err)
		}

		if err := binary.Read(file, binary.BigEndian, &clipInfo.FollowingClipInformationFileName); err != nil {
			return nil, fmt.Errorf("failed to read FollowingClipInformationFileName: %w", err)
		}

		if err := binary.Read(file, binary.BigEndian, &clipInfo.FollowingClipCodecIdentifier); err != nil {
			return nil, fmt.Errorf("failed to read FollowingClipCodecIdentifier: %w", err)
		}

		// Reserve space 1-byte.
		if _, err := file.Seek(1, io.SeekCurrent); err != nil {
			return nil, fmt.Errorf("failed to seek to beyond reserved space: %w", err)
		}
	}

	return clipInfo, nil
}

func (clipInfo *ClipInfo) String() string {
	return fmt.Sprintf(
		"ClipInfo{"+
			"Length: %d, "+
			"ClipStreamType: %d, "+
			"ApplicationType: %d, "+
			"IsCC5: %t, "+
			"TSRecordingRate: %d, "+
			"NumberOfSourcePackets: %d, "+
			"TSTypeInfoBlock: %x, "+
			"FollowingClipStreamType: %d, "+
			"FollowingClipInformationFileName: %s, "+
			"FollowingClipCodecIdentifier: %s,"+
			"}",
		clipInfo.Length,
		clipInfo.ClipStreamType,
		clipInfo.ApplicationType,
		clipInfo.IsCC5,
		clipInfo.TSRecordingRate,
		clipInfo.NumberOfSourcePackets,
		clipInfo.TSTypeInfoBlock,
		clipInfo.FollowingClipStreamType,
		string(clipInfo.FollowingClipInformationFileName[:]),
		string(clipInfo.FollowingClipCodecIdentifier[:]),
	)
}
