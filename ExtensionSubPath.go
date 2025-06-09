package main

import (
	"encoding/binary"
	"io"
)

type ExtensionSubPath struct {
	Length   uint32
	Count    uint16
	SubPaths []*SubPath
}

func (subPathExtension *ExtensionSubPath) Read(file io.ReadSeeker) (err error) {

	if err := binary.Read(file, binary.BigEndian, &subPathExtension.Length); err != nil {
		return err
	}

	if err := binary.Read(file, binary.BigEndian, &subPathExtension.Count); err != nil {
		return err
	}

	subPathExtension.SubPaths = make([]*SubPath, subPathExtension.Count)
	for i := range subPathExtension.SubPaths {
		subPathExtension.SubPaths[i], _ = ReadSubPath(file)
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
