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

type BasicAttributes struct {
	Length           uint8
	StreamCodingType StreamCodingType
}

func (attr *BasicAttributes) SetLength(length uint8) {
	attr.Length = length
}

func (attr *BasicAttributes) SetStreamCodingType(streamCodingType StreamCodingType) {
	attr.StreamCodingType = streamCodingType
}

type PrimaryVideoAttributesH264 struct {
	BasicAttributes
	Format VideoFormatType // 0b11110000
	Rate   VideoRateType   // 0b00001111
}

type PrimaryVideoAttributesHEVC struct {
	BasicAttributes
	Format           VideoFormatType // 0b11110000
	Rate             VideoRateType   // 0b00001111
	DynamicRangeType uint8           // 0b11110000
	ColorSpace       uint8           // 0b00001111
	CRFlag           bool            // 0b10000000
	HDRPlusFlag      bool            // 0b01000000
}

type PrimaryAudioAttributes struct {
	BasicAttributes
	Format       AudioFormatType // 0b11110000
	Rate         AudioRateType   // 0b00001111
	LanguageCode [3]byte
}

type SecondaryAudioExtraAttributes struct {
	NumberOfPrimaryAudioRef uint8
	PrimaryAudioRefs        []uint8
}

type SecondaryAudioAttributes struct {
	PrimaryAudioAttributes
	SecondaryAudioExtraAttributes
}

type SecondaryVideoExtraAttributes struct {
	NumberOfSecondaryAudioRef uint8
	NumberOfPIPPGRef          uint8
	SecondaryAudioRefs        []uint8
	PIPPGRefs                 []uint8
}

type SecondaryVideoAttributes struct {
	PrimaryVideoAttributesH264
	SecondaryVideoExtraAttributes
}

type GraphicsAtttributes struct {
	BasicAttributes
	LanguageCode [3]byte
}

type PGAttributes struct {
	GraphicsAtttributes
}

type IGAttributes struct {
	GraphicsAtttributes
}

type TextAttributes struct {
	GraphicsAtttributes
	CharacterCode uint8
}

// This enables combinations
type StreamAttributes interface {
	Print()
	Read(io.ReadSeeker) error
	SetLength(uint8)
	SetStreamCodingType(StreamCodingType)
}

func ReadStreamAttributes(file io.ReadSeeker, kindOf StreamTypeKindOf) (attr StreamAttributes, err error) {
	//var buffer byte
	var length uint8
	var streamCodingType StreamCodingType

	if err = binary.Read(file, binary.BigEndian, &length); err != nil {
		return nil, fmt.Errorf("failed to read Attributes length: %w", err)
	}
	//var length uint8 = buffer
	//length = buffer

	if err = binary.Read(file, binary.BigEndian, &streamCodingType); err != nil {
		return nil, fmt.Errorf("failed to read Attributes streamCodingType: %w", err)
	}
	//var streamCodingType StreamCodingType = buffer
	//StreamCodingType = buffer

	//// XXX
	//fmt.Printf("DEBUG: StreamAttr.kindOf: %s\n", kindOf)
	//fmt.Printf("DEBUG: StreamAttr.Length: %d\n", length)
	//fmt.Printf("DEBUG: StreamAttr.streamCodingType: %d\n", streamCodingType)

	switch streamCodingType {

	case // PrimaryVideo or SecondaryVideo types
		STREAM_TYPE_VIDEO_MPEG1,
		STREAM_TYPE_VIDEO_MPEG2,
		STREAM_TYPE_VIDEO_H264,
		STREAM_TYPE_VIDEO_H264_MVC,
		STREAM_TYPE_VIDEO_VC1:
		switch kindOf {
		case "PrimaryVideo":
			attr = &PrimaryVideoAttributesH264{}
		case "SecondaryVideo":
			attr = &SecondaryVideoAttributes{}
		default:
			attr = &PrimaryVideoAttributesH264{}
		}

	case // H265 (HEVC) PrimaryVideo
		STREAM_TYPE_VIDEO_HEVC:
		attr = &PrimaryVideoAttributesHEVC{}

	case // Primary Audio
		STREAM_TYPE_AUDIO_LPCM,
		STREAM_TYPE_AUDIO_AC3,
		STREAM_TYPE_AUDIO_DTS,
		STREAM_TYPE_AUDIO_TRUHD,
		STREAM_TYPE_AUDIO_AC3PLUS,
		STREAM_TYPE_AUDIO_DTSHD,
		STREAM_TYPE_AUDIO_DTSHD_MASTER:
		attr = &PrimaryAudioAttributes{}

	case // Secondary Audio
		STREAM_TYPE_AUDIO_AC3PLUS_SECONDARY,
		STREAM_TYPE_AUDIO_DTSHD_SECONDARY:
		attr = &SecondaryAudioAttributes{}

	case // Presentation Graphics
		STREAM_TYPE_SUB_PG:
		attr = &PGAttributes{}

	case // Interactive Graphics
		STREAM_TYPE_SUB_IG:
		attr = &IGAttributes{}

	case // Text (PG) subtitles
		STREAM_TYPE_SUB_TEXT:
		attr = &TextAttributes{}

	default:
		return nil, fmt.Errorf("Unknown Stream Atrribute code type: (%d)", streamCodingType)

	}

	if attr != nil {
		attr.SetLength(length)
		attr.SetStreamCodingType(streamCodingType)
		if err := attr.Read(file); err != nil {
			return nil, fmt.Errorf("failed call attr.Read() on attribute code type (%d) %w", streamCodingType, err)
		}
	}

	return attr, nil
}

