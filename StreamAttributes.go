package main

import (
	"encoding/binary"
	"fmt"
	"io"
)

/*
	NOTES:
	Each Stream Type has specific attributes.
	Each Stream Type has different binary coding rules.
*/

/*
// PROPOSED UNIFIED DATA STRUCTURE
// At the very least the variable names can be made generic.
// StreamInfo represents a Blu-ray stream's attributes.
type StreamInfo struct {
	CodingType uint8    // Stream coding type (e.g., 0x1B for H.264, 0x92 for TextST)
	Format     uint8    // Format (e.g., audio format, video format)
	Rate       uint8    // Rate (e.g., sample rate, frame rate)
	CharCode   uint8    // Character code for text subtitles
	Lang       [4]byte  // Language code (e.g., "eng\0")
	Aspect     uint8    // Aspect ratio for video
	PID        uint16   // Packet ID
}
*/

type BasicAttributes struct {
	Length           uint8
	StreamCodingType uint8
}

func (attr *BasicAttributes) SetLength(length uint8) {
	attr.Length = length
}

func (attr *BasicAttributes) SetStreamCodingType(streamCodingType uint8) {
	attr.StreamCodingType = streamCodingType
}

type BasicAudioVideoAttributes struct {
	BasicAttributes
	Format uint8 // 4 bits high
	Rate   uint8 // 4 bits low
}

type ModernVideoAttributes struct {
	BasicAudioVideoAttributes
	DynamicRangeType uint8 // 4 bits high
	ColorSpace       uint8 // 4 bits low
	CRFlag           bool
	HDRPlusFlag      bool
}

type AudioAttributes struct {
	BasicAudioVideoAttributes
	LanguageCode [3]byte
}

type SubTitleAtttributes struct {
	BasicAttributes
	LanguageCode [3]byte
}

type PGAttributes struct {
	SubTitleAtttributes
}

type IGAttributes struct {
	SubTitleAtttributes
}

type TextAttributes struct {
	SubTitleAtttributes
	CharacterCode uint8
}

// This enables combinations
type StreamAttributes interface {
	Print()
	Read(io.ReadSeeker) error
	SetLength(uint8)
	SetStreamCodingType(uint8)
}

func ReadStreamAttributes(file io.ReadSeeker) (attr StreamAttributes, err error) {
	var buffer byte

	if err = binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return nil, fmt.Errorf("failed to read Attributes length: %w", err)
	}
	var length uint8 = buffer

	if err = binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return nil, fmt.Errorf("failed to read Attributes streamCodingType: %w", err)
	}
	var streamCodingType uint8 = buffer

	switch streamCodingType {
	case 1, 2, 27, 234:
		attr = &BasicAudioVideoAttributes{}

	case 36:
		attr = &ModernVideoAttributes{}

	case 0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0xA1, 0xA2:
		attr = &AudioAttributes{}

	case 0x90: // 144 Pesentation Graphics
		attr = &PGAttributes{}

	case 0x91: // 145 Interactive Graphics
		attr = &IGAttributes{}

	case 0x92: // 146 Text subtitles (goes with Pesentation Graphics)
		attr = &PGAttributes{}

	default:
		return nil, fmt.Errorf("Unknown Stsream Atrribute code type: (%d)", streamCodingType)

	}
	attr.SetLength(length)
	attr.SetStreamCodingType(streamCodingType)
	if err := attr.Read(file); err != nil {
		return nil, fmt.Errorf("failed call attr.Read() on attribute code type (%d) %w", streamCodingType, err)
	}
	return attr, nil
}

func (attr *BasicAudioVideoAttributes) Print() {
	PadPrintln(6, "Stream Attributes:")
	PadPrintf(8, "Length: %d\n", attr.Length)
	PadPrintf(8, "StreamCodingType: %d\n", attr.StreamCodingType)
	PadPrintf(8, "Format: %d\n", attr.Format)
	PadPrintf(8, "Rate: %d\n", attr.Rate)
}

func (attr *ModernVideoAttributes) Print() {
	PadPrintln(6, "Stream Attributes:")
	PadPrintf(8, "Length: %d\n", attr.Length)
	PadPrintf(8, "StreamCodingType: %d\n", attr.StreamCodingType)
	PadPrintf(8, "Format: %d\n", attr.Format)
	PadPrintf(8, "Rate: %d\n", attr.Rate)
	PadPrintf(8, "DynamicRangeType: %d\n", attr.DynamicRangeType)
	PadPrintf(8, "ColorSpace: %d\n", attr.ColorSpace)
	PadPrintf(8, "CRFlag: %v\n", attr.CRFlag)
	PadPrintf(8, "HDRPlusFlag: %v\n", attr.HDRPlusFlag)
}

