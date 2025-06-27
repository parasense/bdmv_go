package indx

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Title struct {
	ObjectType         uint8   // 2-bits 0b11000000 & 0xC0 >> 6
	AccesType          uint8   // 2-bits 0b00110000 & 0x30 >> 4
	PlaybackType       uint8   // 2-bits 0b11000000 & 0xC0 >> 6
	RefToMovieObjectID uint16  // 16-bits
	RefToBDJObjectID   [5]byte // 40-bits
}

func ReadTitle(file io.ReadSeeker) (*Title, error) {
	title := &Title{}

	var buffer byte
	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return nil, fmt.Errorf("failed to read buffer: %w", err)
	}
	title.ObjectType = (buffer & 0xC0) >> 6
	title.AccesType = (buffer & 0x30) >> 4

	// skip 3 byte reserve space
	if _, err := file.Seek(3, io.SeekCurrent); err != nil {
		return nil, fmt.Errorf("failed to seek past reserve space: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return nil, fmt.Errorf("failed to read buffer: %w", err)
	}
	title.PlaybackType = (buffer & 0xC0) >> 6

	// skip 1 byte reserve space
	if _, err := file.Seek(1, io.SeekCurrent); err != nil {
		return nil, fmt.Errorf("failed to seek past reserve space: %w", err)
	}

	switch title.ObjectType {
	case 1: // Movie Object (16-bits + 32-bits)
		if err := binary.Read(file, binary.BigEndian, &title.RefToMovieObjectID); err != nil {
			return nil, fmt.Errorf("failed to read RefToMovieObjectID: %w", err)
		}
		// skip 4 byte reserve space
		if _, err := file.Seek(4, io.SeekCurrent); err != nil {
			return nil, fmt.Errorf("failed to seek past reserve space: %w", err)
		}
	case 2: // BDJ Object (40-bits + 8-bits)
		if err := binary.Read(file, binary.BigEndian, &title.RefToBDJObjectID); err != nil {
			return nil, fmt.Errorf("failed to read RefToBDJObjectID: %w", err)
		}
		// skip 1 byte reserve space
		if _, err := file.Seek(1, io.SeekCurrent); err != nil {
			return nil, fmt.Errorf("failed to seek past reserve space: %w", err)
		}
	default:
		// skip 6 bytes of unknown data
		if _, err := file.Seek(6, io.SeekCurrent); err != nil {
			return nil, fmt.Errorf("failed to seek past reserve space: %w", err)
		}
		return nil, fmt.Errorf("unsupported ObjectType: %d", title.ObjectType)
	}

	return title, nil
}

func (t *Title) String() string {
	return fmt.Sprintf(
		"Title: {ObjectType: %d, AccesType: %d, PlaybackType: %d, RefToMovieObjectID: %d, RefToBDJObjectID: %x}",
		t.ObjectType,
		t.AccesType,
		t.PlaybackType,
		t.RefToMovieObjectID,
		t.RefToBDJObjectID,
	)
}
