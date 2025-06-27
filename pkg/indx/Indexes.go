package indx

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Indexes struct {
	Length             uint32
	FirstPlaybackTitle *Title
	TopMenuTitle       *Title
	NumberOfTitles     uint16
	Titles             []*Title
}

func ReadIndexes(file io.ReadSeeker, offsets *OffsetsUint32) (indexes *Indexes, err error) {
	indexes = &Indexes{}

	// Jump to start address
	if _, err := file.Seek(offsets.Start, io.SeekStart); err != nil {
		return nil, fmt.Errorf("failed to seek to start address: %w\n", err)
	}

	if err := binary.Read(file, binary.BigEndian, &indexes.Length); err != nil {
		return nil, fmt.Errorf("failed to read indexes.Length: %w", err)
	}

	if indexes.FirstPlaybackTitle, err = ReadTitle(file); err != nil {
		return nil, fmt.Errorf("failed to read FirstPlaybackTitle: %w", err)
	}

	if indexes.TopMenuTitle, err = ReadTitle(file); err != nil {
		return nil, fmt.Errorf("failed to read TopMenuTitle: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &indexes.NumberOfTitles); err != nil {
		return nil, fmt.Errorf("failed to read NumberOfTitles: %w", err)
	}

	indexes.Titles = make([]*Title, indexes.NumberOfTitles)
	for i := range indexes.Titles {
		if indexes.Titles[i], err = ReadTitle(file); err != nil {
			return nil, fmt.Errorf("failed to read Title[%d]: %w", i, err)
		}
	}

	return indexes, nil
}

func (indexes *Indexes) String() string {
	return fmt.Sprintf("Indexes{Length: %d, FirstPlaybackTitle: %s, TopMenuTitle: %s, NumberOfTitles: %d, Titles: %s}",
		indexes.Length,
		indexes.FirstPlaybackTitle,
		indexes.TopMenuTitle,
		indexes.NumberOfTitles,
		indexes.Titles)
}
