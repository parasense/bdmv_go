package main

import (
	"encoding/binary"
	"fmt"
	"io"
)

// MPLSHeader represents the 40 byte header of an MPLS file
type MPLSHeader struct {
	TypeIndicator            [4]byte // "MPLS"
	VersionNumber            [4]byte // "0100" or "0200"
	AppInfoStartAddress      uint32  // starts after the header, ends before Playlist
	PlayListStartAddress     uint32
	PlayListMarkStartAddress uint32
	ExtensionStartAddress    uint32
}

func ReadMPLSHeader(file io.ReadSeeker) (*MPLSHeader, error) {
	header := &MPLSHeader{}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("failed to seek to header start address: %v", err)
	}
	if err := binary.Read(file, binary.BigEndian, &header.TypeIndicator); err != nil {
		return nil, fmt.Errorf("failed to read header.TypeIndicator: %v", err)
	}
	if err := binary.Read(file, binary.BigEndian, &header.VersionNumber); err != nil {
		return nil, fmt.Errorf("failed to read header.VersionNumber: %v", err)
	}
	header.AppInfoStartAddress = 40
	if err := binary.Read(file, binary.BigEndian, &header.PlayListStartAddress); err != nil {
		return nil, fmt.Errorf("failed to read header.PlayListStartAddress: %v", err)
	}
	if err := binary.Read(file, binary.BigEndian, &header.PlayListMarkStartAddress); err != nil {
		return nil, fmt.Errorf("failed to read header.PlayListMarkStartAddress: %v", err)
	}
	if err := binary.Read(file, binary.BigEndian, &header.ExtensionStartAddress); err != nil {
		return nil, fmt.Errorf("failed to read header.ExtensionStartAddress: %v", err)
	}
	return header, nil
}

func (header *MPLSHeader) Print() {
	fmt.Println("Header:")
	PadPrintf(2, "Type: %s\n", string(header.TypeIndicator[:]))
	PadPrintf(2, "Version: %s\n", string(header.VersionNumber[:]))
	PadPrintf(2, "Offset AppInfo: %v\n", header.AppInfoStartAddress)
	PadPrintf(2, "Offset PlayList: %v\n", header.PlayListStartAddress)
	PadPrintf(2, "Offset PlayListMark: %v\n", header.PlayListMarkStartAddress)
	PadPrintf(2, "Offset ExtensionData: %v\n", header.ExtensionStartAddress)
	PadPrintln(2, "---")
}
