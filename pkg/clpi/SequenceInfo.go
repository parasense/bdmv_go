package clpi

import (
	"encoding/binary"
	"fmt"
	"io"
)

type SequenceInfo struct {
	Length               uint32
	NumberOfATCSequences uint8
	ATCSequences         []*ATCSequence
}

type ATCSequence struct {
	SPNATCStart          uint32
	NumberOfSTCSequences uint8
	OffsetSTCID          uint8
	STCSequences         []*STCSequence
}

type STCSequence struct {
	PCRPID                uint16
	SPNSTCStart           uint32
	PresentationStartTime uint32
	PresentationEndTime   uint32
}

func ReadSequenceInfo(file io.ReadSeeker, offsets *OffsetsUint32) (sequenceInfo *SequenceInfo, err error) {
	sequenceInfo = &SequenceInfo{}

	// Jump to start address
	if _, err := file.Seek(offsets.Start, io.SeekStart); err != nil {
		return nil, err
	}

	if err := binary.Read(file, binary.BigEndian, &sequenceInfo.Length); err != nil {
		return nil, err
	}

	// 1-byte reserve space
	if _, err := file.Seek(1, io.SeekCurrent); err != nil {
		return nil, err
	}

	if err := binary.Read(file, binary.BigEndian, &sequenceInfo.NumberOfATCSequences); err != nil {
		return nil, err
	}

	sequenceInfo.ATCSequences = make([]*ATCSequence, sequenceInfo.NumberOfATCSequences)
	for i := range sequenceInfo.ATCSequences {
		if sequenceInfo.ATCSequences[i], err = ReadATCSequence(file); err != nil {
			return nil, err
		}
	}

	return sequenceInfo, nil
}

func ReadATCSequence(file io.ReadSeeker) (atcSequence *ATCSequence, err error) {
	atcSequence = &ATCSequence{}

	if err := binary.Read(file, binary.BigEndian, &atcSequence.SPNATCStart); err != nil {
		return nil, err
	}

	if err := binary.Read(file, binary.BigEndian, &atcSequence.NumberOfSTCSequences); err != nil {
		return nil, err
	}

	if err := binary.Read(file, binary.BigEndian, &atcSequence.OffsetSTCID); err != nil {
		return nil, err
	}

	atcSequence.STCSequences = make([]*STCSequence, atcSequence.NumberOfSTCSequences)
	for i := range atcSequence.STCSequences {
		atcSequence.STCSequences[i], err = ReadSTCSequences(file)
	}

	return atcSequence, nil
}

func ReadSTCSequences(file io.ReadSeeker) (stcSequence *STCSequence, err error) {
	stcSequence = &STCSequence{}

	if err := binary.Read(file, binary.BigEndian, &stcSequence.PCRPID); err != nil {
		return nil, err
	}

	if err := binary.Read(file, binary.BigEndian, &stcSequence.SPNSTCStart); err != nil {
		return nil, err
	}

	if err := binary.Read(file, binary.BigEndian, &stcSequence.PresentationStartTime); err != nil {
		return nil, err
	}

	if err := binary.Read(file, binary.BigEndian, &stcSequence.PresentationEndTime); err != nil {
		return nil, err
	}
	return stcSequence, nil
}

func (sequenceInfo *SequenceInfo) String() string {
	return "SequenceInfo{" +
		"Length: " + fmt.Sprintf("%d", sequenceInfo.Length) +
		", NumberOfATCSequences: " + fmt.Sprintf("%d", sequenceInfo.NumberOfATCSequences) +
		", ATCSequences: " + fmt.Sprintf("%v", sequenceInfo.ATCSequences) +
		"}"
}

func (atcSequence *ATCSequence) String() string {
	return "ATCSequence{" +
		"SPNATCStart: " + fmt.Sprintf("%d", atcSequence.SPNATCStart) +
		", NumberOfSTCSequences: " + fmt.Sprintf("%d", atcSequence.NumberOfSTCSequences) +
		", OffsetSTCID: " + fmt.Sprintf("%d", atcSequence.OffsetSTCID) +
		", STCSequences: " + fmt.Sprintf("%v", atcSequence.STCSequences) +
		"}"
}

func (stcSequence *STCSequence) String() string {
	return fmt.Sprintf(
		"STCSequence{PCRPID: %d, "+
			"SPNSTCStart: %d, "+
			"PresentationStartTime: %d, P"+
			"resentationEndTime: %d}",
		stcSequence.PCRPID, stcSequence.SPNSTCStart,
		stcSequence.PresentationStartTime, stcSequence.PresentationEndTime,
	)
}