func (attr *AudioAttributes) Print() {
	PadPrintln(6, "Stream Attributes:")
	PadPrintf(8, "Length: %d\n", attr.Length)
	PadPrintf(8, "StreamCodingType: %d\n", attr.StreamCodingType)
	PadPrintf(8, "Format: %d\n", attr.Format)
	PadPrintf(8, "Rate: %d\n", attr.Rate)
	PadPrintf(8, "LanguageCode: %s\n", attr.LanguageCode)
}

func (attr *PGAttributes) Print() {
	PadPrintln(6, "Stream Attributes:")
	PadPrintf(8, "Length: %d\n", attr.Length)
	PadPrintf(8, "StreamCodingType: %d\n", attr.StreamCodingType)
	PadPrintf(8, "LanguageCode: %s\n", attr.LanguageCode)
}

// type 145
func (attr *IGAttributes) Print() {
	PadPrintln(6, "Stream Attributes:")
	PadPrintf(8, "Length: %d\n", attr.Length)
	PadPrintf(8, "StreamCodingType: %d\n", attr.StreamCodingType)
	PadPrintf(8, "LanguageCode: %s\n", attr.LanguageCode)
}

func (attr *TextAttributes) Print() {
	PadPrintln(6, "Stream Attributes:")
	PadPrintf(8, "Length: %d\n", attr.Length)
	PadPrintf(8, "StreamCodingType: %d\n", attr.StreamCodingType)
	PadPrintf(8, "CharacterCode: %d\n", attr.CharacterCode)
	PadPrintf(8, "LanguageCode: %s\n", attr.LanguageCode)
}

func (attr *BasicAudioVideoAttributes) Read(file io.ReadSeeker) (err error) {
	var buffer byte

	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return fmt.Errorf("failed to read Attributes MPEG buffer: %w", err)
	}
	attr.Format = (buffer & 0xF0) >> 4
	attr.Rate = buffer & 0x0F

	// 3 byte tail padding
	if _, err := file.Seek(3, io.SeekCurrent); err != nil {
		return fmt.Errorf("failed to seek past Attributes MPEG reserve space: %w", err)
	}
	return nil
}

func (attr *ModernVideoAttributes) Read(file io.ReadSeeker) (err error) {
	var buffer byte

	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return fmt.Errorf("failed to read Attributes HEVC buffer: %w", err)
	}
	attr.Format = (buffer & 0xF0) >> 4
	attr.Rate = buffer & 0x0F

	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return fmt.Errorf("failed to read Attributes HEVC buffer: %w", err)
	}
	attr.DynamicRangeType = (buffer & 0xF0) >> 4
	attr.ColorSpace = buffer & 0x0F

	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return fmt.Errorf("failed to read Attributes HEVC buffer: %w", err)
	}
	attr.CRFlag = buffer&0x80 != 0
	attr.HDRPlusFlag = buffer&0x40 != 0

	// 1 byte tail padding
	if _, err := file.Seek(1, io.SeekCurrent); err != nil {
		return fmt.Errorf("failed to seek past Attributes HVEC reserve space: %v\n", err)
	}
	return nil
}

func (attr *AudioAttributes) Read(file io.ReadSeeker) (err error) {
	var buffer byte

	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return fmt.Errorf("failed to read Attributes AUDIO buffer: %w", err)
	}
	attr.Format = (buffer & 0xF0) >> 4
	attr.Rate = buffer & 0x0F

	if err := binary.Read(file, binary.BigEndian, &attr.LanguageCode); err != nil {
		return fmt.Errorf("failed to read Attributes AUDIO LanguageCode: %w", err)
	}
	return nil
}

func (attr *PGAttributes) Read(file io.ReadSeeker) (err error) {
	if err := binary.Read(file, binary.BigEndian, &attr.LanguageCode); err != nil {
		return fmt.Errorf("failed to read Attributes PG LanguageCode: %w", err)
	}

	// 1 byte tail padding
	if _, err := file.Seek(1, io.SeekCurrent); err != nil {
		return fmt.Errorf("failed to seek past Attributes PG reserve space: %w", err)
	}
	return nil

}

// The exact same as PG type
func (attr *IGAttributes) Read(file io.ReadSeeker) (err error) {
	if err := binary.Read(file, binary.BigEndian, &attr.LanguageCode); err != nil {
		return fmt.Errorf("failed to read Attributes IG LanguageCode: %w", err)
	}

	// 1 byte tail padding
	if _, err := file.Seek(1, io.SeekCurrent); err != nil {
		return fmt.Errorf("failed to seek past Attributes IG reserve space: %w", err)
	}
	return nil

}

func (attr *TextAttributes) Read(file io.ReadSeeker) (err error) {
	if err := binary.Read(file, binary.BigEndian, &attr.CharacterCode); err != nil {
		return fmt.Errorf("failed to read Attributes TEXT CharacterCode: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &attr.LanguageCode); err != nil {
		return fmt.Errorf("failed to read Attributes TEXT LanguageCode: %w", err)
	}
	return nil

}
