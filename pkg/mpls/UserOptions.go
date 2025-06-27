package mpls

import (
	"encoding/binary"
	"fmt"
	"io"
)

// UserOptions represents the User Options (UO) section in an MPLS file.
// It contains various flags that indicate the capabilities and options available for user interaction.
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

// ReadUserOptions reads the User Options from the given file.
// It expects the file to be at the correct position where the User Options start.
// The User Options consist of several flags that are read from the file.
// It returns a pointer to the UserOptions struct and an error if any occurs during reading.
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

func (userOptions *UserOptions) String() string {
	return fmt.Sprintf("UserOptions:\n"+
		"  MenuCall: %t\n"+
		"  TitleSearch: %t\n"+
		"  ChapterSearch: %t\n"+
		"  TimeSearch: %t\n"+
		"  SkipToNextPoint: %t\n"+
		"  SkipToPrevPoint: %t\n"+
		"  Stop: %t\n"+
		"  PauseOn: %t\n"+
		"  StillOff: %t\n"+
		"  ForwardPlay: %t\n"+
		"  BackwardPlay: %t\n"+
		"  Resume: %t\n"+
		"  MoveUpSelectedButton: %t\n"+
		"  MoveDownSelectedButton: %t\n"+
		"  MoveLeftSelectedButton: %t\n"+
		"  MoveRightSelectedButton: %t\n"+
		"  SelectButton: %t\n"+
		"  ActivateButton: %t\n"+
		"  SelectAndActivateButton: %t\n"+
		"  PrimaryAudioStreamNumberChange: %t\n"+
		"  AngleNumberChange: %t\n"+
		"  PopupOn: %t\n"+
		"  PopupOff: %t\n"+
		"  PrimaryPGEnableDisable: %t\n"+
		"  PrimaryPGStreamNumberChange: %t\n"+
		"  SecondaryVideoEnableDisable: %t\n"+
		"  SecondaryVideoStreamNumberChange: %t\n"+
		"  SecondaryAudioEnableDisable: %t\n"+
		"  SecondaryAudioStreamNumberChange: %t\n"+
		"  SecondaryPGStreamNumberChange: %t\n",
		userOptions.MenuCall,
		userOptions.TitleSearch,
		userOptions.ChapterSearch,
		userOptions.TimeSearch,
		userOptions.SkipToNextPoint,
		userOptions.SkipToPrevPoint,
		userOptions.Stop,
		userOptions.PauseOn,
		userOptions.StillOff,
		userOptions.ForwardPlay,
		userOptions.BackwardPlay,
		userOptions.Resume,
		userOptions.MoveUpSelectedButton,
		userOptions.MoveDownSelectedButton,
		userOptions.MoveLeftSelectedButton,
		userOptions.MoveRightSelectedButton,
		userOptions.SelectButton,
		userOptions.ActivateButton,
		userOptions.SelectAndActivateButton,
		userOptions.PrimaryAudioStreamNumberChange,
		userOptions.AngleNumberChange,
		userOptions.PopupOn,
		userOptions.PopupOff,
		userOptions.PrimaryPGEnableDisable,
		userOptions.PrimaryPGStreamNumberChange,
		userOptions.SecondaryVideoEnableDisable,
		userOptions.SecondaryVideoStreamNumberChange,
		userOptions.SecondaryAudioEnableDisable,
		userOptions.SecondaryAudioStreamNumberChange,
		userOptions.SecondaryPGStreamNumberChange,
	)
}
