package main

import (
	"fmt"

	bclk "github.com/parasense/bdmv_go/pkg/sound"
)

func HeaderPrint(header *bclk.BCLKHeader) {
	fmt.Println("Header:")
	PadPrintf(2, "Type: %s\n", string(header.TypeIndicator[:]))
	PadPrintf(2, "Version: %s\n", string(header.VersionNumber[:]))
	PadPrintf(2, "Offset: SoundMetaData: [%d:%d]\n", header.SoundMetaData.Start, header.SoundMetaData.Stop)
	PadPrintf(2, "Offset: SoundObjects: [%d:%d]\n", header.SoundObjects.Start, header.SoundObjects.Stop)
	PadPrintf(2, "Offset: Extensions: [%d:%d]\n", header.Extensions.Start, header.Extensions.Stop)
	PadPrintln(2, "---")
}

func SoundMetaDataPrint(soundData *bclk.SoundMetaData) {
	fmt.Println("SoundMetaData:")
	PadPrintf(2, "Length: %d\n", soundData.Length)
	PadPrintf(2, "NumberOfSounds: %d\n", soundData.NumberOfSounds)
	PadPrintln(2, "Sounds:")
	for i, sound := range soundData.SampleAttrs {
		PadPrintf(4, "Sound[%d]:\n", i+1)
		PadPrintf(6, "NumberOfChannels: %d\n", sound.NumberOfChannels)
		PadPrintf(6, "SampleRate: %d Hz\n", sound.SampleRate)
		PadPrintf(6, "BitsPerSample: %d\n", sound.BitsPerSample)
		PadPrintf(6, "SoundDataIndexes: %d\n", sound.SoundDataIndex)
		PadPrintf(6, "NumberOfFrames(perChannel): %d\n", sound.NumberOfFrames)
		PadPrintf(5, "*Duration: %.3f Seconds\n", float32(sound.NumberOfFrames)/float32(sound.SampleRate))
		PadPrintf(5, "*Size: %d Bytes\n", (sound.NumberOfFrames*uint32(sound.BitsPerSample)/8)*uint32(sound.NumberOfChannels))
		PadPrintln(4, "---")
	}
}

func SoundDataPrint(soundData *bclk.SoundData) {
	fmt.Println("SoundData:")
	for i, soundData := range soundData.Data {
		PadPrintf(2, "Data[%d]: (len: %d)\n", i, len(*soundData))
	}
}
