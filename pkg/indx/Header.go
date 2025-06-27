package indx

import (
	"encoding/binary"
	"fmt"
	"io"
)

// MPLSHeader represents the 40 byte header of an MPLS file
type INDXHeader struct {
	TypeIndicator [4]byte // "INDX"
	VersionNumber [4]byte // "0100" or "0200"
	AppInfo       *OffsetsUint32
	Indexes       *OffsetsUint32
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

func ReadINDXHeader(file io.ReadSeeker) (header *INDXHeader, err error) {
	header = &INDXHeader{}

	var eof int64
	if eof, err = file.Seek(0, io.SeekEnd); err != nil {
		return nil, err
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	if err := binary.Read(file, binary.BigEndian, &header.TypeIndicator); err != nil {
		return nil, err
	}

	if err := binary.Read(file, binary.BigEndian, &header.VersionNumber); err != nil {
		return nil, err
	}

	var buffer uint32
	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return nil, err
	}

	header.AppInfo = &OffsetsUint32{Start: 40, Stop: int64(buffer)}
	header.Indexes = &OffsetsUint32{Start: header.AppInfo.Stop}

	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return nil, err
	}

	if buffer == 0 {
		header.Indexes.Stop = eof
		header.Extensions = &OffsetsUint32{Start: 0, Stop: 0}
	} else {
		header.Indexes.Stop = int64(buffer)
		header.Extensions = &OffsetsUint32{Start: header.Indexes.Stop, Stop: eof}
	}

	return header, nil
}

func (header *INDXHeader) String() string {
	return "{TypeIndicator: " + string(header.TypeIndicator[:]) +
		", VersionNumber: " + string(header.VersionNumber[:]) +
		", AppInfo: " + header.AppInfo.String() +
		", Indexes: " + header.Indexes.String() +
		", Extensions: " + header.Extensions.String() + "}"
}
func (offsets *OffsetsUint32) String() string {
	return "{Start: " + fmt.Sprint(offsets.Start) +
		", Stop: " + fmt.Sprint(offsets.Stop) + "}"
}
