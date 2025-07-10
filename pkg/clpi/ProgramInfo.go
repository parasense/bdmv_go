package clpi

import (
	"encoding/binary"
	"fmt"
	"io"
)

type ProgramInfo struct {
	Length           uint32
	NumberOfPrograms uint8
	Programs         []*Program
}

func ReadProgramInfo(file io.ReadSeeker, offsets *OffsetsUint32) (programInfo *ProgramInfo, err error) {
	programInfo = &ProgramInfo{}

	// Jump to start address
	if _, err := file.Seek(offsets.Start, io.SeekStart); err != nil {
		return nil, fmt.Errorf("Failed to Seek offset.Start: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &programInfo.Length); err != nil {
		return nil, fmt.Errorf("Failed to read Length: %w", err)
	}

	// 1-byte reserve space
	if _, err := file.Seek(1, io.SeekCurrent); err != nil {
		return nil, fmt.Errorf("Error seeking reserve space: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &programInfo.NumberOfPrograms); err != nil {
		return nil, fmt.Errorf("Error reading NumberOfPrograms: %w", err)
	}

	programInfo.Programs = make([]*Program, programInfo.NumberOfPrograms)
	for i := range programInfo.Programs {
		if programInfo.Programs[i], err = ReadProgram(file); err != nil {
			return nil, fmt.Errorf("Failed in call to ReadProgram(): %w", err)
		}
	}

	return programInfo, nil
}

func (pi *ProgramInfo) String() string {
	return fmt.Sprintf(
		"ProgramInfo{"+
			"Length: %d, "+
			"NumberOfPrograms: %d, "+
			"Programs: %v}",
		pi.Length, pi.NumberOfPrograms,
		pi.Programs,
	)
}

type Program struct {
	SPNProgramSequenceStart uint32
	ProgramMapPID           uint16
	NumberOfStreamsInPS     uint8
	ProgramStreams          []*ProgramStream
}

func ReadProgram(file io.ReadSeeker) (p *Program, err error) {
	p = &Program{}

	if err := binary.Read(file, binary.BigEndian, &p.SPNProgramSequenceStart); err != nil {
		return nil, fmt.Errorf("Error reading SPNProgramSequenceStart: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &p.ProgramMapPID); err != nil {
		return nil, fmt.Errorf("Error reading ProgramMapPID: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &p.NumberOfStreamsInPS); err != nil {
		return nil, fmt.Errorf("Error reading NumberOfStreamsInPS: %w", err)
	}

	// 1-byte reserve space
	if _, err := file.Seek(1, io.SeekCurrent); err != nil {
		return nil, fmt.Errorf("Error reserve space: %w", err)
	}

	p.ProgramStreams = make([]*ProgramStream, p.NumberOfStreamsInPS)

	for i := range p.ProgramStreams {
		p.ProgramStreams[i], err = ReadProgramStream(file)
		if err != nil {
			return nil, fmt.Errorf("Error returned by ReadProgramStream() %w", err)
		}
	}

	return p, nil
}

func (p *Program) String() string {
	return fmt.Sprintf(
		"Program{"+
			"SPNProgramSequenceStart: %d, "+
			"ProgramMapPID: %d, "+
			"NumberOfStreamsInPS: %d, "+
			"ProgramStreams: %v}",
		p.SPNProgramSequenceStart,
		p.ProgramMapPID,
		p.NumberOfStreamsInPS,
		p.ProgramStreams,
	)
}

type ProgramStream struct {
	StreamPID        uint16
	StreamCodingInfo []StreamCodingInfo
}

func (ps *ProgramStream) String() string {
	return fmt.Sprintf(
		"ProgramStream{"+
			"StreamPID: %d, "+
			"StreamCodingInfo: %v}",
		ps.StreamPID, ps.StreamCodingInfo,
	)
}

type StreamCodingInfo interface {
	Read(io.ReadSeeker) error
	SetLength(uint8)
	SetStreamCodingType(StreamCodingType)
	SetISRCode([12]byte)
}

func ReadProgramStream(file io.ReadSeeker) (p *ProgramStream, err error) {

	p = &ProgramStream{}

	if err := binary.Read(file, binary.BigEndian, &p.StreamPID); err != nil {
		return nil, err
	}

	var length uint8
	if err := binary.Read(file, binary.BigEndian, &length); err != nil {
		return nil, err
	}
	end, _ := CalculateEndOffset(file, length)

	var streamCodingType StreamCodingType
	if err := binary.Read(file, binary.BigEndian, &streamCodingType); err != nil {
		return nil, err
	}

	streamCodingInfo := NewStreamCodingInfo(streamCodingType)
	if streamCodingInfo == nil {
		return nil, fmt.Errorf("Unknown (nil) streamCodingInfo ")
	}
	streamCodingInfo.SetLength(length)
	streamCodingInfo.SetStreamCodingType(streamCodingType)

	if err := streamCodingInfo.Read(file); err != nil {
		return nil, err
	}

	var isrc [12]byte
	if err := binary.Read(file, binary.BigEndian, &isrc); err != nil {
		return nil, err
	}
	streamCodingInfo.SetISRCode(isrc)

	// skip any tail padding
	if _, err := file.Seek(end, io.SeekStart); err != nil {
		return nil, err
	}

	p.StreamCodingInfo = append(p.StreamCodingInfo, streamCodingInfo)

	return p, nil
}

type BaseStreamCodingInfo struct {
	Length           uint8
	StreamCodingType StreamCodingType
	ISRCode          [12]byte // International Standard Recording Code
}

func (base *BaseStreamCodingInfo) SetLength(len uint8) { base.Length = len }
func (base *BaseStreamCodingInfo) SetStreamCodingType(code StreamCodingType) {
	base.StreamCodingType = code
}
func (base *BaseStreamCodingInfo) SetISRCode(code [12]byte) { base.ISRCode = code }

type StreamCodingInfoH264 struct {
	BaseStreamCodingInfo
	VideoFormat      VideoFormatType      // 4-bits 0b11110000
	FrameRate        VideoRateType        // 4-bits 0b00001111
	VideoAspectRatio VideoAspectRatioType // 4-bits 0b11110000
	OCFlag           bool                 // 1-bit  0b00000010
}

type StreamCodingInfoH265 struct {
	BaseStreamCodingInfo
	VideoFormat      VideoFormatType      // 4-bits 0b11110000
	FrameRate        VideoRateType        // 4-bits 0b00001111
	VideoAspectRatio VideoAspectRatioType // 4-bits 0b11110000
	OCFlag           bool                 // 1-bit  0b00000010
	CRFlag           bool                 // 1-bit  0b00000001
	DynamicRangeType uint8                // 4-bits 0b11110000
	ColorSpace       uint8                // 4-bits 0b00001111
	HDRPlusFlag      bool                 // 1-bit  0b10000000
}

type StreamCodingInfoAudio struct {
	BaseStreamCodingInfo
	AudioFormat  AudioFormatType // 4-bits 0b11110000
	SampleRate   AudioRateType   // 4-bits 0b00001111
	LanguageCode [3]byte         // 3-bytes
}

type StreamCodingTypePG struct {
	BaseStreamCodingInfo
	LanguageCode [3]byte // 3-bytes
}

type StreamCodingTypeIG struct {
	BaseStreamCodingInfo
	LanguageCode [3]byte // 3-bytes
}

type StreamCodingTypeText struct {
	BaseStreamCodingInfo
	CharacterCode uint8   // 1-byte
	LanguageCode  [3]byte // 3-bytes
}

func NewStreamCodingInfo(streamCodingType StreamCodingType) StreamCodingInfo {
	switch streamCodingType {

	case // PrimaryVideo or SecondaryVideo types
		STREAM_TYPE_VIDEO_MPEG1,
		STREAM_TYPE_VIDEO_MPEG2,
		STREAM_TYPE_VIDEO_H264,
		STREAM_TYPE_VIDEO_H264_MVC,
		STREAM_TYPE_VIDEO_VC1:
		return &StreamCodingInfoH264{}

	case // H265 (HEVC) PrimaryVideo
		STREAM_TYPE_VIDEO_HEVC:
		return &StreamCodingInfoH265{}

	case // Primary & Secondary Audio
		STREAM_TYPE_AUDIO_LPCM,
		STREAM_TYPE_AUDIO_AC3,
		STREAM_TYPE_AUDIO_DTS,
		STREAM_TYPE_AUDIO_TRUHD,
		STREAM_TYPE_AUDIO_AC3PLUS,
		STREAM_TYPE_AUDIO_DTSHD,
		STREAM_TYPE_AUDIO_DTSHD_MASTER,
		STREAM_TYPE_AUDIO_AC3PLUS_SECONDARY,
		STREAM_TYPE_AUDIO_DTSHD_SECONDARY:
		return &StreamCodingInfoAudio{}

	case // Presentation Graphics
		STREAM_TYPE_SUB_PG:
		return &StreamCodingTypePG{}

	case // Interactive Graphics
		STREAM_TYPE_SUB_IG:
		return &StreamCodingTypeIG{}

	case // Text (PG) subtitles
		STREAM_TYPE_SUB_TEXT:
		return &StreamCodingTypeText{}

	default:
		return nil // XXX - add error handling here

	}

}

func (s *StreamCodingInfoH264) Read(file io.ReadSeeker) error {

	var buffer byte
	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return err
	}
	s.VideoFormat = VideoFormatType((buffer & 0xF0) >> 4)
	s.FrameRate = VideoRateType(buffer & 0x0F)

	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return err
	}

	s.VideoAspectRatio = VideoAspectRatioType((buffer & 0xF0) >> 4)
	s.OCFlag = buffer&0x02 != 0 // 0b00000010

	return nil
}

func (s *StreamCodingInfoH265) Read(file io.ReadSeeker) error {

	var buffer byte
	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return err
	}
	s.VideoFormat = VideoFormatType((buffer & 0xF0) >> 4)
	s.FrameRate = VideoRateType(buffer & 0x0F)

	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return err
	}

	s.VideoAspectRatio = VideoAspectRatioType((buffer & 0xF0) >> 4)
	s.OCFlag = buffer&0x02 != 0 // 0b00000010
	s.CRFlag = buffer&0x01 != 0 // 0b00000001

	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return err
	}
	s.DynamicRangeType = (buffer & 0xF0) >> 4
	s.ColorSpace = (buffer & 0x0F)

	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return err
	}
	s.HDRPlusFlag = buffer&0x80 != 0 // 0b10000000

	return nil
}

