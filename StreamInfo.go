package main

import (
	"encoding/binary"
	"fmt"
	"io"
)

type StreamTypeKindOf string

// Stream type constants
const (
	STREAM_TYPE_PRIMARY_VIDEO   StreamTypeKindOf = "PrimaryVideo"
	STREAM_TYPE_PRIMARY_AUDIO   StreamTypeKindOf = "PrimaryAudio"
	STREAM_TYPE_PG              StreamTypeKindOf = "PresentationGraphics"
	STREAM_TYPE_IG              StreamTypeKindOf = "InteractiveGraphics"
	STREAM_TYPE_SECONDARY_AUDIO StreamTypeKindOf = "SecondaryAudio"
	STREAM_TYPE_SECONDARY_VIDEO StreamTypeKindOf = "SecondaryVideo"
	STREAM_TYPE_PIP             StreamTypeKindOf = "PIP"
	STREAM_TYPE_DV              StreamTypeKindOf = "DolbyVision"
)

// streamKinds defines the set of supported stream types
// The order is important!
var streamKinds = []StreamTypeKindOf{
	STREAM_TYPE_PRIMARY_VIDEO,
	STREAM_TYPE_PRIMARY_AUDIO,
	STREAM_TYPE_PG,
	STREAM_TYPE_IG,
	STREAM_TYPE_SECONDARY_AUDIO,
	STREAM_TYPE_SECONDARY_VIDEO,
	STREAM_TYPE_PIP,
	STREAM_TYPE_DV,
}

// StreamItem holds a counter, StreamType label, and streams.
type StreamItem struct {
	NumberOf uint8
	KindOf   StreamTypeKindOf
	Streams  []*Stream
}

// StreamTable contains information about streams in a PlayItem.
type StreamTable struct {
	Length uint16
	Items  []*StreamItem
}

// ReadStreamWrapper populates a slice of Stream pointers from an io.ReaderSeeker.
func ReadStreamWrapper(file io.ReadSeeker, streamItem *StreamItem) (err error) {
	if streamItem.NumberOf != 0 {
		streamItem.Streams = make([]*Stream, streamItem.NumberOf)
		for i := range streamItem.NumberOf {
			if streamItem.Streams[i], err = ReadStream(file, streamItem.KindOf); err != nil {
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

	end, err := CalculateEndOffset(file, streamTable.Length)
	if err != nil {
		return nil, fmt.Errorf("failed calling CalculateEndOffset(): %w", err)
	}

	// XXX
	fmt.Printf("DEBUG: StreamTable.Length: %d\n", streamTable.Length)
	fmt.Printf("DEBUG: StreamTable.end: %d\n", end)

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

	// XXX
	//fmt.Printf("DEBUG: StreamTable(): \n%+v\n\n", streamTable)

	// Seek past reserved 4-byte space
	if _, err := file.Seek(4, io.SeekCurrent); err != nil {
		return nil, fmt.Errorf("failed to seek past reserved space: %w", err)
	}

	// Read the Streams
	for _, item := range streamTable.Items {
		if err := ReadStreamWrapper(file, item); err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", item.KindOf, err)
		}
	}

	// Skip to the end
	if _, err = file.Seek(end, io.SeekStart); err != nil {
		return nil, fmt.Errorf("StreamTable: failed to seek end offset: %w", err)
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
		if item.NumberOf != 0 {
			for j, stream := range item.Streams {
				PadPrintf(8, "%s Stream [%d]:\n", item.KindOf, j+1)
				stream.Print()
				PadPrintln(8, "---")
			}
		} else {
			PadPrintln(8, "[skip]")
		}
	}
}
