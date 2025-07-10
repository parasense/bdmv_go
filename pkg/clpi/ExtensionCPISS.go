package clpi

import (
	"fmt"
	"io"
)

// ExtensionCPISS implements the ExtensionEntryData interface.
// ExtensionCPISS is an alias to CPI.
type ExtensionCPISS = CPI

func (cpi *ExtensionCPISS) Read(file io.ReadSeeker, offsets *OffsetsUint32, entryMeta *ExtensionEntryMetaData) (err error) {

	// Jump to the start offset
	if _, err := file.Seek(offsets.Start+int64(entryMeta.ExtDataStartAddress), io.SeekStart); err != nil {
		return fmt.Errorf("failed to seek Entry Offset Start: (%d); error: %w", entryMeta.ExtDataStartAddress, err)
	}

	// Call to ReadCPI()
	// This entails crafting a custom OffsetsUint32 struct to pass-in
	offset32 := &OffsetsUint32{
		Start: offsets.Start + int64(entryMeta.ExtDataStartAddress),
		Stop:  offsets.Start + int64(entryMeta.ExtDataStartAddress+entryMeta.ExtDataLength),
	}

	result, err := ReadCPI(file, offset32)
	if err != nil {
		return fmt.Errorf("Call to ReadProgramInfo returned error: %w\n", err)
	}

	*cpi = *result

	return nil
}
