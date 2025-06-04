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
		return nil, fmt.Errorf("failed to read StreamEntry: %v\n", err)
	}
	stream.Attr, err = ReadStreamAttributes(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read StreamAttributes: %v\n", err)
	}
	return stream, nil
}

func (stream Stream) Print() {
	stream.Entry.Print()
	stream.Attr.Print()
}
