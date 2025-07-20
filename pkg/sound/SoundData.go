package sound

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Linear Pulse Code Modulation with 16-bits.
// codec for modeling audio oscilations.
type LPCM16 uint16
type Samples []LPCM16
type SoundData struct {
	Data []*Samples
}

// Reads ALL sounds from the given array of sound attribute entries.
func ReadSoundData(file io.ReadSeeker, offsets *OffsetsUint32, soundMetaData *SoundMetaData) (soundData *SoundData, err error) {
	soundData = &SoundData{}

	soundData.Data = make([]*Samples, soundMetaData.NumberOfSounds)
	for i, attr := range soundMetaData.SampleAttrs {
		fmt.Printf("attr[%d]: %+v\n", i, attr)
		soundData.Data[i], err = ReadSoundDataBlock(file, offsets, attr)
	}

	return soundData, err
}

// Reads ONE sound from the given sound attribute entry.
func ReadSoundDataBlock(file io.ReadSeeker, offsets *OffsetsUint32, attr *SampleAttributes) (data *Samples, err error) {

	// Jump to the start address
	if _, err := file.Seek(offsets.Start+int64(attr.SoundDataIndex), io.SeekStart); err != nil {
		return nil, fmt.Errorf("failed to seek to SoundDataIndex: %w", err)
	}

	data = &Samples{}
	*(data) = make(Samples, attr.NumberOfFrames*uint32(attr.NumberOfChannels))

	for i := range *data {
		if err := binary.Read(file, binary.BigEndian, &(*data)[i]); err != nil {
			return nil, fmt.Errorf("failed to read sample data: %w", err)
		}
	}

	return data, err
}
