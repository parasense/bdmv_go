package main

import (
	"encoding/binary"
	"fmt"
	"io"
)

type UserOptions struct {
	MenuCall                         bool
	TitleSearch                      bool
	ChapterSearch                    bool
	TimeSearch                       bool
	SkipToNextPoint                  bool
	SkipToPrevPoint                  bool
	Stop                             bool
	PauseOn                          bool
	StillOff                         bool
	ForwardPlay                      bool
	BackwardPlay                     bool
	Resume                           bool
	MoveUpSelectedButton             bool
	MoveDownSelectedButton           bool
	MoveLeftSelectedButton           bool
	MoveRightSelectedButton          bool
	SelectButton                     bool
	ActivateButton                   bool
	SelectAndActivateButton          bool
	PrimaryAudioStreamNumberChange   bool
	AngleNumberChange                bool
	PopupOn                          bool
	PopupOff                         bool
	PrimaryPGEnableDisable           bool
	PrimaryPGStreamNumberChange      bool
	SecondaryVideoEnableDisable      bool
	SecondaryVideoStreamNumberChange bool
	SecondaryAudioEnableDisable      bool
	SecondaryAudioStreamNumberChange bool
	SecondaryPGStreamNumberChange    bool
}

func (userOptions *UserOptions) getFlag(data *uint8, mask uint8) bool {
	return *data&mask != 0
}

func ReadUserOptions(file io.ReadSeeker) (userOptions *UserOptions, err error) {
	userOptions = &UserOptions{}

	var flagBuffer byte
	if err := binary.Read(file, binary.BigEndian, &flagBuffer); err != nil {
		return nil, fmt.Errorf("failed to read UO mask table: %w", err)
	}
	userOptions.MenuCall = userOptions.getFlag(&flagBuffer, 0x80)
	userOptions.TitleSearch = userOptions.getFlag(&flagBuffer, 0x40)
	userOptions.ChapterSearch = userOptions.getFlag(&flagBuffer, 0x20)
	userOptions.TimeSearch = userOptions.getFlag(&flagBuffer, 0x10)
	userOptions.SkipToNextPoint = userOptions.getFlag(&flagBuffer, 0x08)
	userOptions.SkipToPrevPoint = userOptions.getFlag(&flagBuffer, 0x04)
	userOptions.Stop = userOptions.getFlag(&flagBuffer, 0x01)

	if err := binary.Read(file, binary.BigEndian, &flagBuffer); err != nil {
		return nil, fmt.Errorf("failed to read UO mask table: %w", err)
	}
	userOptions.PauseOn = userOptions.getFlag(&flagBuffer, 0x80)
	userOptions.StillOff = userOptions.getFlag(&flagBuffer, 0x20)
	userOptions.ForwardPlay = userOptions.getFlag(&flagBuffer, 0x10)
	userOptions.BackwardPlay = userOptions.getFlag(&flagBuffer, 0x08)
	userOptions.Resume = userOptions.getFlag(&flagBuffer, 0x04)
	userOptions.MoveUpSelectedButton = userOptions.getFlag(&flagBuffer, 0x02)
	userOptions.MoveDownSelectedButton = userOptions.getFlag(&flagBuffer, 0x01)

	if err := binary.Read(file, binary.BigEndian, &flagBuffer); err != nil {
		return nil, fmt.Errorf("failed to read UO mask table: %w", err)
	}
	userOptions.MoveLeftSelectedButton = userOptions.getFlag(&flagBuffer, 0x80)
	userOptions.MoveRightSelectedButton = userOptions.getFlag(&flagBuffer, 0x40)
	userOptions.SelectButton = userOptions.getFlag(&flagBuffer, 0x20)
	userOptions.ActivateButton = userOptions.getFlag(&flagBuffer, 0x10)
	userOptions.SelectAndActivateButton = userOptions.getFlag(&flagBuffer, 0x08)
	userOptions.PrimaryAudioStreamNumberChange = userOptions.getFlag(&flagBuffer, 0x04)
	userOptions.AngleNumberChange = userOptions.getFlag(&flagBuffer, 0x01)

	if err := binary.Read(file, binary.BigEndian, &flagBuffer); err != nil {
		return nil, fmt.Errorf("failed to read UO mask table: %w", err)
	}
	userOptions.PopupOn = userOptions.getFlag(&flagBuffer, 0x80)
	userOptions.PopupOff = userOptions.getFlag(&flagBuffer, 0x40)
	userOptions.PrimaryPGEnableDisable = userOptions.getFlag(&flagBuffer, 0x20)
	userOptions.PrimaryPGStreamNumberChange = userOptions.getFlag(&flagBuffer, 0x10)
	userOptions.SecondaryVideoEnableDisable = userOptions.getFlag(&flagBuffer, 0x08)
	userOptions.SecondaryVideoStreamNumberChange = userOptions.getFlag(&flagBuffer, 0x04)
	userOptions.SecondaryAudioEnableDisable = userOptions.getFlag(&flagBuffer, 0x02)
	userOptions.SecondaryAudioStreamNumberChange = userOptions.getFlag(&flagBuffer, 0x01)

	if err := binary.Read(file, binary.BigEndian, &flagBuffer); err != nil {
		return nil, fmt.Errorf("failed to read UO mask table: %w", err)
	}
	userOptions.SecondaryPGStreamNumberChange = userOptions.getFlag(&flagBuffer, 0x40)

	// skip 3 byte reserve padding space
	if _, err := file.Seek(3, io.SeekCurrent); err != nil {
		return nil, fmt.Errorf("failed to seek past reserve space: %w", err)
	}

	return userOptions, nil
}

