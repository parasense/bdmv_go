package mpls

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

// ReadSubPath reads a SubPath from the provided io.ReadSeeker.
// It expects the file (io.ReadSeeker) to be positioned at the start of the
// SubPath structure.
func ReadSubPath(file io.ReadSeeker) (subPath *SubPath, err error) {
	subPath = &SubPath{}

	if err := binary.Read(file, binary.BigEndian, &subPath.Length); err != nil {
		return nil, fmt.Errorf("failed to read stream info length: %w", err)
	}

	end, err := CalculateEndOffset(file, subPath.Length)
	if err != nil {
		return nil, fmt.Errorf("failed calling CalculateEndOffset(): %w", err)
	}

	// Skip 1-byte reserve space
	if _, err := file.Seek(1, io.SeekCurrent); err != nil {
		return nil, fmt.Errorf("failed to seek past reserve space: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &subPath.SubPathType); err != nil {
		return nil, fmt.Errorf("failed to read stream info SubPathType: %w", err)
	}

	// Skip 1-byte reserve space
	if _, err := file.Seek(1, io.SeekCurrent); err != nil {
		return nil, fmt.Errorf("failed to seek past reserve space: %w", err)
	}

	var buffer byte
	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return nil, fmt.Errorf("failed to read Attributes streamCodingType: %w", err)
	}
	subPath.IsRepeatSubPath = buffer&0x01 != 0

	// Skip 1-byte reserve space
	if _, err := file.Seek(1, io.SeekCurrent); err != nil {
		return nil, fmt.Errorf("failed to seek past reserve space: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &subPath.NumberOfSubPlayItems); err != nil {
		return nil, fmt.Errorf("failed to read stream info NumberOfSubPlayItems: %w", err)
	}

	// Create the container of SubPlayItems
	subPath.SubPlayItems = make([]*SubPlayItem, subPath.NumberOfSubPlayItems)
	for i := range subPath.SubPlayItems {

		if subPath.SubPlayItems[i], err = ReadSubPlayItem(file); err != nil {
			return nil, fmt.Errorf("failed to read SubPlayItem: %w", err)
		}

		// 1-byte reserve space
		if _, err := file.Seek(1, io.SeekCurrent); err != nil {
			return nil, fmt.Errorf("failed to seek past reserve space: %w", err)
		}

	}

	// Skip to the end
	if _, err := file.Seek(end, io.SeekStart); err != nil {
		return nil, fmt.Errorf("failed to seek end offset: %w", err)
	}

	return subPath, nil
}

func (subPath *SubPath) String() string {
	return fmt.Sprintf("SubPath:\n"+
		"  Length: %d\n"+
		"  SubPathType: %d\n"+
		"  IsRepeatSubPath: %t\n"+
		"  NumberOfSubPlayItems: %d\n"+
		"  SubPlayItems: %v",
		subPath.Length,
		subPath.SubPathType,
		subPath.IsRepeatSubPath,
		subPath.NumberOfSubPlayItems,
		subPath.SubPlayItems)
}
