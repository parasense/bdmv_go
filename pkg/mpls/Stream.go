package mpls

import (
	"fmt"
	"io"
)

// Stream represents a media stream in an MPLS file.
// It consists of a StreamEntry and StreamAttributes.
// The StreamEntry contains metadata about the stream, while StreamAttributes
// contains specific attributes related to the stream type.
type Stream struct {
	Entry StreamEntry
	Attr  StreamAttributes
}

// ReadStream reads a Stream from the provided io.ReadSeeker.
// It expects the file (io.ReadSeeker) to be positioned at the start of the Stream structure.
// The Stream consists of a StreamEntry and StreamAttributes.
// It returns a pointer to the Stream and an error if any occurs during reading.
func ReadStream(file io.ReadSeeker, kindOf StreamTypeKindOf) (stream *Stream, err error) {
	stream = &Stream{}

	stream.Entry, err = ReadStreamEntry(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read StreamEntry: %w", err)
	}

	stream.Attr, err = ReadStreamAttributes(file, kindOf)
	if err != nil {
		return nil, fmt.Errorf("failed to read StreamAttributes: %w", err)
	}

	return stream, nil
}

// String returns a string representation of the Stream.
func (s *Stream) String() string {
	return fmt.Sprintf("Stream{Entry: %s, Attr: %s}", s.Entry, s.Attr)
}