func (userOptions *UserOptions) Print() {
	PadPrintln(4, "UserOptions:")
	PadPrintf(6, "MenuCall: %v\n", userOptions.MenuCall)
	PadPrintf(6, "TitleSearch: %v\n", userOptions.TitleSearch)
	PadPrintf(6, "ChapterSearch: %v\n", userOptions.ChapterSearch)
	PadPrintf(6, "TimeSearch: %v\n", userOptions.TimeSearch)
	PadPrintf(6, "SkipToNextPoint: %v\n", userOptions.SkipToNextPoint)
	PadPrintf(6, "SkipToPrevPoint: %v\n", userOptions.SkipToPrevPoint)
	PadPrintf(6, "Stop: %v\n", userOptions.Stop)
	PadPrintf(6, "PauseOn: %v\n", userOptions.PauseOn)
	PadPrintf(6, "StillOff %v\n", userOptions.StillOff)
	PadPrintf(6, "ForwardPlay: %v\n", userOptions.ForwardPlay)
	PadPrintf(6, "BackwardPlay: %v\n", userOptions.BackwardPlay)
	PadPrintf(6, "Resume: %v\n", userOptions.Resume)
	PadPrintf(6, "MoveUpSelectedButton: %v\n", userOptions.MoveUpSelectedButton)
	PadPrintf(6, "MoveDownSelectedButton: %v\n", userOptions.MoveDownSelectedButton)
	PadPrintf(6, "MoveLeftSelectedButton: %v\n", userOptions.MoveLeftSelectedButton)
	PadPrintf(6, "MoveRightSelectedButton: %v\n", userOptions.MoveRightSelectedButton)
	PadPrintf(6, "SelectButton: %v\n", userOptions.SelectButton)
	PadPrintf(6, "ActivateButton: %v\n", userOptions.ActivateButton)
	PadPrintf(6, "SelectAndActivateButton: %v\n", userOptions.SelectAndActivateButton)
	PadPrintf(6, "PrimaryAudioStreamNumberChange: %v\n", userOptions.PrimaryAudioStreamNumberChange)
	PadPrintf(6, "AngleNumberChange: %v\n", userOptions.AngleNumberChange)
	PadPrintf(6, "PopupOn: %v\n", userOptions.PopupOn)
	PadPrintf(6, "PopupOff: %v\n", userOptions.PopupOff)
	PadPrintf(6, "PrimaryPGEnableDisable: %v\n", userOptions.PrimaryPGEnableDisable)
	PadPrintf(6, "PrimaryPGStreamNumberChange: %v\n", userOptions.PrimaryPGStreamNumberChange)
	PadPrintf(6, "SecondaryVideoEnableDisable: %v\n", userOptions.SecondaryVideoEnableDisable)
	PadPrintf(6, "SecondaryVideoStreamNumberChange: %v\n", userOptions.SecondaryVideoStreamNumberChange)
	PadPrintf(6, "SecondaryAudioEnableDisable: %v\n", userOptions.SecondaryAudioEnableDisable)
	PadPrintf(6, "SecondaryAudioStreamNumberChange: %v\n", userOptions.SecondaryAudioStreamNumberChange)
	PadPrintf(6, "SecondaryPGStreamNumberChange: %v\n", userOptions.SecondaryPGStreamNumberChange)
}