func (s *StreamCodingInfoAudio) Read(file io.ReadSeeker) error {

	var buffer byte
	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return err
	}
	s.AudioFormat = AudioFormatType((buffer & 0xF0) >> 4)
	s.SampleRate = AudioRateType(buffer & 0x0F)

	if err := binary.Read(file, binary.BigEndian, &s.LanguageCode); err != nil {
		return err
	}

	return nil
}

func (s *StreamCodingTypePG) Read(file io.ReadSeeker) error {

	if err := binary.Read(file, binary.BigEndian, &s.LanguageCode); err != nil {
		return err
	}

	return nil
}

func (s *StreamCodingTypeIG) Read(file io.ReadSeeker) error {

	if err := binary.Read(file, binary.BigEndian, &s.LanguageCode); err != nil {
		return err
	}

	return nil
}

func (s *StreamCodingTypeText) Read(file io.ReadSeeker) error {

	if err := binary.Read(file, binary.BigEndian, &s.CharacterCode); err != nil {
		return err
	}

	if err := binary.Read(file, binary.BigEndian, &s.LanguageCode); err != nil {
		return err
	}

	return nil
}

func (sci *StreamCodingInfoH264) String() string {
	return fmt.Sprintf(
		"StreamCodingInfoH264{"+
			"Length: %d, "+
			"StreamCodingType: %d, "+
			"VideoFormat: %d, "+
			"FrameRate: %d, "+
			"VideoAspectRatio: %d, "+
			"OCFlag: %t, "+
			"ISRCode: %s}",
		sci.Length, sci.StreamCodingType,
		sci.VideoFormat, sci.FrameRate,
		sci.VideoAspectRatio, sci.OCFlag, sci.ISRCode,
	)
}