func (attr *PrimaryVideoAttributesH264) Print() {
	PadPrintln(10, "Attributes:")
	PadPrintf(12, "Length: %d\n", attr.Length)
	PadPrintf(12, "StreamCodingType: %d [%s]\n", attr.StreamCodingType, StreamCodec(attr.StreamCodingType))
	PadPrintf(12, "Format: %d [%s]\n", attr.Format, VideoFormat(attr.Format))
	PadPrintf(12, "Rate: %d [%s] \n", attr.Rate, VideoRate(attr.Rate))
}

func (attr *PrimaryVideoAttributesHEVC) Print() {
	PadPrintln(10, "Attributes:")
	PadPrintf(12, "Length: %d\n", attr.Length)
	PadPrintf(12, "StreamCodingType: %d [%s]\n", attr.StreamCodingType, StreamCodec(attr.StreamCodingType))
	PadPrintf(12, "Format: %d [%s]\n", attr.Format, VideoFormat(attr.Format))
	PadPrintf(12, "Rate: %d [%s]\n", attr.Rate, VideoRate(attr.Rate))
	PadPrintf(12, "DynamicRangeType: %d\n", attr.DynamicRangeType)
	PadPrintf(12, "ColorSpace: %d\n", attr.ColorSpace)
	PadPrintf(12, "CRFlag: %v\n", attr.CRFlag)
	PadPrintf(12, "HDRPlusFlag: %v\n", attr.HDRPlusFlag)
}

func (attr *PrimaryAudioAttributes) Print() {
	PadPrintln(10, "Attributes:")
	PadPrintf(12, "Length: %d\n", attr.Length)
	PadPrintf(12, "StreamCodingType: %d [%s]\n", attr.StreamCodingType, StreamCodec(attr.StreamCodingType))
	PadPrintf(12, "Format: %d [%s]\n", attr.Format, AudioFormat(attr.Format))
	PadPrintf(12, "Rate: %d [%s]\n", attr.Rate, AudioRate(attr.Rate))
	eng, nat := LanguageCode(attr.LanguageCode)
	PadPrintf(12, "LanguageCode: %s [%s, %s]\n", attr.LanguageCode, eng, nat)
}

func (attr *SecondaryAudioAttributes) Print() {
	PadPrintln(10, "Attributes:")
	PadPrintf(12, "Length: %d\n", attr.Length)
	PadPrintf(12, "StreamCodingType: %d [%s]\n", attr.StreamCodingType, StreamCodec(attr.StreamCodingType))
	PadPrintf(12, "Format: %d [%s]\n", attr.Format, AudioFormat(attr.Format))
	PadPrintf(12, "Rate: %d [%s]\n", attr.Rate, AudioRate(attr.Rate))
	eng, nat := LanguageCode(attr.LanguageCode)
	PadPrintf(12, "LanguageCode: %s [%s, %s]\n", attr.LanguageCode, eng, nat)
	PadPrintf(12, "NumberOfPrimaryAudioRef: %+v\n", attr.NumberOfPrimaryAudioRef)

}

func (attr *SecondaryVideoAttributes) Print() {
	PadPrintln(10, "Attributes:")
	PadPrintf(12, "Length: %d\n", attr.Length)
	PadPrintf(12, "StreamCodingType: %d [%s]\n", attr.StreamCodingType, StreamCodec(attr.StreamCodingType))
	PadPrintf(12, "Format: %d [%s]\n", attr.Format, VideoFormat(attr.Format))
	PadPrintf(12, "Rate: %d [%s]\n", attr.Rate, VideoRate(attr.Rate))
}

