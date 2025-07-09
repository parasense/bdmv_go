package clpi

import (
	"encoding/binary"
	"fmt"
	"io"
)

// CLPIHeader represents the 40 byte header of a CLPI file
type CLPIHeader struct {
	TypeIndicator [4]byte // "HDMV"
	VersionNumber [4]byte // "0100" or "0200"
	ClipInfo      *OffsetsUint32
	SequenceInfo  *OffsetsUint32
	ProgramInfo   *OffsetsUint32
	CPI           *OffsetsUint32
	ClipMarks     *OffsetsUint32
	Extensions    *OffsetsUint32
}

// OffsetsUint32 represents the start and stop offsets of a section in the CLPI file.
// The CLPI file format uses 32-bit unsigned integers for (start) offsets, but Go's
// io.Seeker interface requires int64 for seeking.
// Therefore, we use int64 to represent the offsets, even though they are conceptually
// 32-bit unsigned integers. This avoids issues with the io.Seeker interface.
// The Start and Stop fields represent the start and stop offsets of a section in the MPLS
// If a section has no data, both Start and Stop will be 0.
type OffsetsUint32 struct {
	Start,
	Stop int64
}

func ReadCLPIHeader(file io.ReadSeeker) (header *CLPIHeader, err error) {
	header = &CLPIHeader{}

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

	header.ClipInfo = &OffsetsUint32{Start: 40, Stop: int64(buffer)}
	header.SequenceInfo = &OffsetsUint32{Start: header.ClipInfo.Stop}

	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return nil, fmt.Errorf("failed to read header.Playlist.start: %w", err)
	}

	header.SequenceInfo.Stop = int64(buffer)
	header.ProgramInfo = &OffsetsUint32{Start: header.SequenceInfo.Stop}

	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return nil, fmt.Errorf("failed to read header.Marks.start: %w", err)
	}

	header.ProgramInfo.Stop = int64(buffer)
	header.CPI = &OffsetsUint32{Start: header.ProgramInfo.Stop}

	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return nil, fmt.Errorf("failed to read header.Extensions.Start: %w", err)
	}

	header.CPI.Stop = int64(buffer)
	header.ClipMarks = &OffsetsUint32{Start: header.CPI.Stop}

	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return nil, fmt.Errorf("failed to read header.Extensions.Start: %w", err)
	}

	if buffer == 0 {
		header.ClipMarks.Stop = eof
		header.Extensions = &OffsetsUint32{Start: 0, Stop: 0}
	} else {
		header.ClipMarks.Stop = int64(buffer)
		header.Extensions = &OffsetsUint32{Start: header.ClipMarks.Stop, Stop: eof}
	}

	return header, nil
}

func (offsets *OffsetsUint32) String() string {
	return fmt.Sprintf("{Start: %d, Stop: %d}", offsets.Start, offsets.Stop)
}

// String returns a string representation of the CLPIHeader.
func (header *CLPIHeader) String() string {
	return fmt.Sprintf("CLPIHeader{\n"+
		"  TypeIndicator: %s\n"+
		"  VersionNumber: %s\n"+
		"  ClipInfo: %s\n"+
		"  SequenceInfo: %s\n"+
		"  ProgramInfo: %s\n"+
		"  CPI: %s\n"+
		"  ClipMark: %s\n"+
		"  Extensions: %s\n"+
		"}",
		header.TypeIndicator, header.VersionNumber,
		header.ClipInfo, header.SequenceInfo,
		header.ProgramInfo, header.CPI,
		header.ClipMarks, header.Extensions)
}
