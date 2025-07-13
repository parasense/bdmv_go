package mobj

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Why this is not part of the standard library boggles the mind.
// Seek gives your current position when you use it. So if you seek zero
// ahead, Seek() returns your current position.
func ftell(file io.ReadSeeker) (int64, error) {
	return file.Seek(0, io.SeekCurrent)
}

func PadPrintf(indent int, format string, args ...any) {
	fmt.Printf(strings.Repeat(" ", indent)+format, args...)
}

func PadPrintln(indent int, args ...any) {
	fmt.Print(strings.Repeat(" ", indent))
	fmt.Println(args...)
}

func CalculateEndOffset[U uint8 | uint16 | uint32](file io.ReadSeeker, length U) (int64, error) {
	currentPos, err := file.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, fmt.Errorf("failed to get current position: %w", err)
	}
	return currentPos + int64(length), nil
}

// ParseMPLS parses an MPLS file and returns the playlist details
func ParseMOBJ(filePath string) (
	header *MOBJHeader,
	movieObjects *MovieObjects,
	extensiondata *Extensions,
	err error,
) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Header
	if header, err = ReadMOBJHeader(file); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to read header: %w", err)
	}

	// MovieObjects
	if movieObjects, err = ReadMovieObjects(file, header.MovieObjects); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to read appinfo: %w", err)
	}

	// Extensions
	if header.Extensions.Start != 0 && header.Extensions.Stop != 0 {
		if extensiondata, err = ReadExtensions(file, header.Extensions); err != nil {
			return nil, nil, nil, fmt.Errorf("failed to read Extension Data: %w", err)
		}
	}

	return header, movieObjects, extensiondata, nil
}
