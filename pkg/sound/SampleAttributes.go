package sound

import (
	"encoding/binary"
	"fmt"
	"io"
)

type SampleAttributes struct {
	NumberOfChannels uint8  // 4-bits 0b11110000
	SampleRate       uint32 // 4-bits 0b00001111
	BitsPerSample    uint8  // 4-bits 0b11000000
	SoundDataIndex   uint32
	NumberOfFrames   uint32
}

type AudioChannelType uint8

const (
	AUDIO_CHANNELS_MONO   AudioChannelType = 1
	AUDIO_CHANNELS_STEREO AudioChannelType = 3
)

func audioChannelType(code AudioChannelType) (channelType uint8, err error) {
	switch code {
	case AUDIO_CHANNELS_MONO:
		return 1, nil
	case AUDIO_CHANNELS_STEREO:
		return 2, nil
	default:
		return 0, fmt.Errorf("unknown channel type: %v", code)
	}
}

func ReadSampleAttributes(file io.ReadSeeker) (sampleAttr *SampleAttributes, err error) {
	sampleAttr = &SampleAttributes{}

	var buffer uint8
	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return nil, fmt.Errorf("failed to read buffer: %w", err)
	}

	// 4-bits 0b11110000
	if sampleAttr.NumberOfChannels, err = audioChannelType(
		AudioChannelType(
			(buffer & 0xF0) >> 4)); err != nil {
		return nil, err
	}

	// 4-bits 0b00001111
	switch buffer & 0x0F {
	default:
		fallthrough
	case 1:
		sampleAttr.SampleRate = 48000
	}

	if err := binary.Read(file, binary.BigEndian, &buffer); err != nil {
		return nil, fmt.Errorf("failed to read buffer: %w", err)
	}

	// 2-bits 0b11000000
	switch (buffer & 0xC0) >> 6 {
	default:
		fallthrough
	case 1:
		sampleAttr.BitsPerSample = 16
	}

	if err := binary.Read(file, binary.BigEndian, &sampleAttr.SoundDataIndex); err != nil {
		return nil, fmt.Errorf("failed to read SoundDataIndexes: %w", err)
	}

	// Read value into the destination as temporary
	if err := binary.Read(file, binary.BigEndian, &sampleAttr.NumberOfFrames); err != nil {
		return nil, fmt.Errorf("failed to read NumberOfFrames: %w", err)
	}

	// This puts the frames in terms of bytes
	sampleAttr.NumberOfFrames /= uint32((sampleAttr.BitsPerSample / 8) * sampleAttr.NumberOfChannels)

	return sampleAttr, err
}
