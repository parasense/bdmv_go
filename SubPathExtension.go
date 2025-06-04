package main

import (
	"encoding/binary"
	"fmt"
	"io"
)

type SubPathExtension struct {
	Length   uint32
	Count    uint16
	SubPaths []*SubPath
}

// XXX - fix error handling
func (subPathExtension *SubPathExtension) Read(file io.ReadSeeker) (err error) {

	if err := binary.Read(file, binary.BigEndian, &subPathExtension.Length); err != nil {
		return fmt.Errorf("failed to read markEntry type: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &subPathExtension.Count); err != nil {
		return fmt.Errorf("failed to read markEntry type: %w", err)
	}

	subPathExtension.SubPaths = make([]*SubPath, subPathExtension.Count)
	for i := range subPathExtension.SubPaths {
		subPathExtension.SubPaths[i], _ = ReadSubPath(file)
	}

	return nil
}

func (subPathExtension SubPathExtension) Print() {
	fmt.Println("SubPathExtension")
	fmt.Printf("  Length: %d\n", subPathExtension.Length)
	fmt.Printf("  Count: %d\n", subPathExtension.Count)
	for i, subPath := range subPathExtension.SubPaths {
		fmt.Printf("  subPath [%d]:\n", i)
		subPath.Print()
	}
}