func (attr *PGAttributes) Print() {
	PadPrintln(10, "Attributes:")
	PadPrintf(12, "Length: %d\n", attr.Length)
	PadPrintf(12, "StreamCodingType: %d [%s]\n", attr.StreamCodingType, StreamCodec(attr.StreamCodingType))
	eng, nat := LanguageCode(attr.LanguageCode)
	PadPrintf(12, "LanguageCode: %s [%s, %s]\n", attr.LanguageCode, eng, nat)
}

// type 145
func (attr *IGAttributes) Print() {
	PadPrintln(10, "Attributes:")
	PadPrintf(12, "Length: %d\n", attr.Length)
	PadPrintf(12, "StreamCodingType: %d [%s]\n", attr.StreamCodingType, StreamCodec(attr.StreamCodingType))
	eng, nat := LanguageCode(attr.LanguageCode)
	PadPrintf(12, "LanguageCode: %s [%s, %s]\n", attr.LanguageCode, eng, nat)
}

func (attr *TextAttributes) Print() {
	PadPrintln(10, "Attributes:")
	PadPrintf(12, "Length: %d\n", attr.Length)
	PadPrintf(12, "StreamCodingType: %d [%s]\n", attr.StreamCodingType, StreamCodec(attr.StreamCodingType))
	PadPrintf(12, "CharacterCode: %d [%s]\n", attr.CharacterCode, CharacterCode(attr.CharacterCode))
	eng, nat := LanguageCode(attr.LanguageCode)
	PadPrintf(12, "LanguageCode: %s [%s, %s]\n", attr.LanguageCode, eng, nat)
}

func (attr *PrimaryVideoAttributesH264) Read(file io.ReadSeeker) (err error) {
	var buffer byte

	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return fmt.Errorf("failed to read Attributes MPEG buffer: %w", err)
	}
	attr.Format = stnFormat[VideoFormatType](&buffer) // 0b11110000
	attr.Rate = stnRate[VideoRateType](&buffer)       // 0b00001111

	// 3 byte tail padding
	if _, err := file.Seek(3, io.SeekCurrent); err != nil {
		return fmt.Errorf("failed to seek past Attributes MPEG reserve space: %w", err)
	}
	return nil
}

func (attr *PrimaryVideoAttributesHEVC) Read(file io.ReadSeeker) (err error) {
	var buffer byte

	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return fmt.Errorf("failed to read Attributes HEVC buffer: %w", err)
	}
	attr.Format = stnFormat[VideoFormatType](&buffer) // 0b11110000
	attr.Rate = stnRate[VideoRateType](&buffer)       // 0b00001111

	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return fmt.Errorf("failed to read Attributes HEVC buffer: %w", err)
	}
	attr.DynamicRangeType = stnDynamicRangeType(&buffer) // 0b11110000
	attr.ColorSpace = stnColorSpace(&buffer)             // 0b00001111

	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return fmt.Errorf("failed to read Attributes HEVC buffer: %w", err)
	}
	attr.CRFlag = buffer&0x80 != 0      // 0b10000000
	attr.HDRPlusFlag = buffer&0x40 != 0 // 0b01000000

	// 1 byte tail padding
	if _, err := file.Seek(1, io.SeekCurrent); err != nil {
		return fmt.Errorf("failed to seek past Attributes HVEC reserve space: %v\n", err)
	}
	return nil
}

func (attr *PrimaryAudioAttributes) Read(file io.ReadSeeker) (err error) {
	var buffer byte

	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return fmt.Errorf("failed to read Attributes AUDIO buffer: %w", err)
	}
	attr.Format = stnFormat[AudioFormatType](&buffer) // 0b11110000
	attr.Rate = stnRate[AudioRateType](&buffer)       // 0b00001111

	if err := binary.Read(file, binary.BigEndian, &attr.LanguageCode); err != nil {
		return fmt.Errorf("failed to read Attributes AUDIO LanguageCode: %w", err)
	}
	return nil
}

