package mpls

import (
	"encoding/binary"
	"fmt"
	"io"
)

// MPLSHeader represents the 40 byte header of an MPLS file
type MPLSHeader struct {
	TypeIndicator [4]byte // "MPLS"
	VersionNumber [4]byte // "0100" or "0200"
	AppInfo       *OffsetsUint32
	Playlist      *OffsetsUint32
	Marks         *OffsetsUint32
	Extensions    *OffsetsUint32
}

// OffsetsUint32 represents the start and stop offsets of a section in the MPLS file.
// The MPLS file format uses 32-bit unsigned integers for (start) offsets, but Go's
// io.Seeker interface requires int64 for seeking.
// Therefore, we use int64 to represent the offsets, even though they are conceptually
// 32-bit unsigned integers. This avoids issues with the io.Seeker interface.
// The Start and Stop fields represent the start and stop offsets of a section in the MPLS
// If a section has no data, both Start and Stop will be 0.
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

	header.AppInfo = &OffsetsUint32{Start: 40, Stop: int64(buffer)}
	header.Playlist = &OffsetsUint32{Start: header.AppInfo.Stop}

	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return nil, fmt.Errorf("failed to read header.Marks.start: %w", err)
	}

	header.Playlist.Stop = int64(buffer)
	header.Marks = &OffsetsUint32{Start: header.Playlist.Stop}

	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return nil, fmt.Errorf("failed to read header.Extensions.Start: %w", err)
	}

	if buffer == 0 {
		header.Marks.Stop = eof
		header.Extensions = &OffsetsUint32{Start: 0, Stop: 0}
	} else {
		header.Marks.Stop = int64(buffer)
		header.Extensions = &OffsetsUint32{Start: header.Marks.Stop, Stop: eof}

	}

	return header, nil
}

func (offsets *OffsetsUint32) String() string {
	return fmt.Sprintf("{Start: %d, Stop: %d}", offsets.Start, offsets.Stop)
}

// String returns a string representation of the MPLSHeader.
func (header *MPLSHeader) String() string {
	return fmt.Sprintf("Header: \n"+
		"  Type: %s, \n"+
		"  Version: %s, \n"+
		"  Offset AppInfo: %s, \n"+
		"  Offset PlayList: %s, \n"+
		"  Offset Marks: %s, \n"+
		"  Offset Extensions: %s,\n"+
		string(header.TypeIndicator[:]),
		string(header.VersionNumber[:]),
		header.AppInfo,
		header.Playlist,
		header.Marks,
		header.Extensions,
	)
}
