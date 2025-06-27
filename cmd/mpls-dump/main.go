package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	mpls "github.com/parasense/bdmv_go/pkg/mpls"
)

func PadPrintf(indent int, format string, args ...any) {
	fmt.Printf(strings.Repeat(" ", indent)+format, args...)
}

func PadPrintln(indent int, args ...any) {
	fmt.Print(strings.Repeat(" ", indent))
	fmt.Println(args...)
}

// Parse45KhzTimestamp converts a timestamp in 45kHz units to a duration
// returns a quantity of time in seconds.
func Parse45KhzTimestamp(timestamp uint32) time.Duration {
	return time.Duration(uint32(timestamp) / 45000)
}

// returns a quantity of time in seconds.
func Convert45KhzTimeToSeconds(timestamp uint32) uint32 {
	return uint32(timestamp / 45000)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: mpls-parser <mpls-file>")
		os.Exit(1)
	}

	mplsPath := os.Args[1]
	header, appinfo, playlist, chapterMarks, extData, err := mpls.ParseMPLS(mplsPath)
	if err != nil {
		fmt.Printf("Error parsing MPLS file: %+v\n", err)
		os.Exit(1)
	}

	PadPrintf(0, "MPLS File: %s\n", mplsPath)
	HeaderPrint(header)
	AppInfoPrint(appinfo)
	PlayListPrint(playlist)
	PlaylistMarksPrint(chapterMarks)
	if header.Extensions.Start != 0 && header.Extensions.Stop != 0 {
		ExtensionsPrint(extData)
	}
}
