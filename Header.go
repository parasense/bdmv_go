package main

import (
	"encoding/binary"
	"fmt"
	"io"
)

// MPLSHeader represents the 40 byte header of an MPLS file
type MPLSHeader struct {
	TypeIndicator [4]byte // "MPLS"
	VersionNumber [4]byte // "0100" or "0200"
	AppInfo       OffsetsUint32
	Playlist      OffsetsUint32
	Marks         OffsetsUint32
	Extensions    OffsetsUint32
}

// XXX - consider changing this to int64
// Because they are used in io.Seeker which only uses int64
type OffsetsUint32 struct {
	Start,
	Stop int64
}

func ReadMPLSHeader(file io.ReadSeeker) (header *MPLSHeader, err error) {
	header = &MPLSHeader{}

	var eof int64
	if eof, err = file.Seek(0, io.SeekEnd); err != nil {
		return nil, fmt.Errorf("failed to seek to file end address: %w", err)
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("failed to seek to file start address: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &header.TypeIndicator); err != nil {
		return nil, fmt.Errorf("failed to read header.TypeIndicator: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &header.VersionNumber); err != nil {
		return nil, fmt.Errorf("failed to read header.VersionNumber: %w", err)
	}

	var buffer uint32
	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return nil, fmt.Errorf("failed to read header.Playlist.start: %w", err)
	}

	header.AppInfo.Start = 40
	header.AppInfo.Stop = int64(buffer)
	header.Playlist.Start = header.AppInfo.Stop

	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return nil, fmt.Errorf("failed to read header.Marks.start: %w", err)
	}

	header.Playlist.Stop = int64(buffer)
	header.Marks.Start = header.Playlist.Stop

	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return nil, fmt.Errorf("failed to read header.Extensions.Start: %w", err)
	}

	if buffer == 0 {
		header.Marks.Stop = eof
		header.Extensions.Start = 0
		header.Extensions.Stop = 0
	} else {
		header.Marks.Stop = int64(buffer)
		header.Extensions.Start = header.Marks.Stop
		header.Extensions.Stop = eof
	}

	return header, nil
}

func (header *MPLSHeader) Print() {
	fmt.Println("Header:")
	PadPrintf(2, "Type: %s\n", string(header.TypeIndicator[:]))
	PadPrintf(2, "Version: %s\n", string(header.VersionNumber[:]))
	PadPrintf(2, "Offset: AppInfo: [%d:%d]\n", header.AppInfo.Start, header.AppInfo.Stop)
	PadPrintf(2, "Offset: PlayList: [%d:%d]\n", header.Playlist.Start, header.Playlist.Stop)
	PadPrintf(2, "Offset: Marks: [%d:%d]\n", header.Marks.Start, header.Marks.Stop)
	PadPrintf(2, "Offset: Extensions: [%d:%d]\n", header.Extensions.Start, header.Extensions.Stop)
	PadPrintln(2, "---")
}
