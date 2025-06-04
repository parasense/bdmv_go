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
}

func OLDReadStreamAttributes(file io.ReadSeeker) (attr StreamAttributes, err error) {
	var buffer byte

	if err = binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return nil, fmt.Errorf("failed to read Attributes length: %v\n", err)
	}
	var length uint8 = buffer

	if err = binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return nil, fmt.Errorf("failed to read Attributes streamCodingType: %v\n", err)
	}
	var streamCodingType uint8 = buffer

	//if attr, err = NewStreamAttributes(streamCodingType); err != nil {
	//	return nil, fmt.Errorf("failed to Create Attributes data structure: %v\n", err)
	//}

	switch streamCodingType {
	case 1, 2, 27, 234:
		attr := &BasicAudioVideoAttributes{}
		attr.Length = length
		attr.StreamCodingType = streamCodingType
		if err := attr.Read(file); err != nil {
			return nil, fmt.Errorf("failed call attr.Read() on MPEG attributes: %v\n", err)
		}
		return attr, nil

	case 36:
		attr := &ModernVideoAttributes{}
		attr.Length = length
		attr.StreamCodingType = streamCodingType
		if err := attr.Read(file); err != nil {
			return nil, fmt.Errorf("failed call attr.Read() on HEVC attributes: %v\n", err)
		}
		return attr, nil

	case 0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0xA1, 0xA2:
		attr := &AudioAttributes{}
		attr.Length = length
		attr.StreamCodingType = streamCodingType
		if err := attr.Read(file); err != nil {
			return nil, fmt.Errorf("failed call attr.Read() on AUDIO attributes: %v\n", err)
		}
		return attr, nil

	case 0x90: // 144
		attr := &PGAttributes{}
		attr.Length = length
		attr.StreamCodingType = streamCodingType
		if err := attr.Read(file); err != nil {
			return nil, fmt.Errorf("failed call attr.Read() on PG attributes: %v\n", err)
		}
		return attr, nil

	case 0x91: // 145
		attr := &IGAttributes{}
		attr.Length = length
		attr.StreamCodingType = streamCodingType
		if err := attr.Read(file); err != nil {
			return nil, fmt.Errorf("failed call attr.Read() on IG attributes: %v\n", err)
		}
		return attr, nil

	case 0x92: // 146
		attr := &PGAttributes{}
		attr.Length = length
		attr.StreamCodingType = streamCodingType
		if err := attr.Read(file); err != nil {
			return nil, fmt.Errorf("failed call attr.Read() on IG attributes: %v\n", err)
		}
		return attr, nil

	default:
		return nil, fmt.Errorf("failed to read TEXT stsream Atrributes\n")

	}
}

// Factory function
func NewStreamAttributes(length uint8, streamCodingType uint8) (StreamAttributes, error) {
	switch streamCodingType {

	// Note: 0x20 was added to support 3d video.
	// 1, 2, 27, 234, 0x20(3d video)
	case 0x01, 0x02, 0x1b, 0xea, 0x20:
		StreamAttributes := &BasicAudioVideoAttributes{}
		StreamAttributes.Length = length
		StreamAttributes.StreamCodingType = streamCodingType
		return StreamAttributes, nil

	case 36:
		StreamAttributes := &ModernVideoAttributes{}
		StreamAttributes.Length = length
		StreamAttributes.StreamCodingType = streamCodingType
		return StreamAttributes, nil

	case 0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0xA1, 0xA2:
		StreamAttributes := &AudioAttributes{}
		StreamAttributes.Length = length
		StreamAttributes.StreamCodingType = streamCodingType
		return StreamAttributes, nil

	case 0x90:
		StreamAttributes := &PGAttributes{}
		StreamAttributes.Length = length
		StreamAttributes.StreamCodingType = streamCodingType
		return StreamAttributes, nil

	case 0x91:
		StreamAttributes := &IGAttributes{}
		StreamAttributes.Length = length
		StreamAttributes.StreamCodingType = streamCodingType
		return StreamAttributes, nil

	case 0x92:
		StreamAttributes := &TextAttributes{}
		StreamAttributes.Length = length
		StreamAttributes.StreamCodingType = streamCodingType
		return StreamAttributes, nil

	default:
		return nil, fmt.Errorf("Unknown Stream Attribute type: [%v]\n", streamCodingType)
	}
}

