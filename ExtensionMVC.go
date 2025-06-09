package main

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Testing:
// Valerian 3D is used for validation.

// TODO:
// Solve dangling 6-byte leftover data at the end.
// * After parsing the video STN there remains mysterious 6-bytes.
//   Of which, one byte is non-zero that indicating non-reserve padding)

// MVC (Multi View Coding) is used for 3D
type ExtensionMVCStream struct {
	numberOfItems           uint32    // XXX - not parsed, caclulated in a bad way
	LengthOfItems           uint16    //
	FixedOffsetPopUpFlag    uint8     // 0x0=false or 0x80=true
	MVCStreams              []*Stream //
	NumberOfOffsetSequences uint8     // up to 32 e.g 0x20
	remainingBytes          [6]byte   // XXX
}

func (mvcStreamExtension *ExtensionMVCStream) Read(file io.ReadSeeker) (err error) {

	if err := binary.Read(file, binary.BigEndian, &mvcStreamExtension.LengthOfItems); err != nil {
		return err
	}

	if err := binary.Read(file, binary.BigEndian, &mvcStreamExtension.FixedOffsetPopUpFlag); err != nil {
		return err
	}

	// 1-byte reserve space
	if _, err := file.Seek(1, io.SeekCurrent); err != nil {
		return err
	}

	// numberOfItems is for the loop control variable for reading STNs
	mvcStreamExtension.numberOfItems = mvcStreamExtension.numberOfItems / uint32(mvcStreamExtension.LengthOfItems)

	// TODO - fix error handling

	// StereoScopic (3D) Video Stream
	if err := ReadStreamWrapper(file, uint8(mvcStreamExtension.numberOfItems), &mvcStreamExtension.MVCStreams); err != nil {
		//return nil, err
	}

	file.Seek(1, io.SeekCurrent) // skip 1-byte reserve space
	binary.Read(file, binary.BigEndian, &mvcStreamExtension.NumberOfOffsetSequences)

	tmp, _ := ftell(file)
	fmt.Printf("mvcStreamExtension: End of 3d Stream ftell: %d \n\n", tmp)

	// XXX - This is a hack
	//       There appears to a data structure in the remaining 6-bytes.
	//       The 3rd byte is set to 0x01, the rest are all zero.
	// TODO: Remove this when the mystery is solved, and there is a legitimate struct.
	binary.Read(file, binary.BigEndian, &mvcStreamExtension.remainingBytes)

	// XXX - there should not be any remaining bytes.
	//     - The 3rd byte is set as "1", so it cannot be reserve space.
	//     - Maybe a PG structure? (no idea)
	fmt.Printf("REMAINING BYTES: %+v \n\n", mvcStreamExtension.remainingBytes)

	// This is a hack... remove in the future.
	// It shows when the file pointer is at the end of the structure.
	tmp, _ = ftell(file)
	fmt.Printf("mvcStreamExtension: End of 3d Stream ftell: %d \n\n", tmp)

	return nil

}

// XXX - I hate this - there must be another way to get the NumberOf?
//   - Speculation: There might be an assumption one only one item.
func (mvcStreamExtension *ExtensionMVCStream) SetNumberOfItems(length uint32) {
	mvcStreamExtension.numberOfItems = (length - 2) / uint32(mvcStreamExtension.LengthOfItems)
}

func (mvcStreamExtension *ExtensionMVCStream) Print() {
	PadPrintln(4, "MVC(3D)Extension:")
	PadPrintf(6, "numberOfItems: %d\n", mvcStreamExtension.numberOfItems)
	PadPrintf(6, "LengthOfItems: %d\n", mvcStreamExtension.LengthOfItems)
	PadPrintf(6, "FixedOffsetPopUpFlag: %d\n", mvcStreamExtension.FixedOffsetPopUpFlag)
	for i := uint32(0); i < mvcStreamExtension.numberOfItems; i++ {
		mvcStreamExtension.MVCStreams[i].Print()
	}
	PadPrintf(6, "NumberOfOffsetSequences: %d\n", mvcStreamExtension.NumberOfOffsetSequences)
}
