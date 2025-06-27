package mpls

import (
	"encoding/binary"
	"fmt"
	"io"
)

// BasicAttributes is a base structure for all stream attributes.
type BasicAttributes struct {
	Length           uint8
	StreamCodingType StreamCodingType
}

// SetLength sets the Length for the BasicAttributes.
func (attr *BasicAttributes) SetLength(length uint8) {
	attr.Length = length
}

// SetStreamCodingType sets the StreamCodingType for the BasicAttributes.
func (attr *BasicAttributes) SetStreamCodingType(streamCodingType StreamCodingType) {
	attr.StreamCodingType = streamCodingType
}

// PrimaryVideoAttributesH264 is used for Primary Video streams.
// It contains the basic attributes and specific attributes for H264 (mpeg4 AVC),
// and other video codecs like MPEG1, MPEG2, and VC1.
type PrimaryVideoAttributesH264 struct {
	BasicAttributes
	Format VideoFormatType // 0b11110000
	Rate   VideoRateType   // 0b00001111
}

// PrimaryVideoAttributesHEVC is used for Primary Video streams with HEVC (H265) codec.
// HEVC is High Efficiency Video Coding, also known as H.265.
// It contains the basic attributes and specific attributes for HEVC.
// It includes dynamic range type, color space, and flags for CR and HDR+.
type PrimaryVideoAttributesHEVC struct {
	BasicAttributes
	Format           VideoFormatType // 0b11110000
	Rate             VideoRateType   // 0b00001111
	DynamicRangeType uint8           // 0b11110000
	ColorSpace       uint8           // 0b00001111
	CRFlag           bool            // 0b10000000
	HDRPlusFlag      bool            // 0b01000000
}

// PrimaryAudioAttributes is used for Primary Audio streams.
type PrimaryAudioAttributes struct {
	BasicAttributes
	Format       AudioFormatType // 0b11110000
	Rate         AudioRateType   // 0b00001111
	LanguageCode [3]byte
}

// SecondaryAudioExtraAttributes is used for Secondary Audio streams.
type SecondaryAudioExtraAttributes struct {
	NumberOfPrimaryAudioRef uint8
	PrimaryAudioRefs        []uint8
}

// SecondaryAudioAttributes is used for Secondary Audio streams.
type SecondaryAudioAttributes struct {
	PrimaryAudioAttributes
	SecondaryAudioExtraAttributes
}

// SecondaryVideoExtraAttributes is used for Secondary Video streams.
type SecondaryVideoExtraAttributes struct {
	NumberOfSecondaryAudioRef uint8
	NumberOfPIPPGRef          uint8
	SecondaryAudioRefs        []uint8
	PIPPGRefs                 []uint8
}

// SecondaryVideoAttributes is used for Secondary Video streams.
type SecondaryVideoAttributes struct {
	PrimaryVideoAttributesH264
	SecondaryVideoExtraAttributes
}

// GraphicsAttributes is a base structure for all graphics-related attributes.
type GraphicsAttributes struct {
	BasicAttributes
	LanguageCode [3]byte
}

// PGAttributes is used for Presentation Graphics (PG) subtitles.
type PGAttributes struct {
	GraphicsAttributes
}

// IGAttributes is used for Interactive Graphics (IG) subtitles.
type IGAttributes struct {
	GraphicsAttributes
}

// TextAttributes is used for Text (PG) subtitles.
type TextAttributes struct {
	GraphicsAttributes
	CharacterCode uint8
}

// StreamAttributes is an interface that defines the methods for reading and setting
type StreamAttributes interface {
	Read(io.ReadSeeker) error
	SetLength(uint8)
	SetStreamCodingType(StreamCodingType)
}

// ReadStreamAttributes reads the stream attributes from the provided io.ReadSeeker.
// It expects the file (io.ReadSeeker) to be positioned at the start of the
// StreamAttributes structure.
// It determines the type of stream attributes based on the StreamCodingType byte and
// reads the corresponding structure accordingly.
// It returns a StreamAttributes interface and an error if any occurs during reading.
func ReadStreamAttributes(file io.ReadSeeker, kindOf StreamTypeKindOf) (attr StreamAttributes, err error) {
	var length uint8
	var streamCodingType StreamCodingType

	if err = binary.Read(file, binary.BigEndian, &length); err != nil {
		return nil, fmt.Errorf("failed to read Attributes length: %w", err)
	}

	if err = binary.Read(file, binary.BigEndian, &streamCodingType); err != nil {
		return nil, fmt.Errorf("failed to read Attributes streamCodingType: %w", err)
	}

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

// Read implements the StreamAttributes interface for PrimaryVideoAttributesH264.
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

// Read implements the StreamAttributes interface for PrimaryVideoAttributesH264.
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

// Read implements the StreamAttributes interface for PrimaryAudioAttributes.
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

// Read implements the StreamAttributes interface for SecondaryAudioAttributes.
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

// Read implements the StreamAttributes interface for SecondaryVideoAttributes.
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

// Read implements the StreamAttributes interface for PGAttributes.
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

// Read implements the StreamAttributes interface for IGAttributes.
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

// Read implements the StreamAttributes interface for TextAttributess.
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
