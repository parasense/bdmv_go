package main

import (
	"fmt"

	"github.com/parasense/bdmv_go/pkg/mpls"
)

func HeaderPrint(header *mpls.MPLSHeader) {
	fmt.Println("Header:")
	PadPrintf(2, "Type: %s\n", string(header.TypeIndicator[:]))
	PadPrintf(2, "Version: %s\n", string(header.VersionNumber[:]))
	PadPrintf(2, "Offset: AppInfo: [%d:%d]\n", header.AppInfo.Start, header.AppInfo.Stop)
	PadPrintf(2, "Offset: PlayList: [%d:%d]\n", header.Playlist.Start, header.Playlist.Stop)
	PadPrintf(2, "Offset: Marks: [%d:%d]\n", header.Marks.Start, header.Marks.Stop)
	PadPrintf(2, "Offset: Extensions: [%d:%d]\n", header.Extensions.Start, header.Extensions.Stop)
	PadPrintln(2, "---")
}

func AppInfoPrint(appinfo *mpls.AppInfo) {
	PadPrintln(0, "AppInfo:")
	PadPrintf(2, "Length: %d\n", appinfo.Length)
	PadPrintf(2, "PlaybackType: %d\n", appinfo.PlaybackType)
	PadPrintf(2, "PlaybackCount: %d\n", appinfo.PlaybackCount)
	USerOptionsPrint(appinfo.UserOptions)
	PadPrintf(2, "RandomAccessFlag: %v\n", appinfo.RandomAccessFlag)
	PadPrintf(2, "AudioMixFlag: %v\n", appinfo.AudioMixFlag)
	PadPrintf(2, "LosslessBypassFlag: %v\n", appinfo.LosslessBypassFlag)
	PadPrintf(2, "MVCBaseViewRFlag: %v\n", appinfo.MVCBaseViewRFlag)
	PadPrintf(2, "SDRConversionNotificationFlag: %v\n", appinfo.SDRConversionNotificationFlag)
	PadPrintln(2, "---")
}

func PlayListPrint(playlist *mpls.PlayList) {
	PadPrintln(0, "PlayList:")
	PadPrintf(2, "Length: %d\n", playlist.Length)
	PadPrintf(2, "NumberOfPlayItems: %d\n", playlist.NumberOfPlayItems)
	PadPrintf(2, "NumberOfSubPaths: %d\n", playlist.NumberOfSubPaths)
	PadPrintln(2)

	var totalDuration uint32
	for i, playItem := range playlist.PlayItems {
		inTime := Convert45KhzTimeToSeconds(playItem.INTime)
		outTime := Convert45KhzTimeToSeconds(playItem.OUTTime)
		duration := outTime - inTime
		totalDuration += duration

		PadPrintf(2, "PlayItem [%d]:\n", i)
		PlayItemPrint(playItem)
		PadPrintln(2, "---")
	}

	for i, subPath := range playlist.SubPaths {
		PadPrintf(2, "SubPath [%d]:\n", i)
		SubPathPrint(subPath)
		PadPrintln(2, "---")
	}
	PadPrintln(0, "---")
}

func PlayItemPrint(playItem *mpls.PlayItem) {
	inTime := Convert45KhzTimeToSeconds(playItem.INTime)
	outTime := Convert45KhzTimeToSeconds(playItem.OUTTime)
	PadPrintf(4, "Length: %d\n", playItem.Length)
	PadPrintf(4, "Clip File: %s\n", playItem.ClipInformationFileName)
	PadPrintf(4, "Codec ID: %s\n", playItem.ClipCodecIdentifier)
	PadPrintf(4, "Multi-angle: %v\n", playItem.IsMultiAngle)
	PadPrintf(4, "ConnectionCondition: %d\n", playItem.ConnectionCondition)
	PadPrintf(4, "RefToSTCID: %d\n", playItem.RefToSTCID)
	PadPrintf(4, "InTime: %d (%d)\n", playItem.INTime, inTime)
	PadPrintf(4, "OUTime: %d (%d)\n", playItem.OUTTime, outTime)
	PadPrintf(6, "*Duration: %v\n", outTime-inTime)
	USerOptionsPrint(playItem.UserOptions)
	PadPrintf(4, "PlayItemRandomAccessFlag: %v\n", playItem.PlayItemRandomAccessFlag)
	PadPrintf(4, "StillMode: %v\n", playItem.StillMode)
	PadPrintf(4, "StillTime: %v\n", playItem.StillTime)
	PadPrintf(4, "NumberOfAngles: %v\n", playItem.NumberOfAngles)
	PadPrintf(4, "IsDifferentAudios: %v\n", playItem.IsDifferentAudios)
	PadPrintf(4, "IsSeamlessAngleChange: %v\n", playItem.IsSeamlessAngleChange)
	PadPrintf(4, "Angles:\n")
	for i := range playItem.Angles {
		PadPrintf(6, "Angle: [%d]:\n", i+1)
		PlayItemEntryPrint(playItem.Angles[i])
	}
	StreamTablePrint(playItem.StreamTable)
}

