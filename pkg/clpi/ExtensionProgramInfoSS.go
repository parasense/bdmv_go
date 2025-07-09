package clpi

import (
	"fmt"
	"io"
)

// ExtensionProgramInfoSS implements the ExtensionEntryData interface.
// ExtensionProgramInfoSS is an alias to ProgramInfo.
type ExtensionProgramInfoSS = ProgramInfo

func (pi *ExtensionProgramInfoSS) Read(file io.ReadSeeker, offsets *OffsetsUint32, entryMeta *ExtensionEntryMetaData) (err error) {

	// Jump to the start offset
	if _, err := file.Seek(offsets.Start+int64(entryMeta.ExtDataStartAddress), io.SeekStart); err != nil {
		return fmt.Errorf("failed to seek Entry Offset Start: (%d); error:  %w", offsets.Start+int64(entryMeta.ExtDataStartAddress), err)
	}

	// Call to ReadProgramInfo()
	// This entails crafting a custom OffsetsUint32 struct to pass-in
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
