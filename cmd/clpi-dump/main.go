package main

import (
	"fmt"
	"os"
	"strings"

	clpi "github.com/parasense/bdmv_go/pkg/clpi"
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
		fmt.Println("Usage: clpi-parser <clpi-file>")
		os.Exit(1)
	}

	clpiPath := os.Args[1]
	header, clipinfo, sequenceinfo, programinfo, cpi, clipMarks, extensions, err := clpi.ParseCLPI(clpiPath)
	if err != nil {
		fmt.Printf("Error parsing CLPI file: %+v\n", err)
		os.Exit(1)
	}

	PadPrintf(0, "CLPI File: %s\n", clpiPath)
	PadPrintln(0, "")
	HeaderPrint(header)
	PadPrintln(0, "")
	ClipInfoPrint(clipinfo)
	PadPrintln(0, "")
	SequenceInfoPrint(sequenceinfo)
	PadPrintln(0, "")
	ProgramInfoPrint(programinfo)
	PadPrintln(0, "")
	CPIPrint(cpi)
	PadPrintln(0, "")
	ClipMarksPrint(clipMarks)
	PadPrintln(0, "")
	if header.Extensions.Start != 0 && header.Extensions.Stop != 0 {
		ExtensionsPrint(extensions)
	}

}
