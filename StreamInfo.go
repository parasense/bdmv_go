package main

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Stream type constants
const (
	STREAM_TYPE_VIDEO           = "Video"
	STREAM_TYPE_AUDIO           = "Audio"
	STREAM_TYPE_PG              = "PresentationGraphics"
	STREAM_TYPE_IG              = "InteractiveGraphics"
	STREAM_TYPE_SECONDARY_AUDIO = "SecondaryAudio"
	STREAM_TYPE_SECONDARY_VIDEO = "SecondaryVideo"
	STREAM_TYPE_PIP             = "PIP"
	STREAM_TYPE_DV              = "DolbyVision"
)

// streamKinds defines the supported stream types
var streamKinds = []string{
	STREAM_TYPE_VIDEO,
	STREAM_TYPE_AUDIO,
	STREAM_TYPE_PG,
	STREAM_TYPE_IG,
	STREAM_TYPE_SECONDARY_AUDIO,
	STREAM_TYPE_SECONDARY_VIDEO,
	STREAM_TYPE_PIP,
	STREAM_TYPE_DV,
}

// StreamItem holds a counter, type, and streams.
type StreamItem struct {
	NumberOf uint8
	KindOf   string
	Streams  []*Stream
}

// StreamTable contains information about streams in a PlayItem.
type StreamTable struct {
	Length uint16
	Items  []*StreamItem
}

// ReadStreamWrapper populates a slice of Stream pointers from an io.ReaderSeeker.
func ReadStreamWrapper(file io.ReadSeeker, count uint8, streamPointer *[]*Stream) (err error) {
	if count != 0 {
		*streamPointer = make([]*Stream, count)
		for i := uint8(0); i < count; i++ {
			if (*streamPointer)[i], err = ReadStream(file); err != nil {
				return fmt.Errorf("failed to read Stream: %w", err)
			}
		}
	}
	return nil
}

func ReadStreamTable(file io.ReadSeeker) (*StreamTable, error) {
	// Initialize StreamTable with pre-allocated StreamItems
	streamTable := &StreamTable{
		Items: make([]*StreamItem, len(streamKinds)),
	}
	for i, kind := range streamKinds {
		streamTable.Items[i] = &StreamItem{KindOf: kind}
	}

	// Read Length
	if err := binary.Read(file, binary.BigEndian, &streamTable.Length); err != nil {
		return nil, fmt.Errorf("failed to read stream table length: %w", err)
	}

	// Seek past reserved 2-byte space
	if _, err := file.Seek(2, io.SeekCurrent); err != nil {
		return nil, fmt.Errorf("failed to seek past reserved space: %w", err)
	}

	// Read the counter fields
	for _, item := range streamTable.Items {
		if err := binary.Read(file, binary.BigEndian, &item.NumberOf); err != nil {
			return nil, fmt.Errorf("failed to read NumberOf%s: %w", item.KindOf, err)
		}
	}

	// Seek past reserved 4-byte space
	if _, err := file.Seek(4, io.SeekCurrent); err != nil {
		return nil, fmt.Errorf("failed to seek past reserved space: %w", err)
	}

	// Read the Streams
	for _, item := range streamTable.Items {
		if err := ReadStreamWrapper(file, item.NumberOf, &item.Streams); err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", item.KindOf, err)
		}
	}

	// Validate the StreamTable
	if err := streamTable.Assert(); err != nil {
		return nil, fmt.Errorf("stream table validation failed: %w", err)
	}

	return streamTable, nil
}

func (streamTable *StreamTable) Assert() error {
	// Check if Length is zero
	if streamTable.Length == 0 {
		return fmt.Errorf("streamTable.Length cannot be zero")
	}

	// Check if there is at least one stream
	var totalStreams uint8
	for _, item := range streamTable.Items {
		totalStreams += item.NumberOf
	}
	if totalStreams == 0 {
		return fmt.Errorf("At least one stream must be present")
	}

	// Validate slice lengths against their corresponding NumberOf fields
	for _, item := range streamTable.Items {
		if len(item.Streams) != int(item.NumberOf) {
			return fmt.Errorf("%s slice length (%d) does not match NumberOf%s (%d)", item.KindOf, len(item.Streams), item.KindOf, item.NumberOf)
		}
	}

	return nil
}

func (streamTable *StreamTable) Print() {
	PadPrintln(4, "StreamTable:")
	PadPrintf(6, "Length: %d\n", streamTable.Length)

	for _, item := range streamTable.Items {
		PadPrintf(6, "NumberOf%s: %d\n", item.KindOf, item.NumberOf)
		for j, stream := range item.Streams {
			PadPrintf(6, "%s Stream [%d]:\n", item.KindOf, j)
			stream.Print()
		}
	}
}
