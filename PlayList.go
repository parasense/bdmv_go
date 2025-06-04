package main

import (
	"encoding/binary"
	"fmt"
	"io"
)

// PlayList represents the main playlist structure
type PlayList struct {
	Length            uint32
	NumberOfPlayItems uint16
	NumberOfSubPaths  uint16
	PlayItems         []*PlayItem
	SubPaths          []*SubPath
}

func ReadPlayList(file io.ReadSeeker) (playlist *PlayList, err error) {
	playlist = &PlayList{}
	if err := binary.Read(file, binary.BigEndian, &playlist.Length); err != nil {
		return nil, fmt.Errorf("failed to read playlist.Length: %v", err)
	}
	// Skip past reserve space between Length and NumberOfPlayItems
	if _, err := file.Seek(2, io.SeekCurrent); err != nil {
		return nil, fmt.Errorf("failed to seek past playlist reserve space: %v", err)
	}
	if err := binary.Read(file, binary.BigEndian, &playlist.NumberOfPlayItems); err != nil {
		return nil, fmt.Errorf("failed to read playlist.NumberOfPlayItems: %v", err)
	}
	if err := binary.Read(file, binary.BigEndian, &playlist.NumberOfSubPaths); err != nil {
		return nil, fmt.Errorf("failed to read playlist.NumberOfSubPaths: %v", err)
	}

	playlist.PlayItems = make([]*PlayItem, playlist.NumberOfPlayItems)
	for i := range uint16(playlist.NumberOfPlayItems) {
		if playlist.PlayItems[i], err = ReadPlayItem(file); err != nil {
			return nil, fmt.Errorf("failed to read PlayListItem: %v\n", err)
		}
		playlist.PlayItems[i].Assert()
	}

	playlist.SubPaths = make([]*SubPath, playlist.NumberOfSubPaths)
	for i := range uint16(playlist.NumberOfSubPaths) {
		if playlist.SubPaths[i], err = ReadSubPath(file); err != nil {
			return nil, fmt.Errorf("failed to read SubPath: %v\n", err)
		}
		//playlist.SubPaths[i].Assert()
	}

	return playlist, nil
}

func (playlist *PlayList) Print() {
	PadPrintln(0, "PlayList:")
	PadPrintf(2, "Length: %d\n", playlist.Length)
	PadPrintf(2, "NumberOfPlayItems: %d\n", playlist.NumberOfPlayItems)
	PadPrintf(2, "NumberOfSubPaths: %d\n", playlist.NumberOfSubPaths)
	PadPrintln(2)

	var totalDuration uint32
	for i, playItem := range playlist.PlayItems {
		inTime := Convert45KhzTimeToSeconds(playItem.INTime)
		outTime := Convert45KhzTimeToSeconds(playItem.OUTTime)
		duration := outTime - inTime
		totalDuration += duration

		PadPrintf(2, "PlayItem [%d]:\n", i)
		playItem.Print()
		PadPrintln(2, "---")

	}

	for i, subPath := range playlist.SubPaths {
		PadPrintf(2, "SubPath [%d]:\n", i)
		subPath.Print()
		PadPrintln(2, "---")
	}
	PadPrintln(0, "---")
}
