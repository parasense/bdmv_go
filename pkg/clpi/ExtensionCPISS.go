package clpi

import (
	"fmt"
	"io"
)

// ExtensionCPISS implements the ExtensionEntryData interface.
// ExtensionCPISS is an alias to CPI.
type ExtensionCPISS = CPI

func (cpi *ExtensionCPISS) Read(file io.ReadSeeker, offsets *OffsetsUint32, entryMeta *ExtensionEntryMetaData) (err error) {

	// Call to ReadCPI()
	// This entails crafting a custom OffsetsUint32 struct to pass-in.
	// ReadCPI() will jump to the start offset passed in.
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
