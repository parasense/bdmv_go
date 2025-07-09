package clpi

import (
	"fmt"
	"io"
	"os"
	"strings"
)

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

func ParseCLPI(filePath string) (
	header *CLPIHeader,
	clipInfo *ClipInfo,
	sequenceInfo *SequenceInfo,
	programInfo *ProgramInfo,
	cpi *CPI,
	clipMarks *ClipMarks,
	extensiondata *Extensions,
	err error,
) {

	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Header
	if header, err = ReadCLPIHeader(file); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("failed to read header: %w", err)
	}

	// ClipInfo
	if clipInfo, err = ReadClipInfo(file, header.ClipInfo); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("failed to read clipinfo: %w", err)
	}

	// SequenceInfo
	if sequenceInfo, err = ReadSequenceInfo(file, header.SequenceInfo); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("failed to read sequence info: %w", err)
	}

	// ProgramInfo
	if programInfo, err = ReadProgramInfo(file, header.ProgramInfo); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("failed to read program info: %w", err)
	}

	// CPI
	if cpi, err = ReadCPI(file, header.CPI); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("failed to read CPI: %w", err)
	}

	// Clip Marks
	if clipMarks, err = ReadClipMarks(file, header.ClipMarks); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("failed to read ClipMarks: %w", err)
	}

	// Extensions
	if header.Extensions.Start != 0 && header.Extensions.Stop != 0 {
		if extensiondata, err = ReadExtensions(file, header.Extensions); err != nil {
			return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("failed to read Extension Data: %w", err)
		}
	}

	return header, clipInfo, sequenceInfo, programInfo, cpi, clipMarks, extensiondata, nil
}