func SubPathPrint(subPath *mpls.SubPath) {
	PadPrintf(4, "Length: %d\n", subPath.Length)
	PadPrintf(4, "SubPathType: %d [%s]\n", subPath.SubPathType, mpls.SubPathType(subPath.SubPathType))
	PadPrintf(4, "IsRepeatSubPath: %v\n", subPath.IsRepeatSubPath)
	PadPrintf(4, "NumberOfSubPlayItems: %d\n", subPath.NumberOfSubPlayItems)
	PadPrintf(4, "SubPlayItems:\n")
	for i := range subPath.SubPlayItems {
		PadPrintf(6, "Angle [%d]:\n", i)
		SubPlayItemPrint(subPath.SubPlayItems[i])
	}
}

func USerOptionsPrint(userOptions *mpls.UserOptions) {
	PadPrintln(4, "UserOptions:")
	PadPrintf(6, "MenuCall: %v\n", userOptions.MenuCall)
	PadPrintf(6, "TitleSearch: %v\n", userOptions.TitleSearch)
	PadPrintf(6, "ChapterSearch: %v\n", userOptions.ChapterSearch)
	PadPrintf(6, "TimeSearch: %v\n", userOptions.TimeSearch)
	PadPrintf(6, "SkipToNextPoint: %v\n", userOptions.SkipToNextPoint)
	PadPrintf(6, "SkipToPrevPoint: %v\n", userOptions.SkipToPrevPoint)
	PadPrintf(6, "Stop: %v\n", userOptions.Stop)
	PadPrintf(6, "PauseOn: %v\n", userOptions.PauseOn)
	PadPrintf(6, "StillOff %v\n", userOptions.StillOff)
	PadPrintf(6, "ForwardPlay: %v\n", userOptions.ForwardPlay)
	PadPrintf(6, "BackwardPlay: %v\n", userOptions.BackwardPlay)
	PadPrintf(6, "Resume: %v\n", userOptions.Resume)
	PadPrintf(6, "MoveUpSelectedButton: %v\n", userOptions.MoveUpSelectedButton)
	PadPrintf(6, "MoveDownSelectedButton: %v\n", userOptions.MoveDownSelectedButton)
	PadPrintf(6, "MoveLeftSelectedButton: %v\n", userOptions.MoveLeftSelectedButton)
	PadPrintf(6, "MoveRightSelectedButton: %v\n", userOptions.MoveRightSelectedButton)
	PadPrintf(6, "SelectButton: %v\n", userOptions.SelectButton)
	PadPrintf(6, "ActivateButton: %v\n", userOptions.ActivateButton)
	PadPrintf(6, "SelectAndActivateButton: %v\n", userOptions.SelectAndActivateButton)
	PadPrintf(6, "PrimaryAudioStreamNumberChange: %v\n", userOptions.PrimaryAudioStreamNumberChange)
	PadPrintf(6, "AngleNumberChange: %v\n", userOptions.AngleNumberChange)
	PadPrintf(6, "PopupOn: %v\n", userOptions.PopupOn)
	PadPrintf(6, "PopupOff: %v\n", userOptions.PopupOff)
	PadPrintf(6, "PrimaryPGEnableDisable: %v\n", userOptions.PrimaryPGEnableDisable)
	PadPrintf(6, "PrimaryPGStreamNumberChange: %v\n", userOptions.PrimaryPGStreamNumberChange)
	PadPrintf(6, "SecondaryVideoEnableDisable: %v\n", userOptions.SecondaryVideoEnableDisable)
	PadPrintf(6, "SecondaryVideoStreamNumberChange: %v\n", userOptions.SecondaryVideoStreamNumberChange)
	PadPrintf(6, "SecondaryAudioEnableDisable: %v\n", userOptions.SecondaryAudioEnableDisable)
	PadPrintf(6, "SecondaryAudioStreamNumberChange: %v\n", userOptions.SecondaryAudioStreamNumberChange)
	PadPrintf(6, "SecondaryPGStreamNumberChange: %v\n", userOptions.SecondaryPGStreamNumberChange)
}

