package sound

import (
	"fmt"
	"os"
)

/*
	Remarks:
		BDMV/AUXDATA/sound.bdmv
*/

// ParseMPLS parses an MPLS file and returns the playlist details
func ParseBCLK(filePath string) (
	header *BCLKHeader,
	soundMetaData *SoundMetaData,
	soundData *SoundData,
	//extensiondata *Extensions,
	err error,
) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Header
	if header, err = ReadBCLKHeader(file); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to read header: %w", err)
	}

	// Metadata
	if soundMetaData, err = ReadSoundMetaData(file, header.SoundMetaData); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to read appinfo: %w", err)
	}

	// Data
	if soundMetaData.NumberOfSounds > 0 {
		soundData, err = ReadSoundData(file, header.SoundObjects, soundMetaData)
	}

	//// Extensions
	//if header.Extensions.Start != 0 && header.Extensions.Stop != 0 {
	//	if extensiondata, err = ReadExtensions(file, header.Extensions); err != nil {
	//		return nil, nil, nil, fmt.Errorf("failed to read Extension Data: %w", err)
	//	}
	//}

	return header, soundMetaData, soundData, nil
}
