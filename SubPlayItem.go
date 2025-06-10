package main

import (
	"encoding/binary"
	"fmt"
	"io"
)

type SubPlayItem struct {
	Length                   uint16
	FileName                 [5]byte
	Codec                    [4]byte // 27-bits reserve space after
	ConnectionCondition      uint8   // 4-bits
	IsMultiClipEntries       bool
	RefToSTCID               uint8
	INTime                   uint32
	OUTTime                  uint32
	SyncPlaytItemID          uint16
	SyncStartPTS             uint32
	NumberOfMultiClipEntries uint8
	MultiClipEntries         []*PlayItemEntry
}

func ReadSubPlayItem(file io.ReadSeeker) (subPlayItem *SubPlayItem, err error) {
	subPlayItem = &SubPlayItem{}

	if err := binary.Read(file, binary.BigEndian, &subPlayItem.Length); err != nil {
		return nil, fmt.Errorf("failed to read stream info length: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &subPlayItem.FileName); err != nil {
		return nil, fmt.Errorf("failed to read stream info FileName: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &subPlayItem.Codec); err != nil {
		return nil, fmt.Errorf("failed to read stream info Codec: %w", err)
	}

	// Skip 3-byte reserve space
	if _, err := file.Seek(3, io.SeekCurrent); err != nil {
		return nil, fmt.Errorf("failed to seek past reserve space: %w", err)
	}

	var buffer byte
	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return nil, fmt.Errorf("failed to read byte buffer: %w", err)
	}

	subPlayItem.ConnectionCondition = (buffer & 0x1e) >> 1 // 4-bits, 0b00011110
	subPlayItem.IsMultiClipEntries = buffer&0x01 != 0      // 1-bit,  0b00000001

	if err := binary.Read(file, binary.BigEndian, &subPlayItem.INTime); err != nil {
		return nil, fmt.Errorf("failed to read INTime: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &subPlayItem.OUTTime); err != nil {
		return nil, fmt.Errorf("failed to read OUTTime: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &subPlayItem.SyncPlaytItemID); err != nil {
		return nil, fmt.Errorf("failed to read SyncPlaytItemID: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &subPlayItem.SyncStartPTS); err != nil {
		return nil, fmt.Errorf("failed to read SyncStartPTS: %w", err)
	}

	if subPlayItem.IsMultiClipEntries {
		if err := binary.Read(file, binary.BigEndian, &subPlayItem.NumberOfMultiClipEntries); err != nil {
			return nil, fmt.Errorf("failed to read NumberOfMultiClipEntries: %w", err)
		}

		if subPlayItem.NumberOfMultiClipEntries < 1 {
			subPlayItem.NumberOfMultiClipEntries = 1
		}

		subPlayItem.MultiClipEntries = make([]*PlayItemEntry, subPlayItem.NumberOfMultiClipEntries)

		subPlayItem.MultiClipEntries[0] = &PlayItemEntry{
			FileName:   subPlayItem.FileName,
			Codec:      subPlayItem.Codec,
			RefToSTCID: subPlayItem.RefToSTCID,
		}
		for i := uint8(1); i < subPlayItem.NumberOfMultiClipEntries; i++ {
			subPlayItem.MultiClipEntries[i], err = ReadPlayItemEntry(file)
			if err != nil {
				return nil, fmt.Errorf("failed to read PlayItemEntry: %w", err)
			}
		}
	}

	return subPlayItem, nil
}

func (subPlayItem *SubPlayItem) Print() {
	inTime := Convert45KhzTimeToSeconds(subPlayItem.INTime)
	outTime := Convert45KhzTimeToSeconds(subPlayItem.OUTTime)
	PadPrintf(8, "Length: %d\n", subPlayItem.Length)
	PadPrintf(8, "FileName: %s\n", subPlayItem.FileName)
	PadPrintf(8, "Codec: %s\n", subPlayItem.Codec)
	PadPrintf(8, "ConnectionCondition: %d\n", subPlayItem.ConnectionCondition)
	PadPrintf(8, "IsMultiClipEntries: %v\n", subPlayItem.IsMultiClipEntries)
	PadPrintf(8, "RefToSTCID: %d\n", subPlayItem.RefToSTCID)
	PadPrintf(8, "InTime: %d (%d)\n", subPlayItem.INTime, inTime)
	PadPrintf(8, "OUTime: %d (%d)\n", subPlayItem.OUTTime, outTime)
	PadPrintf(10, "*Duration: %v\n", outTime-inTime)
	PadPrintf(8, "SyncPlaytItemID: %d\n", subPlayItem.SyncPlaytItemID)
	PadPrintf(8, "SyncStartPTS: %d\n", subPlayItem.SyncStartPTS)
	PadPrintf(8, "NumberOfMultiClipEntries: %d\n", subPlayItem.NumberOfMultiClipEntries)
	PadPrintf(8, "MultiClipEntries:\n")
	for i := uint8(0); i < uint8(len(subPlayItem.MultiClipEntries)); i++ {
		PadPrintf(6, "MultiClipEntry [%d]:\n", i)
		subPlayItem.MultiClipEntries[i].Print()
	}

}
