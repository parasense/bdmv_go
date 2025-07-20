package sound

import (
	"encoding/binary"
	"fmt"
	"io"
)

type SoundMetaData struct {
	Length         uint32
	NumberOfSounds uint8
	SampleAttrs    []*SampleAttributes
}

func ReadSoundMetaData(file io.ReadSeeker, offsets *OffsetsUint32) (soundData *SoundMetaData, err error) {

	// Jump to start address
	if _, err := file.Seek(offsets.Start, io.SeekStart); err != nil {
		return nil, fmt.Errorf("failed to seek to start address: %w", err)
	}

	soundData = &SoundMetaData{}

	if err := binary.Read(file, binary.BigEndian, &soundData.Length); err != nil {
		return nil, fmt.Errorf("failed to read play item length: %w", err)
	}

	// skip 1-byte reserve space
	if _, err := file.Seek(1, io.SeekCurrent); err != nil {
		return nil, fmt.Errorf("failed to seek past reserve space: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &soundData.NumberOfSounds); err != nil {
		return nil, fmt.Errorf("failed to read NumberOfSounds: %w", err)
	}

	// loop through the sound attributes
	soundData.SampleAttrs = make([]*SampleAttributes, soundData.NumberOfSounds)
	for i := range soundData.SampleAttrs {
		if soundData.SampleAttrs[i], err = ReadSampleAttributes(file); err != nil {
			return nil, err
		}
	}

	return soundData, err
}
