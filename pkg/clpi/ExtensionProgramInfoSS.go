package clpi

import (
	"fmt"
	"io"
)

// ExtensionProgramInfoSS implements the ExtensionEntryData interface.
// ExtensionProgramInfoSS is an alias to ProgramInfo.
type ExtensionProgramInfoSS = ProgramInfo

func (pi *ExtensionProgramInfoSS) Read(file io.ReadSeeker, offsets *OffsetsUint32, entryMeta *ExtensionEntryMetaData) (err error) {

	// Call to ReadProgramInfo()
	// This entails crafting a custom OffsetsUint32 struct to pass-in.
	// ReadProgramInfo() will jump to the start offset passed in.
	offset32 := &OffsetsUint32{
		Start: offsets.Start + int64(entryMeta.ExtDataStartAddress),
		Stop:  offsets.Start + int64(entryMeta.ExtDataStartAddress+entryMeta.ExtDataLength),
	}

	result, err := ReadProgramInfo(file, offset32)
	if err != nil {
		return fmt.Errorf("Call to ReadProgramInfo returned error: %w\n", err)
	}

	*pi = *result

	return nil
}