func ReadStreamAttributes(file io.ReadSeeker) (StreamAttributes, error) {
	var buffer byte

	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return nil, fmt.Errorf("failed to read Attributes length: %v\n", err)
	}
	var length uint8 = buffer

	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return nil, fmt.Errorf("failed to read Attributes streamCodingType: %v\n", err)
	}
	var streamCodingType uint8 = buffer

	attr, err := NewStreamAttributes(length, streamCodingType)
	if err != nil {
		return nil, fmt.Errorf("failed to Create Attributes data structure: %v\n", err)
	}

	if err := attr.Read(file); err != nil {
		return nil, fmt.Errorf("failed call attr.Read() on attributes: %v\n", err)
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
		return fmt.Errorf("failed to read Attributes MPEG buffer: %v\n", err)
	}
	attr.Format = (buffer & 0xF0) >> 4
	attr.Rate = buffer & 0x0F

	// 3 byte tail padding
	if _, err := file.Seek(3, io.SeekCurrent); err != nil {
		return fmt.Errorf("failed to seek past Attributes MPEG reserve space: %v\n", err)
	}
	return nil
}

func (attr *ModernVideoAttributes) Read(file io.ReadSeeker) (err error) {
	var buffer byte

	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return fmt.Errorf("failed to read Attributes HEVC buffer: %v\n", err)
	}
	attr.Format = (buffer & 0xF0) >> 4
	attr.Rate = buffer & 0x0F

	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return fmt.Errorf("failed to read Attributes HEVC buffer: %v\n", err)
	}
	attr.DynamicRangeType = (buffer & 0xF0) >> 4
	attr.ColorSpace = buffer & 0x0F

	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return fmt.Errorf("failed to read Attributes HEVC buffer: %v\n", err)
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
		return fmt.Errorf("failed to read Attributes AUDIO buffer: %v\n", err)
	}
	attr.Format = (buffer & 0xF0) >> 4
	attr.Rate = buffer & 0x0F

	if err := binary.Read(file, binary.BigEndian, &attr.LanguageCode); err != nil {
		return fmt.Errorf("failed to read Attributes AUDIO LanguageCode: %v\n", err)
	}
	return nil
}

func (attr *PGAttributes) Read(file io.ReadSeeker) (err error) {
	if err := binary.Read(file, binary.BigEndian, &attr.LanguageCode); err != nil {
		return fmt.Errorf("failed to read Attributes PG LanguageCode: %v\n", err)
	}

	// 1 byte tail padding
	if _, err := file.Seek(1, io.SeekCurrent); err != nil {
		return fmt.Errorf("failed to seek past Attributes PG reserve space: %v\n", err)
	}
	return nil

}

// The exact same as PG type
func (attr *IGAttributes) Read(file io.ReadSeeker) (err error) {
	if err := binary.Read(file, binary.BigEndian, &attr.LanguageCode); err != nil {
		return fmt.Errorf("failed to read Attributes IG LanguageCode: %v\n", err)
	}

	// 1 byte tail padding
	if _, err := file.Seek(1, io.SeekCurrent); err != nil {
		return fmt.Errorf("failed to seek past Attributes IG reserve space: %v\n", err)
	}
	return nil

}

func (attr *TextAttributes) Read(file io.ReadSeeker) (err error) {
	if err := binary.Read(file, binary.BigEndian, &attr.CharacterCode); err != nil {
		return fmt.Errorf("failed to read Attributes TEXT CharacterCode: %v\n", err)
	}

	if err := binary.Read(file, binary.BigEndian, &attr.LanguageCode); err != nil {
		return fmt.Errorf("failed to read Attributes TEXT LanguageCode: %v\n", err)
	}
	return nil

}