func (attr *SecondaryAudioAttributes) Read(file io.ReadSeeker) (err error) {
	var buffer byte

	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return fmt.Errorf("failed to read Attributes AUDIO buffer: %w", err)
	}
	attr.Format = stnFormat[AudioFormatType](&buffer) // 0b11110000
	attr.Rate = stnRate[AudioRateType](&buffer)       // 0b00001111

	if err := binary.Read(file, binary.BigEndian, &attr.LanguageCode); err != nil {
		return fmt.Errorf("failed to read Attributes AUDIO LanguageCode: %w", err)
	}

	//
	// Extra Attributes
	//
	if err := binary.Read(file, binary.BigEndian, &attr.NumberOfPrimaryAudioRef); err != nil {
		return fmt.Errorf("failed to read SECONDARY AUDIO NumberOfPrimaryAudioRef: %w", err)
	}

	if _, err := file.Seek(1, io.SeekCurrent); err != nil {
		return fmt.Errorf("failed to seek past Attributes reserve space: %w", err)
	}

	if attr.NumberOfPrimaryAudioRef > 0 {
		attr.PrimaryAudioRefs = make([]uint8, attr.NumberOfPrimaryAudioRef)
		for i := range attr.NumberOfPrimaryAudioRef {
			if err := binary.Read(file, binary.BigEndian, &attr.PrimaryAudioRefs[i]); err != nil {
				return fmt.Errorf("failed to read SECONDARY Video  SecondaryAudioRefs: %w", err)
			}
		}

		// If the NumberOf is odd, then 1-byte of tail padding/reserve
		if attr.NumberOfPrimaryAudioRef%2 != 0 {
			if _, err := file.Seek(1, io.SeekCurrent); err != nil {
				return fmt.Errorf("failed to seek past reserve space: %w", err)
			}
		}
	}

	return nil
}

func (attr *SecondaryVideoAttributes) Read(file io.ReadSeeker) (err error) {
	var buffer byte

	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return fmt.Errorf("failed to read Attributes MPEG buffer: %w", err)
	}
	attr.Format = stnFormat[VideoFormatType](&buffer) // 0b11110000
	attr.Rate = stnRate[VideoRateType](&buffer)       // 0b00001111

	// 3 byte tail padding
	if _, err := file.Seek(3, io.SeekCurrent); err != nil {
		return fmt.Errorf("failed to seek past Attributes MPEG reserve space: %w", err)
	}

	//
	// Extra Attributes
	//
	if err := binary.Read(file, binary.BigEndian, &attr.NumberOfSecondaryAudioRef); err != nil {
		return fmt.Errorf("failed to read SECONDARY Video NumberOfSecondaryAudioRef: %w", err)
	}

	if _, err := file.Seek(1, io.SeekCurrent); err != nil {
		return fmt.Errorf("failed to seek past reserve space: %w", err)
	}

	if attr.NumberOfSecondaryAudioRef > 0 {
		attr.SecondaryAudioRefs = make([]uint8, attr.NumberOfSecondaryAudioRef)
		for i := range attr.NumberOfSecondaryAudioRef {
			if err := binary.Read(file, binary.BigEndian, &attr.SecondaryAudioRefs[i]); err != nil {
				return fmt.Errorf("failed to read SECONDARY Video  SecondaryAudioRefs: %w", err)
			}
		}

		// If the NumberOf is odd, then 1-byte of tail padding/reserve
		if attr.NumberOfSecondaryAudioRef%2 != 0 {
			if _, err := file.Seek(1, io.SeekCurrent); err != nil {
				return fmt.Errorf("failed to seek past reserve space: %w", err)
			}
		}
	}

	if err := binary.Read(file, binary.BigEndian, &attr.NumberOfPIPPGRef); err != nil {
		return fmt.Errorf("failed to read SECONDARY Video NumberOfPIPPGRef: %w", err)
	}

	if _, err := file.Seek(1, io.SeekCurrent); err != nil {
		return fmt.Errorf("failed to seek past reserve space: %w", err)
	}

	if attr.NumberOfPIPPGRef > 0 {
		attr.PIPPGRefs = make([]uint8, attr.NumberOfPIPPGRef)
		for i := range attr.NumberOfPIPPGRef {
			if err := binary.Read(file, binary.BigEndian, &attr.PIPPGRefs[i]); err != nil {
				return fmt.Errorf("failed to read SECONDARY Video PIPPGRefs: %w", err)
			}
		}

		// If the NumberOf is odd, then 1-byte of tail padding/reserve
		if attr.NumberOfPIPPGRef%2 != 0 {
			if _, err := file.Seek(1, io.SeekCurrent); err != nil {
				return fmt.Errorf("failed to seek past reserve space: %w", err)
			}
		}
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

//
// helper functions
//

// Common bit reads
func fourBitsHigh(buffer *byte) uint8 { return (*buffer & 0xF0) >> 4 }
func fourBitsLow(buffer *byte) uint8  { return (*buffer & 0x0F) }

// Wrappers for the sake of readability
func stnFormat[formatType ~uint8](buffer *byte) formatType { return formatType(fourBitsHigh(buffer)) }
func stnRate[rateType ~uint8](buffer *byte) rateType       { return rateType(fourBitsLow(buffer)) }
func stnDynamicRangeType(buffer *byte) uint8               { return fourBitsHigh(buffer) }
func stnColorSpace(buffer *byte) uint8                     { return fourBitsLow(buffer) }
