package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

/*
DEV NOTE:

	Stream Entry & Attribute structures need work
	* video streams don't need audio fields.
	* audio streams don't need video fields.
	* each Entry & attribute is generic, and they could be more specific.

	Alignment issues still potentially lurk.

	Code has been cleaned
	* all printing code now lives near the related data struct.
	* printing code now calls other printing code all the way down the leaf structures.

	TODO:
	* Extension data stuff (in-progress)
	  - Subpaths (in progress)
	  - 3d STNs (not started)
	  - Static metadata (not started)
	  - PIP metadata (not started)
	* Runtime assertions & deffensive coding
	* encoding/binary optimizations
	  - Try to avoid any reflection code paths.
	  - Look at "uvariant" functions
	* AUDIT the bit/byte-wise operatiosn for shift errors.
	  - Found a mistake when parsing multi-bit most-significan bits (lack of right shifting)
	  - For example `FOO = buffer & 0xF0` ==> `FOO = (buffer & 0xF0) >> 4`
*/

// Why this is not part of the standard library boggles the mind.
// Seek gives your current position when you use it. So if you seek zero
// ahead, Seek() returns your current position.
func ftell(file io.ReadSeeker) (int64, error) {
	return file.Seek(0, io.SeekCurrent)
}

// Parse45KhzTimestamp converts a timestamp in 45kHz units to a duration
// returns a quantity of time in seconds.
func Parse45KhzTimestamp(timestamp uint32) time.Duration {
	return time.Duration(uint32(timestamp) / 45000)
}

// returns a quantity of time in seconds.
func Convert45KhzTimeToSeconds(timestamp uint32) uint32 {
	return uint32(timestamp / 45000)
}

func PadPrintf(indent int, format string, args ...any) {
	fmt.Printf(strings.Repeat(" ", indent)+format, args...)
}

func PadPrintln(indent int, args ...any) {
	fmt.Print(strings.Repeat(" ", indent))
	fmt.Println(args...)
}

// ParseMPLS parses an MPLS file and returns the playlist details
func ParseMPLS(filePath string) (
	header *MPLSHeader,
	appinfo *AppInfo,
	playlist *PlayList,
	chapterMarks *PlaylistMarks,
	extensiondata *Extensions,
	err error,
) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Header
	if header, err = ReadMPLSHeader(file); err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to read header: %w", err)
	}

	// AppInfo
	if appinfo, err = ReadAppInfo(file, &header.AppInfo); err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to read appinfo: %w", err)
	}

	// Playlist
	if playlist, err = ReadPlayList(file, &header.Playlist); err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to read PlayList: %w", err)
	}

	// Marks
	if chapterMarks, err = ReadMarks(file, &header.Marks); err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to read Chapter Marks: %w", err)
	}

	// Extensions
	if header.Extensions.Start != 0 && header.Extensions.Stop != 0 {
		if extensiondata, err = ReadExtensions(file, &header.Extensions); err != nil {
			return nil, nil, nil, nil, nil, fmt.Errorf("failed to read Extension Data: %w", err)
		}
	}

	return header, appinfo, playlist, chapterMarks, extensiondata, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: mpls-parser <mpls-file>")
		os.Exit(1)
	}

	mplsPath := os.Args[1]
	header, appinfo, playlist, chapterMarks, extData, err := ParseMPLS(mplsPath)
	if err != nil {
		fmt.Errorf("Error parsing MPLS file: %w\n", err)
		os.Exit(1)
	}

	PadPrintf(0, "MPLS File: %s\n", mplsPath)
	header.Print()
	appinfo.Print()
	playlist.Print()
	chapterMarks.Print()
	if header.Extensions.Start != 0 && header.Extensions.Stop != 0 {
		extData.Print()
	}
}