func PlayItemEntryPrint(playItemEntry *mpls.PlayItemEntry) {
	PadPrintf(8, "FileName: %s\n", playItemEntry.FileName)
	PadPrintf(8, "Codec: %s\n", playItemEntry.Codec)
	PadPrintf(8, "RefToSTCID: %d\n", playItemEntry.RefToSTCID)
}

func SubPlayItemPrint(subPlayItem *mpls.SubPlayItem) {
	inTime := Convert45KhzTimeToSeconds(subPlayItem.INTime)
	outTime := Convert45KhzTimeToSeconds(subPlayItem.OUTTime)
	PadPrintf(8, "Length: %d\n", subPlayItem.Length)
	PadPrintf(8, "FileName: %s\n", subPlayItem.FileName)
	PadPrintf(8, "Codec: %s\n", subPlayItem.Codec)
	PadPrintf(8, "ConnectionCondition: %d\n", subPlayItem.ConnectionCondition)
	PadPrintf(8, "IsMultiClipEntries: %v\n", subPlayItem.IsMultiClipEntries)
	PadPrintf(8, "RefToSTCID: %d\n", subPlayItem.RefToSTCID)
	PadPrintf(8, "InTime: %d (%d)\n", subPlayItem.INTime, inTime)
	PadPrintf(8, "OUTime: %d (%d)\n", subPlayItem.OUTTime, outTime)
	PadPrintf(10, "*Duration: %v\n", outTime-inTime)
	PadPrintf(8, "SyncPlaytItemID: %d\n", subPlayItem.SyncPlaytItemID)
	PadPrintf(8, "SyncStartPTS: %d\n", subPlayItem.SyncStartPTS)
	PadPrintf(8, "NumberOfMultiClipEntries: %d\n", subPlayItem.NumberOfMultiClipEntries)
	PadPrintf(8, "MultiClipEntries:\n")
	for i := range subPlayItem.MultiClipEntries {
		PadPrintf(6, "MultiClipEntry [%d]:\n", i)
		PlayItemEntryPrint(subPlayItem.MultiClipEntries[i])
	}

}

func StreamTablePrint(streamTable *mpls.StreamTable) {
	PadPrintln(4, "StreamTable:")
	PadPrintf(6, "Length: %d\n", streamTable.Length)

	for _, item := range streamTable.Items {
		PadPrintf(6, "NumberOf%s: %d\n", item.KindOf, item.NumberOf)
		if item.NumberOf != 0 {
			for j, stream := range item.Streams {
				PadPrintf(8, "%s Stream [%d]:\n", item.KindOf, j+1)
				StreamPrint(*stream)
				PadPrintln(8, "---")
			}
		} else {
			PadPrintln(8, "[skip]")
		}
	}
}

//
// Stream printing is an interface
//

// StreamPrint prints the stream information, including its entry and attributes.
func StreamPrint(stream mpls.Stream) {
	if stream.Entry != nil {
		StreamEntryPrint(stream.Entry)
	}

	if stream.Attr != nil {
		StreamAttrPrint(stream.Attr)
	}

}

// StreamEntryPrint prints the stream entry based on its type.
func StreamEntryPrint(streamEntry mpls.StreamEntry) {
	switch streamEntryType := streamEntry.(type) {
	case *mpls.StreamEntryTypeI:
		StreamEntryTypeIPrint(streamEntryType)
	case *mpls.StreamEntryTypeII:
		StreamEntryTypeIIPrint(streamEntryType)
	case *mpls.StreamEntryTypeIII:
		StreamEntryTypeIIIPrint(streamEntryType)
	}
}

