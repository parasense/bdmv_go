package mpls

import (
	"encoding/binary"
	"fmt"
	"io"
)

// PlayItemEntry represents a play item entry in an MPLS file.
// It consists of a 5-byte clip name, a 4-byte codec identifier,
// and a 1-byte reference to the Stream Table ID (STCID).
// The PlayItemEntry is used to identify a specific media file and its codec,
// along with a reference to the Stream Table ID that contains additional information
// about the media stream associated with this play item.
// The FileName is a 5-byte string that typically represents the name of the media file
// (e.g., "00063").
// The Codec is a 4-byte string that identifies the codec used for the media file
// (e.g., "M2TS").
type PlayItemEntry struct {
	FileName   [5]byte
	Codec      [4]byte
	RefToSTCID uint8
}

// ReadPlayItemEntry reads a PlayItemEntry from the provided io.ReadSeeker.
// It expects the file (io.ReadSeeker) to be positioned at the start of the PlayItemEntry structure.
// The PlayItemEntry consists of a 5-byte clip name, a 4-byte codec identifier,
// and a 1-byte reference to the Stream Table ID (STCID).
// It returns a pointer to the PlayItemEntry and an error if any occurs during reading.
func ReadPlayItemEntry(file io.ReadSeeker) (*PlayItemEntry, error) {
	playItemEntry := &PlayItemEntry{}

	// The 5 bytes clip name
	if err := binary.Read(file, binary.BigEndian, &playItemEntry.FileName); err != nil {
		return nil, fmt.Errorf("failed to read clip info filename: %w", err)
	}

	// The 4 byte codec should be something like "M2TS"
	if err := binary.Read(file, binary.BigEndian, &playItemEntry.Codec); err != nil {
		return nil, fmt.Errorf("failed to read clip codec identifier: %w", err)
	}

	if err := binary.Read(file, binary.BigEndian, &playItemEntry.RefToSTCID); err != nil {
		return nil, fmt.Errorf("failed to read play item length: %w", err)
	}

	return playItemEntry, nil
}

func (p *PlayItemEntry) String() string {
	return fmt.Sprintf("PlayItemEntry{FileName: %s, Codec: %s, RefToSTCID: %d}",
		p.FileName, p.Codec, p.RefToSTCID)
}
