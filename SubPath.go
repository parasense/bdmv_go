package main

/*
Dev Notes:

While testing a particular bluray title: Never Ending Story (1984).
It was throwing an error, because apparently that title pads SubPlayItems.
That was strange, until then no other test title did that!

There was exactly 1-byte of zero fill padding at the end, causing alignment havok.
So Please note that apparently it is allowed to have padding here.
Documentation is scarce, but I've seen no indication this allows padding reserve space.
IT DOES ALLOW FOR ARBITRARY PADDING or RESERVE SPACE.
*/

import (
	"encoding/binary"
	"fmt"
	"io"
)

type SubPath struct {
	Length               uint32
	SubPathType          uint8
	IsRepeatSubPath      bool
	NumberOfSubPlayItems uint8
	SubPlayItems         []*SubPlayItem
}

func ReadSubPath(file io.ReadSeeker) (subPath *SubPath, err error) {
	subPath = &SubPath{}

	if err := binary.Read(file, binary.BigEndian, &subPath.Length); err != nil {
		return nil, fmt.Errorf("failed to read stream info length: %v\n", err)
	}

	end, err := CalculateEndOffset(file, subPath.Length)
	if err != nil {
		return nil, err
	}

	// Skip 1-byte reserve space
	if _, err := file.Seek(1, io.SeekCurrent); err != nil {
		return nil, fmt.Errorf("failed to seek past reserve space: %v\n", err)
	}

	if err := binary.Read(file, binary.BigEndian, &subPath.SubPathType); err != nil {
		return nil, fmt.Errorf("failed to read stream info SubPathType: %v\n", err)
	}

	// Skip 1-byte reserve space
	if _, err := file.Seek(1, io.SeekCurrent); err != nil {
		return nil, fmt.Errorf("failed to seek past reserve space: %v\n", err)
	}

	var buffer byte
	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return nil, fmt.Errorf("failed to read Attributes streamCodingType: %v\n", err)
	}
	subPath.IsRepeatSubPath = buffer&0x01 != 0

	// Skip 1-byte reserve space
	if _, err := file.Seek(1, io.SeekCurrent); err != nil {
		return nil, fmt.Errorf("failed to seek past reserve space: %v\n", err)
	}

	if err := binary.Read(file, binary.BigEndian, &subPath.NumberOfSubPlayItems); err != nil {
		return nil, fmt.Errorf("failed to read stream info NumberOfSubPlayItems: %v\n", err)
	}

	// Create the container of SubPlayItems
	subPath.SubPlayItems = make([]*SubPlayItem, subPath.NumberOfSubPlayItems)
	for i := uint8(0); i < uint8(len(subPath.SubPlayItems)); i++ {
		if subPath.SubPlayItems[i], err = ReadSubPlayItem(file); err != nil {
			return nil, fmt.Errorf("failed to read SubPlayItem: %v\n", err)
		}
	}

	// Skip to the end
	if _, err = file.Seek(end, io.SeekStart); err != nil {
		return nil, fmt.Errorf("failed to seek end offset: %v\n", err)
	}

	return subPath, nil
}

func (subPath *SubPath) Print() {
	PadPrintf(4, "Length: %d\n", subPath.Length)
	PadPrintf(4, "SubPathType: %d\n", subPath.SubPathType)
	PadPrintf(4, "IsRepeatSubPath: %v\n", subPath.IsRepeatSubPath)
	PadPrintf(4, "NumberOfSubPlayItems: %d\n", subPath.NumberOfSubPlayItems)
	PadPrintf(4, "SubPlayItems:\n")
	for i := uint8(0); i < uint8(len(subPath.SubPlayItems)); i++ {
		PadPrintf(6, "Angle [%d]:\n", i)
		subPath.SubPlayItems[i].Print()
	}
}
