package mpls

import (
	"encoding/binary"
	"fmt"
	"io"
)

// SubPlayItem represents a sub-play item in an MPLS file.
type SubPlayItem struct {
	Length                   uint16
	FileName                 [5]byte
	Codec                    [4]byte // 27-bits reserve space after
	ConnectionCondition      uint8   // 0b00011110
	IsMultiClipEntries       bool    // 0b00000001
	RefToSTCID               uint8
	INTime                   uint32
	OUTTime                  uint32
	SyncPlaytItemID          uint16
	SyncStartPTS             uint32
	NumberOfMultiClipEntries uint8
	MultiClipEntries         []*PlayItemEntry
}

// ReadSubPlayItem reads a SubPlayItem from the provided io.ReadSeeker.
// It expects the file (io.ReadSeeker) to be positioned at the start of the
// SubPlayItem structure.
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

	subPlayItem.ConnectionCondition = (buffer & 0x1e) >> 1 // 0b00011110
	subPlayItem.IsMultiClipEntries = buffer&0x01 != 0      // 0b00000001

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

func (subPlayItem *SubPlayItem) String() string {
	return fmt.Sprintf("SubPlayItem:\n"+
		"  Length: %d\n"+
		"  FileName: %s\n"+
		"  Codec: %s\n"+
		"  ConnectionCondition: %d\n"+
		"  IsMultiClipEntries: %t\n"+
		"  RefToSTCID: %d\n"+
		"  INTime: %d\n"+
		"  OUTTime: %d\n"+
		"  SyncPlaytItemID: %d\n"+
		"  SyncStartPTS: %d\n"+
		"  NumberOfMultiClipEntries: %d",
		subPlayItem.Length,
		subPlayItem.FileName,
		subPlayItem.Codec,
		subPlayItem.ConnectionCondition,
		subPlayItem.IsMultiClipEntries,
		subPlayItem.RefToSTCID,
		subPlayItem.INTime,
		subPlayItem.OUTTime,
		subPlayItem.SyncPlaytItemID,
		subPlayItem.SyncStartPTS,
		subPlayItem.NumberOfMultiClipEntries)
}
