package mpls

import (
	"encoding/binary"
	"fmt"
	"io"
)

// PlayList represents the main playlist structure
// It contains the length of the playlist, number of play items, number of sub paths,
// and slices for play items and sub paths.
type PlayList struct {
	Length            uint32
	NumberOfPlayItems uint16
	NumberOfSubPaths  uint16
	PlayItems         []*PlayItem
	SubPaths          []*SubPath
}

// ReadPlayList reads a playlist from the provided file at the specified offsets
// and returns a PlayList object.
// It expects the file (io.ReadSeeker) to be positioned at the start of the PlayList structure.
// The PlayList consists of a length, number of play items, number of sub paths,
// followed by the play items and sub paths themselves.
// It returns a pointer to the PlayList and an error if any occurs during reading.
func ReadPlayList(file io.ReadSeeker, offsets *OffsetsUint32) (playlist *PlayList, err error) {
	playlist = &PlayList{}

	// Jump to start address
	if _, err := file.Seek(offsets.Start, io.SeekStart); err != nil {
		return nil, fmt.Errorf("failed to seek to start address: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &playlist.Length); err != nil {
		return nil, fmt.Errorf("failed to read playlist.Length: %w", err)
	}

	// Skip past reserve space between Length and NumberOfPlayItems
	if _, err := file.Seek(2, io.SeekCurrent); err != nil {
		return nil, fmt.Errorf("failed to seek past playlist reserve space: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &playlist.NumberOfPlayItems); err != nil {
		return nil, fmt.Errorf("failed to read playlist.NumberOfPlayItems: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &playlist.NumberOfSubPaths); err != nil {
		return nil, fmt.Errorf("failed to read playlist.NumberOfSubPaths: %w", err)
	}

	playlist.PlayItems = make([]*PlayItem, playlist.NumberOfPlayItems)
	for i := range uint16(playlist.NumberOfPlayItems) {
		if playlist.PlayItems[i], err = ReadPlayItem(file); err != nil {
			return nil, fmt.Errorf("failed to read PlayListItem: %w", err)
		}
		playlist.PlayItems[i].Assert()
	}

	playlist.SubPaths = make([]*SubPath, playlist.NumberOfSubPaths)
	for i := range uint16(playlist.NumberOfSubPaths) {
		if playlist.SubPaths[i], err = ReadSubPath(file); err != nil {
			return nil, fmt.Errorf("failed to read SubPath: %w", err)
		}
		//playlist.SubPaths[i].Assert()
	}

	return playlist, nil
}

func (playlist *PlayList) String() string {
	return fmt.Sprintf("PlayList{Length: %d, NumberOfPlayItems: %d, NumberOfSubPaths: %d, PlayItems: %v, SubPaths: %v}",
		playlist.Length, playlist.NumberOfPlayItems, playlist.NumberOfSubPaths, playlist.PlayItems, playlist.SubPaths)
}
