package main

import (
	"fmt"

	indx "github.com/parasense/bdmv_go/pkg/indx"
)

func HeaderPrint(header *indx.INDXHeader) {
	fmt.Println("Header:")
	PadPrintf(2, "Type: %s\n", string(header.TypeIndicator[:]))
	PadPrintf(2, "Version: %s\n", string(header.VersionNumber[:]))
	PadPrintf(2, "Offset: AppInfo: [%d:%d]\n", header.AppInfo.Start, header.AppInfo.Stop)
	PadPrintf(2, "Offset: Indexes: [%d:%d]\n", header.Indexes.Start, header.Indexes.Stop)
	PadPrintf(2, "Offset: Extensions: [%d:%d]\n", header.Extensions.Start, header.Extensions.Stop)
	PadPrintln(2, "---")
}

func AppInfoPrint(appinfo *indx.AppInfo) {
	PadPrintln(0, "AppInfo:")
	PadPrintf(2, "Length: %d\n", appinfo.Length)
	PadPrintf(2, "InitialOutputModePreference: %d\n", appinfo.InitialOutputModePreference)
	PadPrintf(2, "SSContentExistFlag: %t\n", appinfo.SSContentExistFlag)
	PadPrintf(2, "InitialDynamicRangeType: %v\n", appinfo.InitialDynamicRangeType)
	PadPrintf(2, "VideoFormat: %v\n", appinfo.VideoFormat)
	PadPrintf(2, "FrameRate: %v\n", appinfo.FrameRate)
	PadPrintf(2, "UserData: %v\n", appinfo.UserData)
	PadPrintln(2, "---")
}

func IndexesPrint(indexes *indx.Indexes) {
	PadPrintln(0, "Indexes:")
	PadPrintf(2, "Length: %d\n", indexes.Length)
	PadPrintf(2, "FirstPlaybackTitle: \n")
	TitlePrint(indexes.FirstPlaybackTitle)
	PadPrintf(2, "TopMenuTitle: \n")
	TitlePrint(indexes.TopMenuTitle)
	PadPrintf(2, "NumberOfTitles: %d\n", indexes.NumberOfTitles)
	PadPrintf(2, "Titles: \n")
	for i, title := range indexes.Titles {
		PadPrintf(4, "Title[%d]: \n", i+1)
		TitlePrint(title)
	}
}

func TitlePrint(title *indx.Title) {
	PadPrintln(2, "Title:")
	PadPrintf(4, "ObjectType: %d\n", title.ObjectType)
	PadPrintf(4, "AccesType: %d\n", title.AccesType)
	PadPrintf(4, "PlaybackType: %d\n", title.PlaybackType)

	switch title.PlaybackType {
	case 0: // hdmv presentation
		PadPrintf(4, "RefToMovieObjectID: %v\n", title.RefToMovieObjectID)
	case 1: // hdmv interactive
		PadPrintf(4, "RefToMovieObjectID: %v\n", title.RefToMovieObjectID)
	case 2: // bdj presentation
		PadPrintf(4, "RefToBDJObjectID: %v\n", title.RefToBDJObjectID)
	case 3: // bdj interactive
		PadPrintf(4, "RefToBDJObjectID: %v\n", title.RefToBDJObjectID)
	}

	PadPrintln(2, "---")
}
