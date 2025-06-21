package main

import (
	"encoding/binary"
	"fmt"
	"io"
)

// ExtensionSubPath implements the ExtensionEntryData interface.
type ExtensionSubPath struct {
	Length   uint32
	Count    uint16
	SubPaths []*SubPath
}

func (extensionSubPath *ExtensionSubPath) Read(file io.ReadSeeker, offsets *OffsetsUint32, entryMeta *ExtensionEntryMetaData) (err error) {

	PadPrintln(0, "SubPath Extension:")
	PadPrintln(2, "---")

	// Calculate the Start/Stop offsets for this extension.
	offsetStart := offsets.Start + int64(entryMeta.ExtDataStartAddress)
	offsetStop := offsetStart + int64(entryMeta.ExtDataLength)
	PadPrintf(2, "offsetStart == %d\n", offsetStart)
	PadPrintf(2, "offsetStop  == %d\n", offsetStop)
	PadPrintln(2, "---")

	// Jump to the start offset
	if _, err := file.Seek(offsetStart, io.SeekStart); err != nil {
		return fmt.Errorf("failed to seek Entry Offset Start: (%d); error: %w", offsetStart, err)
	}

	// XXX - DEBUG block
	PadPrintln(0, "Extensions Entry DEBUG:")
	PadPrintf(2, "ExtDataType == %d\n", entryMeta.ExtDataType)
	PadPrintf(2, "ExtDataVersion == %d\n", entryMeta.ExtDataVersion)
	PadPrintf(2, "ExtDataStartAddress == %d\n", entryMeta.ExtDataStartAddress)
	PadPrintf(2, "ExtDataLength == %d\n", entryMeta.ExtDataLength)
	fmt.Println("---")
	// XXX - EO DEBUG block

	if err := binary.Read(file, binary.BigEndian, &extensionSubPath.Length); err != nil {
		return fmt.Errorf("failed to read extensionSubPath.Length: %w", err)
	}
	PadPrintf(4, "extensionSubPath.Length: %d\n", extensionSubPath.Length)

	if err := binary.Read(file, binary.BigEndian, &extensionSubPath.Count); err != nil {
		return fmt.Errorf("failed to read extensionSubPath.Count: %w", err)
	}
	PadPrintf(4, "extensionSubPath.Count: %d\n", extensionSubPath.Count)

	extensionSubPath.SubPaths = make([]*SubPath, extensionSubPath.Count)
	for i := range extensionSubPath.SubPaths {
		PadPrintf(4, "extensionSubPath.SubPaths[%d]\n", i)
		if extensionSubPath.SubPaths[i], err = ReadSubPath(file); err != nil {
			return fmt.Errorf("failed calling ReadSubPath() in ExtensionSubPath.Read(): %w", err)
		}
		extensionSubPath.SubPaths[i].Print()

		//// 1-byte reserve space
		//if _, err := file.Seek(end, io.SeekStart); err != nil {
		//	return fmt.Errorf("failed to seek past reserve space: %w", err)
		//}
	}

	return nil
}

func (subPathExtension *ExtensionSubPath) Print() {
	PadPrintln(4, "SubPathExtension")
	PadPrintf(6, "Length: %d\n", subPathExtension.Length)
	PadPrintf(6, "Count: %d\n", subPathExtension.Count)
	for i, subPath := range subPathExtension.SubPaths {
		PadPrintf(6, "SubPath [%d]:\n", i+1)
		subPath.Print()
	}
}
