package mobj

import (
	"encoding/binary"
	"fmt"
	"io"
)

type MovieObjects struct {
	Length               uint32
	NumberOfMovieObjects uint16
	MovieObjects         []*MovieObject
}

type MovieObject struct {
	ResumeIntentionFlag        bool // 0b10000000
	MenuCallMask               bool // 0b01000000
	TitleSearchMask            bool // 0b00100000
	NumberOfNavigationCommands uint16
	NavigationCommands         []*NavigationCommand
}

type NavigationCommand struct {
	OperandCount           uint8 // 0b11100000
	CommandGroup           uint8 // 0b00011000
	CommandSubGroup        uint8 // 0b00000111
	ImmediateValueFlagDest bool  // 0b10000000
	ImmediateValueFlagSrc  bool  // 0b01000000
	BranchOption           uint8 // 0b00001111
	_                      uint8 // 0b11110000
	CompareOption          uint8 // 0b00001111
	_                      uint8 // 0b11100000
	SetOption              uint8 // 0b00011111
	Destination            uint32
	Source                 uint32
}

func ReadMovieObjects(file io.ReadSeeker, offsets *OffsetsUint32) (mobjs *MovieObjects, err error) {

	// Jump to start address
	if _, err := file.Seek(offsets.Start, io.SeekStart); err != nil {
		return nil, fmt.Errorf("failed to seek to start address: %w", err)
	}

	mobjs = &MovieObjects{}

	if err := binary.Read(file, binary.BigEndian, &mobjs.Length); err != nil {
		return nil, fmt.Errorf("failed to read play item length: %w", err)
	}

	// skip 4-bytes reserve space
	if _, err := file.Seek(4, io.SeekCurrent); err != nil {
		return nil, fmt.Errorf("failed to seek past reserve space: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &mobjs.NumberOfMovieObjects); err != nil {
		return nil, fmt.Errorf("failed to read play item length: %w", err)
	}

	mobjs.MovieObjects = make([]*MovieObject, mobjs.NumberOfMovieObjects)

	for i := range mobjs.MovieObjects {
		if mobjs.MovieObjects[i], err = ReadMovieObject(file); err != nil {
			return nil, err
		}
	}

	return mobjs, err
}

func ReadMovieObject(file io.ReadSeeker) (mobj *MovieObject, err error) {
	mobj = &MovieObject{}

	var buffer uint8
	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return nil, fmt.Errorf("failed to read play item length: %w", err)
	}
	mobj.ResumeIntentionFlag = (buffer&0x80 != 0)
	mobj.MenuCallMask = (buffer&040 != 0)
	mobj.TitleSearchMask = (buffer&0x20 != 0)

	// skip 1-byte reserve space
	if _, err := file.Seek(1, io.SeekCurrent); err != nil {
		return nil, fmt.Errorf("failed to seek past reserve space: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &mobj.NumberOfNavigationCommands); err != nil {
		return nil, fmt.Errorf("failed to read play item length: %w", err)
	}

	mobj.NavigationCommands = make([]*NavigationCommand, mobj.NumberOfNavigationCommands)
	for i := range mobj.NavigationCommands {
		if mobj.NavigationCommands[i], err = ReadNavCmd(file); err != nil {
			return nil, err
		}
	}

	return mobj, err
}

func ReadNavCmd(file io.ReadSeeker) (nav *NavigationCommand, err error) {
	nav = &NavigationCommand{}

	var buffer uint8
	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return nil, fmt.Errorf("failed to read buffer: %w", err)
	}
	nav.OperandCount = (buffer & 0xE0) >> 5 // 0b11100000
	nav.CommandGroup = (buffer & 0x18) >> 3 // 0b00011000
	nav.CommandSubGroup = (buffer & 0x03)   // 0b00000011

	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return nil, fmt.Errorf("failed to read buffer: %w", err)
	}
	nav.ImmediateValueFlagDest = (buffer&0x80 != 0) // 0b10000000
	nav.ImmediateValueFlagSrc = (buffer&0x40 != 0)  // 0b01000000
	nav.BranchOption = (buffer & 0x0F)              // 0b00001111

	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return nil, fmt.Errorf("failed to read buffer: %w", err)
	}
	nav.CompareOption = (buffer & 0x0F) // 0b00001111

	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return nil, fmt.Errorf("failed to read buffer: %w", err)
	}
	nav.SetOption = (buffer & 0x1F)

	if err := binary.Read(file, binary.BigEndian, &nav.Destination); err != nil {
		return nil, fmt.Errorf("failed to read Destination: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &nav.Source); err != nil {
		return nil, fmt.Errorf("failed to read Source: %w", err)
	}

	return nav, err
}
