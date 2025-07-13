package mobj

import (
	"encoding/binary"
	"fmt"
	"io"
)

// MOBJHeader represents the 40 byte header of an MOBJ file
type MOBJHeader struct {
	TypeIndicator [4]byte // "MOBJ
	VersionNumber [4]byte // "0100" or "0200"
	MovieObjects  *OffsetsUint32
	Extensions    *OffsetsUint32
}

type OffsetsUint32 struct {
	Start,
	Stop int64
}

func ReadMOBJHeader(file io.ReadSeeker) (header *MOBJHeader, err error) {
	header = &MOBJHeader{}

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

	if buffer == 0 {
		header.MovieObjects = &OffsetsUint32{Start: 40, Stop: eof}
		header.Extensions = &OffsetsUint32{Start: 0, Stop: 0}
	} else {
		header.MovieObjects = &OffsetsUint32{Start: 40, Stop: int64(buffer)}
		header.Extensions = &OffsetsUint32{Start: header.MovieObjects.Stop, Stop: eof}
	}

	return header, nil
}

func (offsets *OffsetsUint32) String() string {
	return fmt.Sprintf("{Start: %d, Stop: %d}", offsets.Start, offsets.Stop)
}

// String returns a string representation of the MPLSHeader.
func (header *MOBJHeader) String() string {
	return fmt.Sprintf(
		"Header{"+
			"Type: %s, "+
			"Version: %s, "+
			"Offset MovieObjects: %s, "+
			"Offset Extensions: %s, "+
			"}",
		string(header.TypeIndicator[:]),
		string(header.VersionNumber[:]),
		header.MovieObjects,
		header.Extensions,
	)
}
