package main

import (
	"fmt"

	foo "github.com/parasense/bdmv_go/pkg/mobj"
)

func HeaderPrint(header *foo.MOBJHeader) {
	fmt.Println("Header:")
	PadPrintf(2, "Type: %s\n", string(header.TypeIndicator[:]))
	PadPrintf(2, "Version: %s\n", string(header.VersionNumber[:]))
	PadPrintf(2, "Offset: MovieObject: [%d:%d]\n", header.MovieObjects.Start, header.MovieObjects.Stop)
	PadPrintf(2, "Offset: Extensions: [%d:%d]\n", header.Extensions.Start, header.Extensions.Stop)
	PadPrintln(2, "---")
}

func ClipInfoPrint(movieObjects *foo.MovieObjects) {
	fmt.Println("MovieObjects:")
	PadPrintf(2, "Length: %d\n", movieObjects.Length)
	PadPrintf(2, "NumberOfMovieObjects: %d\n", movieObjects.NumberOfMovieObjects)
	PadPrintln(2, "MovieObjects:")
	for i, mobj := range movieObjects.MovieObjects {
		PadPrintf(4, "MovieObject[%d]:\n", i+1)
		PadPrintf(6, "ResumeIntentionFlag: %t\n", mobj.ResumeIntentionFlag)
		PadPrintf(6, "MenuCallMask: %t\n", mobj.MenuCallMask)
		PadPrintf(6, "TitleSearchMask: %t\n", mobj.TitleSearchMask)
		PadPrintf(6, "NumberOfNavigationCommands: %d\n", mobj.NumberOfNavigationCommands)
		for j, nav := range mobj.NavigationCommands {
			PadPrintf(8, "Navigation Command[%d]:\n", j+1)
			PadPrintf(10, "OperandCount: %d\n", nav.OperandCount)
			PadPrintf(10, "CommandGroup: %d\n", nav.CommandGroup)
			PadPrintf(10, "CommandSubGroup: %d\n", nav.CommandSubGroup)
			PadPrintf(10, "ImmediateValueFlagDest: %t\n", nav.ImmediateValueFlagDest)
			PadPrintf(10, "ImmediateValueFlagSrc: %t\n", nav.ImmediateValueFlagSrc)
			PadPrintf(10, "BranchOption: %d\n", nav.BranchOption)
			PadPrintf(10, "CompareOption: %d\n", nav.CompareOption)
			PadPrintf(10, "SetOption: %d\n", nav.SetOption)
			PadPrintf(10, "Destination: %d\n", nav.Destination)
			PadPrintf(10, "Source: %d\n", nav.Source)
			PadPrintf(8, "CMD: %s %+v, %+v\n",
				foo.GetCommand(
					nav.CommandGroup,
					nav.CommandSubGroup,
					nav.BranchOption,
					nav.CompareOption,
					nav.SetOption,
				), nav.Destination, nav.Source,
			)
			PadPrintln(8, "---")
		}
		PadPrintln(4, "---")
	}
	PadPrintln(2, "---")
}

func ExtensionsPrint(extensions *foo.Extensions) {
	PadPrintln(0)
	PadPrintln(0, "Extensions:")
	PadPrintln(2, "Extensions MetaData:")
	ExtensionsMetaDataPrint(extensions.MetaData)
	PadPrintln(2, "---")
	for i, metaData := range extensions.EntriesMetaData {
		PadPrintf(2, "Extension Entry MetaData [%d]:\n", i+1)
		ExtensionEntryMetaDataPrint(metaData)
		PadPrintln(2, "---")
		if extensions.EntriesData[i] == nil {
			PadPrintln(4, "[Empty Extension payload]")
			continue
		} else {
			PadPrintf(2, "Extension Entry payload [%d]:\n", i+1)
			ExtensionsEntryDataPrint(extensions.EntriesData[i])
		}
		PadPrintln(2, "---")
		PadPrintln(2, "---")
	}
}

func ExtensionsMetaDataPrint(metaData *foo.ExtensionsMetaData) {
	PadPrintf(4, "ExtensionsMetaData.Length: %d\n", metaData.Length)
	PadPrintf(4, "ExtensionsMetaData.EntryDataStartAddr: %d\n", metaData.EntryDataStartAddr)
	PadPrintf(4, "ExtensionsMetaData.EntryDataCount: %d\n", metaData.EntryDataCount)
}

func ExtensionEntryMetaDataPrint(entryMetaData *foo.ExtensionEntryMetaData) {
	PadPrintf(4, "ExtensionEntryMetaData.ExtDataType: %d\n", entryMetaData.ExtDataType)
	PadPrintf(4, "ExtensionEntryMetaData.ExtDataVersion: %d\n", entryMetaData.ExtDataVersion)
	PadPrintf(4, "ExtensionEntryMetaData.ExtDataStartAddress: %d\n", entryMetaData.ExtDataStartAddress)
	PadPrintf(4, "ExtensionEntryMetaData.ExtDataLength: %d\n", entryMetaData.ExtDataLength)
}

func ExtensionsEntryDataPrint(entryData foo.ExtensionEntryData) {
	/*
	   switch entryType := entryData.(type) {

	   case //*clpi.ExtensionLPCMDownMixCoefficient:

	   	//ExtensionLPCMDownMixCoefficientPrint(entryType)

	   case //*clpi.ExtensionExtentStartPoints:

	   	//ExtensionExtentStartPointsPrint(entryType)

	   case //*clpi.ExtensionProgramInfoSS:

	   	//ExtensionProgramInfoSSPrint(entryType)

	   case //*clpi.ExtensionCPISS:

	   		//ExtensionCPISSPrint(entryType)
	   	}
	*/
}