// StreamAttrPrint prints the stream attributes based on their specific type.
func StreamAttrPrint(streamAttr mpls.StreamAttributes) {
	switch streamAttrType := streamAttr.(type) {
	case *mpls.PrimaryVideoAttributesH264:
		PrimaryVideoAttributesH264Print(streamAttrType)
	case *mpls.PrimaryVideoAttributesHEVC:
		PrimaryVideoAttributesHEVCPrint(streamAttrType)
	case *mpls.PrimaryAudioAttributes:
		PrimaryAudioAttributesPrint(streamAttrType)
	case *mpls.SecondaryAudioAttributes:
		SecondaryAudioAttributesPrint(streamAttrType)
	case *mpls.SecondaryVideoAttributes:
		SecondaryVideoAttributesPrint(streamAttrType)
	case *mpls.PGAttributes:
		PGAttributesPrint(streamAttrType)
	case *mpls.IGAttributes:
		IGAttributesPrint(streamAttrType)
	case *mpls.TextAttributes:
		TextAttributesPrint(streamAttrType)
	}
}

// StreamEntryTypeIPrint prints the details of a StreamEntryTypeI.
func StreamEntryTypeIPrint(entry *mpls.StreamEntryTypeI) {
	PadPrintln(10, "Entry:")
	PadPrintf(12, "Length: %d\n", entry.Length)
	PadPrintf(12, "StreamType: %d [%s]\n", entry.StreamType, mpls.StreamType(entry.StreamType))
	PadPrintf(12, "RefToStreamPID: %d\n", entry.RefToStreamPID)
}

// StreamEntryTypeIIPrint prints the details of a StreamEntryTypeII.
func StreamEntryTypeIIPrint(entry *mpls.StreamEntryTypeII) {
	PadPrintln(10, "Entry:")
	PadPrintf(12, "Length: %d\n", entry.Length)
	PadPrintf(12, "StreamType: %d [%s]\n", entry.StreamType, mpls.StreamType(entry.StreamType))
	PadPrintf(12, "RefToSubPathID: %d\n", entry.RefToSubPathID)
	PadPrintf(12, "RefToSubClipID: %d\n", entry.RefToSubClipID)
	PadPrintf(12, "RefToStreamPID: %d\n", entry.RefToStreamPID)
}

// StreamEntryTypeIIIPrint prints the details of a StreamEntryTypeIII.
func StreamEntryTypeIIIPrint(entry *mpls.StreamEntryTypeIII) {
	PadPrintln(10, "Entry:")
	PadPrintf(12, "Length: %d\n", entry.Length)
	PadPrintf(12, "StreamType: %d [%s]\n", entry.StreamType, mpls.StreamType(entry.StreamType))
	PadPrintf(12, "RefToSubPathID: %d\n", entry.RefToSubPathID)
	PadPrintf(12, "RefToStreamPID: %d\n", entry.RefToStreamPID)
}

// PrimaryVideoAttributesH264Print prints the attributes of a PrimaryVideoAttributesH264.
func PrimaryVideoAttributesH264Print(attr *mpls.PrimaryVideoAttributesH264) {
	PadPrintln(10, "Attributes:")
	PadPrintf(12, "Length: %d\n", attr.Length)
	PadPrintf(12, "StreamCodingType: %d [%s]\n", attr.StreamCodingType, mpls.StreamCodec(attr.StreamCodingType))
	PadPrintf(12, "Format: %d [%s]\n", attr.Format, mpls.VideoFormat(attr.Format))
	PadPrintf(12, "Rate: %d [%s] \n", attr.Rate, mpls.VideoRate(attr.Rate))
}

// PrimaryVideoAttributesH264Print prints the attributes of a PrimaryVideoAttributesH264.
func PrimaryVideoAttributesHEVCPrint(attr *mpls.PrimaryVideoAttributesHEVC) {
	PadPrintln(10, "Attributes:")
	PadPrintf(12, "Length: %d\n", attr.Length)
	PadPrintf(12, "StreamCodingType: %d [%s]\n", attr.StreamCodingType, mpls.StreamCodec(attr.StreamCodingType))
	PadPrintf(12, "Format: %d [%s]\n", attr.Format, mpls.VideoFormat(attr.Format))
	PadPrintf(12, "Rate: %d [%s]\n", attr.Rate, mpls.VideoRate(attr.Rate))
	PadPrintf(12, "DynamicRangeType: %d\n", attr.DynamicRangeType)
	PadPrintf(12, "ColorSpace: %d\n", attr.ColorSpace)
	PadPrintf(12, "CRFlag: %v\n", attr.CRFlag)
	PadPrintf(12, "HDRPlusFlag: %v\n", attr.HDRPlusFlag)
}

