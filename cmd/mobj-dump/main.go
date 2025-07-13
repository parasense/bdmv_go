package main

import (
	"fmt"
	"os"
	"strings"

	mobj "github.com/parasense/bdmv_go/pkg/mobj"
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
	header, movieObjects, extensions, err := mobj.ParseMOBJ(mobjPath)
	if err != nil {
		fmt.Printf("Error parsing CLPI file: %+v\n", err)
		os.Exit(1)
	}

	PadPrintf(0, "MOBJ File: %s\n", mobjPath)
	PadPrintln(0, "")
	HeaderPrint(header)
	PadPrintln(0, "")
	ClipInfoPrint(movieObjects)
	PadPrintln(0, "")
	if header.Extensions.Start != 0 && header.Extensions.Stop != 0 {
		ExtensionsPrint(extensions)
	}

}
