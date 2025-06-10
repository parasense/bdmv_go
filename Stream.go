package main

import (
	"fmt"
	"io"
)

type Stream struct {
	Entry StreamEntry
	Attr  StreamAttributes
}

func ReadStream(file io.ReadSeeker) (stream *Stream, err error) {
	stream = &Stream{}
	stream.Entry, err = ReadStreamEntry(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read StreamEntry: %w", err)
	}
	stream.Attr, err = ReadStreamAttributes(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read StreamAttributes: %w", err)
	}
	return stream, nil
}

func (stream Stream) Print() {
	stream.Entry.Print()
	stream.Attr.Print()
}
