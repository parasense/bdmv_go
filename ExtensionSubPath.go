package main

import (
	"encoding/binary"
	"fmt"
	"io"
)

type ExtensionSubPath struct {
	Length   uint32
	Count    uint16
	SubPaths []*SubPath
}

func (extensionSubPath *ExtensionSubPath) Read(file io.ReadSeeker) (err error) {

	if err := binary.Read(file, binary.BigEndian, &extensionSubPath.Length); err != nil {
		return fmt.Errorf("failed to read extensionSubPath.Length: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &extensionSubPath.Count); err != nil {
		return fmt.Errorf("failed to read extensionSubPath.Count: %w", err)
	}

	extensionSubPath.SubPaths = make([]*SubPath, extensionSubPath.Count)
	for i := range extensionSubPath.SubPaths {
		if extensionSubPath.SubPaths[i], err = ReadSubPath(file); err != nil {
			return fmt.Errorf("failed calling ReadSubPath() in ExtensionSubPath.Read(): %w", err)
		}
	}

	return nil
}

func (subPathExtension *ExtensionSubPath) Print() {
	PadPrintln(4, "SubPathExtension")
	PadPrintf(6, "Length: %d\n", subPathExtension.Length)
	PadPrintf(6, "Count: %d\n", subPathExtension.Count)
	for i, subPath := range subPathExtension.SubPaths {
		PadPrintf(6, "SubPath [%d]:\n", i)
		subPath.Print()
	}
}
