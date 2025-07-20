package main

import (
	"fmt"
	"os"
	"strings"

	bclk "github.com/parasense/bdmv_go/pkg/sound"
)

func PadPrintf(indent int, format string, args ...any) {
	fmt.Printf(strings.Repeat(" ", indent)+format, args...)
}

func PadPrintln(indent int, args ...any) {
	fmt.Print(strings.Repeat(" ", indent))
	fmt.Println(args...)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: mobj-dump <mobj-file>")
		os.Exit(1)
	}

	mobjPath := os.Args[1]
	header, soundMetaData, soundData, err := bclk.ParseBCLK(mobjPath)
	if err != nil {
		fmt.Printf("Error parsing CLPI file: %+v\n", err)
		os.Exit(1)
	}

	PadPrintf(0, "BCLK File: %s\n", mobjPath)
	PadPrintln(0, "")
	HeaderPrint(header)
	PadPrintln(0, "")
	SoundMetaDataPrint(soundMetaData)
	PadPrintln(0, "")
	SoundDataPrint(soundData)
	PadPrintln(0, "")
	//PadPrintf(2, "%#v\n\n", soundData)
	//if header.Extensions.Start != 0 && header.Extensions.Stop != 0 {
	//	ExtensionsPrint(extensions)
	//}

}