// PrimaryAudioAttributesPrint prints the attributes of a PrimaryAudioAttributes.
func PrimaryAudioAttributesPrint(attr *mpls.PrimaryAudioAttributes) {
	PadPrintln(10, "Attributes:")
	PadPrintf(12, "Length: %d\n", attr.Length)
	PadPrintf(12, "StreamCodingType: %d [%s]\n", attr.StreamCodingType, mpls.StreamCodec(attr.StreamCodingType))
	PadPrintf(12, "Format: %d [%s]\n", attr.Format, mpls.AudioFormat(attr.Format))
	PadPrintf(12, "Rate: %d [%s]\n", attr.Rate, mpls.AudioRate(attr.Rate))
	eng, nat := mpls.LanguageCode(attr.LanguageCode)
	PadPrintf(12, "LanguageCode: %s [%s, %s]\n", attr.LanguageCode, eng, nat)
}

// SecondaryAudioAttributesPrint prints the attributes of a SecondaryAudioAttributes.
// This is used for secondary audio streams in the MPLS structure.
func SecondaryAudioAttributesPrint(attr *mpls.SecondaryAudioAttributes) {
	PadPrintln(10, "Attributes:")
	PadPrintf(12, "Length: %d\n", attr.Length)
	PadPrintf(12, "StreamCodingType: %d [%s]\n", attr.StreamCodingType, mpls.StreamCodec(attr.StreamCodingType))
	PadPrintf(12, "Format: %d [%s]\n", attr.Format, mpls.AudioFormat(attr.Format))
	PadPrintf(12, "Rate: %d [%s]\n", attr.Rate, mpls.AudioRate(attr.Rate))
	eng, nat := mpls.LanguageCode(attr.LanguageCode)
	PadPrintf(12, "LanguageCode: %s [%s, %s]\n", attr.LanguageCode, eng, nat)
	PadPrintf(12, "NumberOfPrimaryAudioRef: %+v\n", attr.NumberOfPrimaryAudioRef)
}

// SecondaryVideoAttributesPrint prints the attributes of a SecondaryVideoAttributes.
// This is used for secondary video streams in the MPLS structure.
func SecondaryVideoAttributesPrint(attr *mpls.SecondaryVideoAttributes) {
	PadPrintln(10, "Attributes:")
	PadPrintf(12, "Length: %d\n", attr.Length)
	PadPrintf(12, "StreamCodingType: %d [%s]\n", attr.StreamCodingType, mpls.StreamCodec(attr.StreamCodingType))
	PadPrintf(12, "Format: %d [%s]\n", attr.Format, mpls.VideoFormat(attr.Format))
	PadPrintf(12, "Rate: %d [%s]\n", attr.Rate, mpls.VideoRate(attr.Rate))
}

// PGAttributesPrint prints the attributes of a PGAttributes.
// This is used for presentation graphics streams in the MPLS structure.
func PGAttributesPrint(attr *mpls.PGAttributes) {
	PadPrintln(10, "Attributes:")
	PadPrintf(12, "Length: %d\n", attr.Length)
	PadPrintf(12, "StreamCodingType: %d [%s]\n", attr.StreamCodingType, mpls.StreamCodec(attr.StreamCodingType))
	eng, nat := mpls.LanguageCode(attr.LanguageCode)
	PadPrintf(12, "LanguageCode: %s [%s, %s]\n", attr.LanguageCode, eng, nat)
}

// IGAttributesPrint prints the attributes of a IGAttributes.
// This is used for interactive graphics streams in the MPLS structure.
func IGAttributesPrint(attr *mpls.IGAttributes) {
	PadPrintln(10, "Attributes:")
	PadPrintf(12, "Length: %d\n", attr.Length)
	PadPrintf(12, "StreamCodingType: %d [%s]\n", attr.StreamCodingType, mpls.StreamCodec(attr.StreamCodingType))
	eng, nat := mpls.LanguageCode(attr.LanguageCode)
	PadPrintf(12, "LanguageCode: %s [%s, %s]\n", attr.LanguageCode, eng, nat)
}

// TextAttributesPrint prints the attributes of a TextAttributes.
// This is used for text subtitle streams in the MPLS structure.
func TextAttributesPrint(attr *mpls.TextAttributes) {
	PadPrintln(10, "Attributes:")
	PadPrintf(12, "Length: %d\n", attr.Length)
	PadPrintf(12, "StreamCodingType: %d [%s]\n", attr.StreamCodingType, mpls.StreamCodec(attr.StreamCodingType))
	PadPrintf(12, "CharacterCode: %d [%s]\n", attr.CharacterCode, mpls.CharacterCode(attr.CharacterCode))
	eng, nat := mpls.LanguageCode(attr.LanguageCode)
	PadPrintf(12, "LanguageCode: %s [%s, %s]\n", attr.LanguageCode, eng, nat)
}

