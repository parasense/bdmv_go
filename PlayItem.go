package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
)

// PlayItem represents a single item in the playlist
type PlayItem struct {
	Length                   uint16
	ClipInformationFileName  [5]byte // 5-digit string like "00063"
	ClipCodecIdentifier      [4]byte // 4-character string like "M2TS"
	IsMultiAngle             bool
	ConnectionCondition      uint8
	RefToSTCID               uint8
	INTime                   uint32 // Timestamp in 45kHz
	OUTTime                  uint32 // Timestamp in 45kHz
	UserOptions              *UserOptions
	PlayItemRandomAccessFlag bool
	StillMode                uint8 // 0x00 == none ; 0x01 == finite still time (StillTime takes a uint16 value) ; 0x02 == infinite still time
	StillTime                uint16
	NumberOfAngles           uint8
	IsDifferentAudios        bool
	IsSeamlessAngleChange    bool
	Angles                   []*PlayItemEntry
	StreamTable              *StreamTable
}

func ReadPlayItem(file io.ReadSeeker) (*PlayItem, error) {

	playItem := &PlayItem{}
	if err := binary.Read(file, binary.BigEndian, &playItem.Length); err != nil {
		return nil, fmt.Errorf("failed to read play item length: %w", err)
	}

	// The 5 bytes clip name
	if err := binary.Read(file, binary.BigEndian, &playItem.ClipInformationFileName); err != nil {
		return nil, fmt.Errorf("failed to read clip info filename: %w", err)
	}

	// The 4 byte codec should be something like "M2TS"
	if err := binary.Read(file, binary.BigEndian, &playItem.ClipCodecIdentifier); err != nil {
		return nil, fmt.Errorf("failed to read clip codec identifier: %w", err)
	}

	// skip 1 byte reserve space
	if _, err := file.Seek(1, io.SeekCurrent); err != nil {
		return nil, fmt.Errorf("failed to seek past reserve space: %w", err)
	}

	// Read 1 byte into a buffer to extract IsMultiAngle bit flag and ConnectionCondition 4 bit number
	var flagBuffer uint8
	if err := binary.Read(file, binary.BigEndian, &flagBuffer); err != nil {
		return nil, fmt.Errorf("failed to read flagBuffer: %w", err)
	}

	playItem.IsMultiAngle = flagBuffer&0x10 != 0     // 0b00010000
	playItem.ConnectionCondition = flagBuffer & 0x0F // 0b00001111

	if err := binary.Read(file, binary.BigEndian, &playItem.RefToSTCID); err != nil {
		return nil, fmt.Errorf("failed to read STC ID: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &playItem.INTime); err != nil {
		return nil, fmt.Errorf("failed to read IN time: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &playItem.OUTTime); err != nil {
		return nil, fmt.Errorf("failed to read OUT time: %w", err)
	}

	var err error

	// This reads 8 bytes
	playItem.UserOptions, err = ReadUserOptions(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read UserOptions: %w", err)
	}

	// Read the random access flag (1 bit)
	if err := binary.Read(file, binary.BigEndian, &flagBuffer); err != nil {
		return nil, fmt.Errorf("failed to read flagBuffer: %w", err)
	}
	playItem.PlayItemRandomAccessFlag = flagBuffer&0x80 != 0

	// Still mode 1 byte
	if err := binary.Read(file, binary.BigEndian, &playItem.StillMode); err != nil {
		return nil, fmt.Errorf("failed to read StillMode: %w", err)
	}

	// Read StillTime if StillMode enabled
	if playItem.StillMode == 1 {

		// Read two bytes of StillTime
		if err := binary.Read(file, binary.BigEndian, &playItem.StillTime); err != nil {
			return nil, fmt.Errorf("failed to read StillTime: %w", err)
		}
	} else {

		// Else, Skip two bytes that would have been StillTime
		if _, err := file.Seek(2, io.SeekCurrent); err != nil {
			return nil, fmt.Errorf("failed to seek past reserve space: %w", err)
		}
	}

	if playItem.IsMultiAngle {

		// Read one byte of NumberOfAngles
		if err := binary.Read(file, binary.BigEndian, &playItem.NumberOfAngles); err != nil {
			return nil, fmt.Errorf("failed to read NumberOfAngles: %w", err)
		}

		// This is what libbluray does - no idea why.
		if playItem.NumberOfAngles < 1 {
			playItem.NumberOfAngles = 1
		}

		// Read the IsDifferentAudios & IsSeamlessAngleChange flags (1 bit each)
		if err := binary.Read(file, binary.BigEndian, &flagBuffer); err != nil {
			return nil, fmt.Errorf("failed to read flagBuffer: %w", err)
		}
		playItem.IsDifferentAudios = flagBuffer&0x02 != 0
		playItem.IsSeamlessAngleChange = flagBuffer&0x01 != 0
	} else {
		playItem.NumberOfAngles = 1
	}

	playItem.Angles = make([]*PlayItemEntry, playItem.NumberOfAngles)

	// Copy The (already parsed) clip information into clip[0]
	// They are the zero'th clip of multiAngle things.
	playItem.Angles[0] = &PlayItemEntry{
		FileName:   playItem.ClipInformationFileName,
		Codec:      playItem.ClipCodecIdentifier,
		RefToSTCID: playItem.RefToSTCID,
	}

	for i := uint8(1); i < playItem.NumberOfAngles; i++ {
		if playItem.Angles[i], err = ReadPlayItemEntry(file); err != nil {
			return nil, fmt.Errorf("failed to read PlayItemEntry: %w", err)
		}
	}

	// Read Stream Table (formerly SteamInfo)
	if playItem.StreamTable, err = ReadStreamTable(file); err != nil {
		return nil, fmt.Errorf("failed to read StreamTable: %w", err)
	}

	// XXX - debug opportunity
	// Print info about file pointer possition in the file.
	// Print info about the size of the play items.
	// Allow to see from output any bytes not processed.

	return playItem, nil
}

func (playItem *PlayItem) Print() {
	inTime := Convert45KhzTimeToSeconds(playItem.INTime)
	outTime := Convert45KhzTimeToSeconds(playItem.OUTTime)
	PadPrintf(4, "Length: %d\n", playItem.Length)
	PadPrintf(4, "Clip File: %s\n", playItem.ClipInformationFileName)
	PadPrintf(4, "Codec ID: %s\n", playItem.ClipCodecIdentifier)
	PadPrintf(4, "Multi-angle: %v\n", playItem.IsMultiAngle)
	PadPrintf(4, "ConnectionCondition: %d\n", playItem.ConnectionCondition)
	PadPrintf(4, "RefToSTCID: %d\n", playItem.RefToSTCID)
	PadPrintf(4, "InTime: %d (%d)\n", playItem.INTime, inTime)
	PadPrintf(4, "OUTime: %d (%d)\n", playItem.OUTTime, outTime)
	PadPrintf(6, "*Duration: %v\n", outTime-inTime)
	playItem.UserOptions.Print()
	PadPrintf(4, "PlayItemRandomAccessFlag: %v\n", playItem.PlayItemRandomAccessFlag)
	PadPrintf(4, "StillMode: %v\n", playItem.StillMode)
	PadPrintf(4, "StillTime: %v\n", playItem.StillTime)
	PadPrintf(4, "NumberOfAngles: %v\n", playItem.NumberOfAngles)
	PadPrintf(4, "IsDifferentAudios: %v\n", playItem.IsDifferentAudios)
	PadPrintf(4, "IsSeamlessAngleChange: %v\n", playItem.IsSeamlessAngleChange)
	PadPrintf(4, "Angles:\n")
	for i := uint8(0); i < uint8(len(playItem.Angles)); i++ {
		PadPrintf(6, "Angle: [%d]:\n", i)
		playItem.Angles[i].Print()
	}
	playItem.StreamTable.Print()
}

func (playItem *PlayItem) Assert() {

	// Length should not be zero
	if playItem.Length == 0 {
		log.Fatal("Assertion: playItem.Length must not equal zero.")
	}

	// FileName should be numeric
	if !isNumeric(playItem.ClipInformationFileName) {
		log.Fatal("Assertion: playItem.ClipInformationFileName must be numeric.")
	}

	// Codec ID should be upper case alpha numeric
	if !isAlphanumericUppercase(playItem.ClipCodecIdentifier) {
		log.Fatal("Assertion: playItem.ClipCodecIdentifier must be upper case ascii")
	}

	// INTime should always be less-than OUTTime
	if playItem.INTime > playItem.OUTTime {
		log.Fatal("Assertion: playItem.INTime must be less-than playItem.OUTTime.")
	}

	// multi-angle true entails NumberOfAngles is non-zero
	if playItem.IsMultiAngle && playItem.NumberOfAngles == 0 {
		log.Fatal("Assertion: NumberOfAngles must not equal zero when IsMultiAngle is true.")
	}

	// When stillMode == 1 (finite time), then StillTime must be greater-than zero
	// When stillMode == 2 (infinite time), then StillTime must (probably) be zero
	if playItem.StillMode == 1 && playItem.StillTime == 0 {
		log.Fatal("Assertion: playItem.StillTime must not equal zero when playItem.StillMode equals one.")
	} else if playItem.StillMode == 2 && playItem.StillTime != 0 {
		log.Fatal("Assertion: playItem.StillTime must equal zero when playItem.StillMode equals two.")
	}

}