func (sci *StreamCodingInfoH265) String() string {
	return fmt.Sprintf(
		"StreamCodingInfoH265{"+
			"Length: %d, "+
			"StreamCodingType: %d, "+
			"VideoFormat: %d, "+
			"FrameRate: %d, "+
			"VideoAspectRatio: %d, "+
			"OCFlag: %t, "+
			"DynamicRangeType: %d, "+
			"ColorSpace: %d, "+
			"HDRPlusFlag: %t, "+
			"ISRCode: %s}",
		sci.Length, sci.StreamCodingType,
		sci.VideoFormat, sci.FrameRate,
		sci.VideoAspectRatio, sci.OCFlag,
		sci.DynamicRangeType, sci.ColorSpace,
		sci.HDRPlusFlag, sci.ISRCode,
	)
}

func (sci *StreamCodingInfoAudio) String() string {
	return fmt.Sprintf(
		"StreamCodingInfoAudio{"+
			"Length: %d, "+
			"StreamCodingType: %d, "+
			"AudioFormat: %d, "+
			"FrameRate: %d, "+
			"LanguageCode: %s, "+
			"ISRCode: %s}",
		sci.Length, sci.StreamCodingType,
		sci.AudioFormat, sci.SampleRate,
		sci.LanguageCode[:], sci.ISRCode,
	)
}

func (sci *StreamCodingTypePG) String() string {
	return fmt.Sprintf(
		"StreamCodingTypePG{"+
			"Length: %d, "+
			"StreamCodingType: %d, "+
			"LanguageCode: %s, "+
			"ISRCode: %s}",
		sci.Length, sci.StreamCodingType,
		sci.LanguageCode[:], sci.ISRCode,
	)
}

func (sci *StreamCodingTypeIG) String() string {
	return fmt.Sprintf(
		"StreamCodingInfoIG{"+
			"Length: %d, "+
			"StreamCodingType: %d, "+
			"LanguageCode: %s, "+
			"ISRCode: %s}",
		sci.Length, sci.StreamCodingType,
		sci.LanguageCode[:], sci.ISRCode,
	)
}

func (sci *StreamCodingTypeText) String() string {
	return fmt.Sprintf(
		"StreamCodingInfoH264{"+
			"Length: %d, "+
			"StreamCodingType: %d, "+
			"LanguageCode: %s, "+
			"CharacterCode: %d, "+
			"ISRCode: %s}",
		sci.Length, sci.StreamCodingType,
		sci.LanguageCode[:], sci.CharacterCode, sci.ISRCode,
	)
}
