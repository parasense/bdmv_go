package main

import (
	"fmt"
	"os"
	"strings"

	indx "github.com/parasense/bdmv_go/pkg/indx"
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
		fmt.Println("Usage: indx-parser <indx-file>")
		os.Exit(1)
	}

	indxPath := os.Args[1]
	header, appinfo, indexes, extData, err := indx.ParseINDX(indxPath)
	if err != nil {
		fmt.Printf("Error parsing INDX file: %+v\n", err)
		os.Exit(1)
	}

	PadPrintf(0, "INDX File: %s\n", indxPath)

	if header.Extensions.Start != 0 && header.Extensions.Stop != 0 {
		PadPrintf(0, "Extensions: %s\n", extData)
	}
	HeaderPrint(header)
	AppInfoPrint(appinfo)
	IndexesPrint(indexes)
	//if header.Extensions.Start != 0 && header.Extensions.Stop != 0 {
	//	//extData.Print()
	//	ExtensionsPrint(extData)
	//}

}