// XXX - There might be other MarkType to print
// PLEASE INVESTIGATE
func PlaylistMarksPrint(playlistMarks *mpls.PlaylistMarks) {
	PadPrintf(0, "Chapter Marks: [%d]\n", len(playlistMarks.Marks))
	for i, mark := range playlistMarks.Marks {
		if mark.MarkType == 1 { // Chapter mark
			timestamp := Parse45KhzTimestamp(mark.MarkTimeStamp)
			PadPrintf(2, "Chapter [%d]: at [%v] (PlayItem: %d)\n", i+1, timestamp, mark.RefToPlayItemID)
			MarkEntryPrint(mark)
			PadPrintln(2, "---")
		}
	}
}

func MarkEntryPrint(markEntry *mpls.MarkEntry) {
	PadPrintf(4, "MarkType: %d\n", markEntry.MarkType)
	PadPrintf(4, "RefToPlayItemID: %d\n", markEntry.RefToPlayItemID)
	PadPrintf(4, "MarkTimeStamp: %d\n", markEntry.MarkTimeStamp)
	PadPrintf(4, "EntryESPID: %d\n", markEntry.EntryESPID)
	PadPrintf(4, "Duration: %d\n", markEntry.Duration)
}

//
// Extensions are abstracted by an interface
//

func ExtensionsPrint(extensions *mpls.Extensions) {
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

func ExtensionsMetaDataPrint(metaData *mpls.ExtensionsMetaData) {
	PadPrintf(4, "ExtensionsMetaData.Length: %d\n", metaData.Length)
	PadPrintf(4, "ExtensionsMetaData.EntryDataStartAddr: %d\n", metaData.EntryDataStartAddr)
	PadPrintf(4, "ExtensionsMetaData.EntryDataCount: %d\n", metaData.EntryDataCount)
}

func ExtensionEntryMetaDataPrint(entryMetaData *mpls.ExtensionEntryMetaData) {
	PadPrintf(4, "ExtensionEntryMetaData.ExtDataType: %d\n", entryMetaData.ExtDataType)
	PadPrintf(4, "ExtensionEntryMetaData.ExtDataVersion: %d\n", entryMetaData.ExtDataVersion)
	PadPrintf(4, "ExtensionEntryMetaData.ExtDataStartAddress: %d\n", entryMetaData.ExtDataStartAddress)
	PadPrintf(4, "ExtensionEntryMetaData.ExtDataLength: %d\n", entryMetaData.ExtDataLength)
}

func ExtensionsEntryDataPrint(entryData mpls.ExtensionEntryData) {
	switch entryType := entryData.(type) {
	case *mpls.ExtensionPIP:
		ExtensionPIPPrint(entryType)
	case *mpls.ExtensionMVCStream:
		ExtensionMVCStreamPrint(entryType)
	case *mpls.ExtensionSubPath:
		ExtensionSubPathPrint(entryType)
	case *mpls.ExtensionStaticMetaData:
		ExtensionStaticMetaDataPrint(entryType)
	}
}

//
// The actual extensions
//

func ExtensionSubPathPrint(subPathExtension *mpls.ExtensionSubPath) {
	PadPrintln(4, "SubPathExtension")
	PadPrintf(6, "Length: %d\n", subPathExtension.Length)
	PadPrintf(6, "Count: %d\n", subPathExtension.Count)
	for i, subPath := range subPathExtension.SubPaths {
		PadPrintf(6, "SubPath [%d]:\n", i+1)
		SubPathPrint(subPath)
	}
}

func ExtensionMVCStreamPrint(extensionMVCStream *mpls.ExtensionMVCStream) {
	PadPrintln(4, "MVC(3D)Extension:")
	for i, mvcStream := range extensionMVCStream.MVCStreams {
		PadPrintf(8, "mvcStream[%d]\n", i+1)
		MVCStreamPrint(mvcStream)
	}
}

func MVCStreamPrint(mvcStream *mpls.MVCStream) {
	PadPrintf(8, "mvcStream.Length: %d\n", mvcStream.Length)
	PadPrintf(8, "mvcStream.FixedOffsetPopUpFlag: %+v\n", mvcStream.FixedOffsetPopUpFlag)
	PadPrintf(8, "mvcStream.Entry:\n")
	if mvcStream.Entry != nil {
		StreamEntryPrint(mvcStream.Entry)
	}

	PadPrintf(8, "mvcStream.Attr:\n")
	if mvcStream.Attr != nil {
		StreamAttrPrint(mvcStream.Attr)
	}
	PadPrintf(8, "mvcStream.NumberOfOffsetSequences: %d\n", mvcStream.NumberOfOffsetSequences)
}

func ExtensionPIPPrint(pip *mpls.ExtensionPIP) {
	PadPrintln(4, "PIP Extension:")
	PadPrintf(6, "Length: %d\n", pip.Length)
	PadPrintf(6, "NumberOfEntries: %d\n", pip.NumberOfEntries)
	for i, pipEntry := range pip.PIPEntries {
		PadPrintf(8, "PIPEntry[%d]:\n", i+1)
		PIPEntryPrint(pipEntry)
		PadPrintln(8, "---")
	}
}

func PIPEntryPrint(pipEntry *mpls.PIPEntry) {
	PadPrintf(10, "ClipRef: %d\n", pipEntry.ClipRef)
	PadPrintf(10, "SecondaryVideoRef: %d\n", pipEntry.SecondaryVideoRef)
	PadPrintf(10, "TimelineType: %d\n", pipEntry.TimelineType)
	PadPrintf(10, "LumaKeyFlag: %v\n", pipEntry.LumaKeyFlag)
	PadPrintf(10, "TrickPlayFlag: %v\n", pipEntry.TrickPlayFlag)
	PadPrintf(10, "UpperLimitLumaKey: %d\n", pipEntry.UpperLimitLumaKey)
	PadPrintf(10, "DataAddress: %d\n", pipEntry.DataAddress)
	if pipEntry.Data != nil {
		PIPDataPrint(pipEntry.Data)
	} else {
		PadPrintln(10, "[Empty PIPData]")
	}
}

func PIPDataPrint(pipData *mpls.PIPData) {
	PadPrintf(12, "NumberOfEntries: %d\n", pipData.NumberOfEntries)
	for i, entry := range pipData.Entries {
		PadPrintf(12, "PIPData[%d]:\n", i+1)
		PIPDataEntryPrint(entry)
		PadPrintln(12, "---")
	}
}

func PIPDataEntryPrint(pipDataEntry *mpls.PIPDataEntry) {
	PadPrintf(14, "Time: %v\n", pipDataEntry.Time)
	PadPrintf(14, "Xpos: %d\n", pipDataEntry.Xpos)
	PadPrintf(14, "Ypos: %d\n", pipDataEntry.Ypos)
	PadPrintf(14, "ScaleFactor: %d [%s]\n", pipDataEntry.ScaleFactor, mpls.PIPScaling(pipDataEntry.ScaleFactor))
}

func ExtensionStaticMetaDataPrint(smExtension *mpls.ExtensionStaticMetaData) {
	PadPrintln(4, "Static MetaData Extension")
	PadPrintf(6, "Length: %d\n", smExtension.Length)
	PadPrintf(6, "Count: %d\n", smExtension.Count)
	for i, entry := range smExtension.Entries {
		PadPrintf(6, "Static MetaData Entry [%d]:\n", i)
		StaticMetaDataEntryPrint(entry)
	}
}

func StaticMetaDataEntryPrint(smEntry *mpls.StaticMetaDataEntry) {
	PadPrintln(4, "Static MetaData Entry")
	PadPrintf(6, "DynamicRangeType: %d\n", smEntry.DynamicRangeType)
	PadPrintf(6, "DisplayPrimariesX: %v\n", smEntry.DisplayPrimariesX)
	PadPrintf(6, "DisplayPrimariesY: %v\n", smEntry.DisplayPrimariesX)
	PadPrintf(6, "WhitePointX: %v\n", smEntry.WhitePointX)
	PadPrintf(6, "WhitePointY: %v\n", smEntry.WhitePointY)
	PadPrintf(6, "MaxDisplayMasteringLuminance: %v\n", smEntry.MaxDisplayMasteringLuminance)
	PadPrintf(6, "MinDisplayMasteringLuminance: %v\n", smEntry.MinDisplayMasteringLuminance)
	PadPrintf(6, "MaxCLL: %v\n", smEntry.MaxCLL)
	PadPrintf(6, "MaxFALL: %v\n", smEntry.MaxFALL)
}
