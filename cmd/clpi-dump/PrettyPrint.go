package main

import (
	"fmt"

	clpi "github.com/parasense/bdmv_go/pkg/clpi"
)

func HeaderPrint(header *clpi.CLPIHeader) {
	fmt.Println("Header:")
	PadPrintf(2, "Type: %s\n", string(header.TypeIndicator[:]))
	PadPrintf(2, "Version: %s\n", string(header.VersionNumber[:]))
	PadPrintf(2, "Offset: ClipInfo: [%d:%d]\n", header.ClipInfo.Start, header.ClipInfo.Stop)
	PadPrintf(2, "Offset: SequenceInfo: [%d:%d]\n", header.SequenceInfo.Start, header.SequenceInfo.Stop)
	PadPrintf(2, "Offset: ProgramInfo: [%d:%d]\n", header.ProgramInfo.Start, header.ProgramInfo.Stop)
	PadPrintf(2, "Offset: CPI: [%d:%d]\n", header.CPI.Start, header.CPI.Stop)
	PadPrintf(2, "Offset: ClipMarks: [%d:%d]\n", header.ClipMarks.Start, header.ClipMarks.Stop)
	PadPrintf(2, "Offset: Extensions: [%d:%d]\n", header.Extensions.Start, header.Extensions.Stop)
	PadPrintln(2, "---")
}

func ClipInfoPrint(clipInfo *clpi.ClipInfo) {
	fmt.Println("ClipInfo:")
	PadPrintf(2, "Length: %d\n", clipInfo.Length)
	PadPrintf(2, "ClipStreamType: %d\n", clipInfo.ClipStreamType)
	PadPrintf(2, "ApplicationType: %d [%s]\n", clipInfo.ApplicationType, clpi.ClipApplication(clipInfo.ApplicationType))
	PadPrintf(2, "IsCC5: %t\n", clipInfo.IsCC5)
	PadPrintf(2, "TSRecordingRate: %d\n", clipInfo.TSRecordingRate)
	PadPrintf(2, "NumberOfSourcePackets: %d\n", clipInfo.NumberOfSourcePackets)
	PadPrintf(2, "TSTypeInfoBlock: %+v\n", clipInfo.TSTypeInfoBlock)
	if clipInfo.IsCC5 {
		PadPrintf(4, "FollowingClipStreamType: %d\n", clipInfo.FollowingClipStreamType)
		PadPrintf(4, "FollowingClipInformationFileName: %+v\n", clipInfo.FollowingClipInformationFileName)
		PadPrintf(4, "FollowingClipCodecIdentifier: %+v\n", clipInfo.FollowingClipCodecIdentifier)
	}
	PadPrintln(2, "---")
}

func SequenceInfoPrint(sequenceInfo *clpi.SequenceInfo) {
	fmt.Println("SequenceInfo:")
	PadPrintf(2, "Length: %d\n", sequenceInfo.Length)
	PadPrintf(2, "NumberOfATCSequences: %d\n", sequenceInfo.NumberOfATCSequences)
	PadPrintln(2, "ATCSequences:")
	for i, atc := range sequenceInfo.ATCSequences {
		PadPrintf(4, "ATCSequence[%d]:\n", i+1)
		PadPrintf(6, "SPNATCStart: %d\n", atc.SPNATCStart)
		PadPrintf(6, "NumberOfSTCSequences: %d\n", atc.NumberOfSTCSequences)
		PadPrintf(6, "OffsetSTCID: %d\n", atc.OffsetSTCID)
		PadPrintln(6, "STCSequences:")
		for j, stc := range atc.STCSequences {
			PadPrintf(8, "STCSequence[%d]:\n", j+1)
			PadPrintf(10, "PCRPID: %d\n", stc.PCRPID)
			PadPrintf(10, "SPNSTCStart: %d\n", stc.SPNSTCStart)
			PadPrintf(10, "PresentationStartTime: %d\n", stc.PresentationStartTime)
			PadPrintf(10, "PresentationEndTime: %d\n", stc.PresentationEndTime)
			PadPrintln(8, "---")
		}
		PadPrintln(4, "---")
	}
	PadPrintln(2, "---")
}

func ProgramInfoPrint(pi *clpi.ProgramInfo) {
	fmt.Println("ProgramInfo:")
	PadPrintf(2, "Length: %d\n", pi.Length)
	PadPrintf(2, "NumberOfPrograms: %d\n", pi.NumberOfPrograms)
	if pi.NumberOfPrograms > 0 {
		PadPrintln(2, "Programs:")
		for i, pgm := range pi.Programs {
			PadPrintf(4, "Program [%d]:\n", i+1)
			PadPrintf(6, "SPNProgramSequenceStart: %d\n", pgm.SPNProgramSequenceStart)
			PadPrintf(6, "ProgramMapPID: %d\n", pgm.ProgramMapPID)
			PadPrintf(6, "NumberOfStreamsInPS: %d\n", pgm.NumberOfStreamsInPS)
			if pgm.NumberOfStreamsInPS > 0 {

				PadPrintln(6, "ProgramStreams:")
				for j, pgmstrm := range pgm.ProgramStreams {
					PadPrintf(6, "Program Stream [%d]:\n", j+1)
					PadPrintf(8, "StreamPID: %d\n", pgmstrm.StreamPID)

					for k, sci := range pgmstrm.StreamCodingInfo {
						switch streamType := sci.(type) {
						case *clpi.StreamCodingInfoH264:
							PadPrintf(10, "Length: %d\n", streamType.Length)
							PadPrintf(10, "StreamCodingType: %d [%s]\n", streamType.StreamCodingType, clpi.StreamCodec(streamType.StreamCodingType))
							PadPrintf(10, "ISRCode: %s\n", streamType.ISRCode[:])
							PadPrintf(10, "Format: %d [%s]\n", streamType.VideoFormat, clpi.VideoFormat(streamType.VideoFormat))
							PadPrintf(10, "Rate: %d [%s] \n", streamType.FrameRate, clpi.VideoRate(streamType.FrameRate))
							PadPrintf(10, "AspectRatio: %d [%s]\n", streamType.VideoAspectRatio, clpi.AspectRatio(streamType.VideoAspectRatio))
							PadPrintf(10, "OCFlag: %t\n", streamType.OCFlag)

						case *clpi.StreamCodingInfoH265:
							PadPrintf(10, "Length: %d\n", streamType.Length)
							PadPrintf(10, "StreamCodingType: %d [%s]\n", streamType.StreamCodingType, clpi.StreamCodec(streamType.StreamCodingType))
							PadPrintf(10, "ISRCode: %s\n", streamType.ISRCode[:])
							PadPrintf(10, "Format: %d [%s]\n", streamType.VideoFormat, clpi.VideoFormat(streamType.VideoFormat))
							PadPrintf(10, "Rate: %d [%s] \n", streamType.FrameRate, clpi.VideoRate(streamType.FrameRate))
							PadPrintf(10, "AspectRatio: %d [%s]\n", streamType.VideoAspectRatio, clpi.AspectRatio(streamType.VideoAspectRatio))
							PadPrintf(10, "OCFlag: %t\n", streamType.OCFlag)
							PadPrintf(10, "CRFlag: %t\n", streamType.CRFlag)
							PadPrintf(10, "DynamicRangeType: %d\n", streamType.DynamicRangeType)
							PadPrintf(10, "ColorSpace: %d\n", streamType.ColorSpace)
							PadPrintf(10, "HDRPlusFlag: %t\n", streamType.HDRPlusFlag)

						case *clpi.StreamCodingInfoAudio:
							PadPrintf(10, "Length: %d\n", streamType.Length)
							PadPrintf(10, "StreamCodingType: %d [%s]\n", streamType.StreamCodingType, clpi.StreamCodec(streamType.StreamCodingType))
							PadPrintf(10, "ISRCode: %s\n", streamType.ISRCode[:])
							PadPrintf(10, "Format: %d [%s]\n", streamType.AudioFormat, clpi.AudioFormat(streamType.AudioFormat))
							PadPrintf(10, "Rate: %d [%s]\n", streamType.SampleRate, clpi.AudioRate(streamType.SampleRate))
							PadPrintf(10, "LanguageCode: %s\n", streamType.LanguageCode[:])

						case *clpi.StreamCodingTypePG:
							PadPrintf(10, "Length: %d\n", streamType.Length)
							PadPrintf(10, "StreamCodingType: %d [%s]\n", streamType.StreamCodingType, clpi.StreamCodec(streamType.StreamCodingType))
							PadPrintf(10, "ISRCode: %s\n", streamType.ISRCode[:])
							PadPrintf(10, "LanguageCode: %s\n", streamType.LanguageCode[:])

						case *clpi.StreamCodingTypeIG:
							PadPrintf(10, "Length: %d\n", streamType.Length)
							PadPrintf(10, "StreamCodingType: %d [%s]\n", streamType.StreamCodingType, clpi.StreamCodec(streamType.StreamCodingType))
							PadPrintf(10, "ISRCode: %s\n", streamType.ISRCode[:])
							PadPrintf(10, "LanguageCode: %s\n", streamType.LanguageCode[:])

						case *clpi.StreamCodingTypeText:
							PadPrintf(10, "Length: %d\n", streamType.Length)
							PadPrintf(10, "StreamCodingType: %d [%s]\n", streamType.StreamCodingType, clpi.StreamCodec(streamType.StreamCodingType))
							PadPrintf(10, "ISRCode: %s\n", streamType.ISRCode[:])
							PadPrintf(10, "CharacterCode: %d\n", streamType.CharacterCode)
							PadPrintf(10, "LanguageCode: %s\n", streamType.LanguageCode[:])

						default:
							PadPrintf(12, "[%d] StreamType: %+v\n", k+1, streamType)
						}
					}
					PadPrintln(6, "---")
				}
			}
			PadPrintln(4, "---")
		}
	}
	PadPrintln(2, "---")
}

func CPIPrint(cpi *clpi.CPI) {
	fmt.Println("CPI:")
	PadPrintf(2, "Length: %d\n", cpi.Length)
	PadPrintf(2, "CPIType: %d\n", cpi.CPIType)
	PadPrintf(2, "NumberOfStreamPIDEntries: %d\n", cpi.NumberOfStreamPIDEntries)
	if cpi.NumberOfStreamPIDEntries > 0 {
		PadPrintln(2, "StreamPIDEntries:")
		for i, streamPidEntry := range cpi.StreamPIDEntries {
			PadPrintf(4, "StreamPIDEntry [%d]:\n", i+1)
			PadPrintf(6, "StreamPID: %d\n", streamPidEntry.StreamPID)
			PadPrintf(6, "EPStreamType: %d\n", streamPidEntry.EPStreamType)
			PadPrintf(6, "NumberOfEPCoarseEntries: %d\n", streamPidEntry.NumberOfEPCoarseEntries)
			PadPrintf(6, "NumberOfEPFineEntries: %d\n", streamPidEntry.NumberOfEPFineEntries)
			PadPrintf(6, "EPMapStreamStartAddr: %d\n", streamPidEntry.EPMapStreamStartAddr)
			PadPrintf(6, "EPFineTableStartAddress: %d\n", streamPidEntry.EPFineTableStartAddress)
			if streamPidEntry.NumberOfEPCoarseEntries > 0 {
				PadPrintln(6, "CourseEntries:")
				for j, courseEntry := range streamPidEntry.CourseEntries {
					PadPrintf(8, "Course Entry [%d]:\n", j+1)
					PadPrintf(10, "RefToEPFineID: %d\n", courseEntry.RefToEPFineID)
					PadPrintf(10, "PTSEPCoarse: %d\n", courseEntry.PTSEPCoarse)
					PadPrintf(10, "SPNEPCoarse: %d\n", courseEntry.SPNEPCoarse)
					PadPrintln(8, "---")
				}
			}
			if streamPidEntry.NumberOfEPFineEntries > 0 {
				PadPrintln(6, "FineEntries:")
				for j, fineEntry := range streamPidEntry.FineEntries {
					PadPrintf(8, "Fine Entry [%d]:\n", j+1)
					PadPrintf(10, "IsAngleChangePoint: %t\n", fineEntry.IsAngleChangePoint)
					PadPrintf(10, "IEndPositionOffset: %d\n", fineEntry.IEndPositionOffset)
					PadPrintf(10, "PTSEPFine: %d\n", fineEntry.PTSEPFine)
					PadPrintf(10, "SPNEPFine: %d\n", fineEntry.SPNEPFine)
					PadPrintln(8, "---")
				}
				PadPrintln(6, "---")
			}
			PadPrintln(4, "---")
		}
	}
	PadPrintln(2, "---")
}

func ClipMarksPrint(clipMarks *clpi.ClipMarks) {
	fmt.Println("ClipMarks:")
	PadPrintf(2, "Length: %d\n", clipMarks.Length)
	PadPrintf(2, "NumberOfClipMarks: %d\n", clipMarks.NumberOfClipMarks)
	PadPrintln(2, "MarkEntries:")
	for i, entry := range clipMarks.MarkEntries {
		PadPrintf(4, "LenMarkEntry[%d]:\n", i)
		PadPrintf(6, "MarkType: %d\n", entry.MarkType)
		PadPrintf(6, "MarkPID: %d\n", entry.MarkPID)
		PadPrintf(6, "MarkTimeStamp: %d\n", entry.MarkTimeStamp)
		PadPrintf(6, "MarkEntryPoint: %d\n", entry.MarkEntryPoint)
		PadPrintf(6, "MarkDuration: %d\n", entry.MarkDuration)
		PadPrintln(4, "---")
	}
	PadPrintln(2, "---")
}

func ExtensionsPrint(extensions *clpi.Extensions) {
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

func ExtensionsMetaDataPrint(metaData *clpi.ExtensionsMetaData) {
	PadPrintf(4, "ExtensionsMetaData.Length: %d\n", metaData.Length)
	PadPrintf(4, "ExtensionsMetaData.EntryDataStartAddr: %d\n", metaData.EntryDataStartAddr)
	PadPrintf(4, "ExtensionsMetaData.EntryDataCount: %d\n", metaData.EntryDataCount)
}

func ExtensionEntryMetaDataPrint(entryMetaData *clpi.ExtensionEntryMetaData) {
	PadPrintf(4, "ExtensionEntryMetaData.ExtDataType: %d\n", entryMetaData.ExtDataType)
	PadPrintf(4, "ExtensionEntryMetaData.ExtDataVersion: %d\n", entryMetaData.ExtDataVersion)
	PadPrintf(4, "ExtensionEntryMetaData.ExtDataStartAddress: %d\n", entryMetaData.ExtDataStartAddress)
	PadPrintf(4, "ExtensionEntryMetaData.ExtDataLength: %d\n", entryMetaData.ExtDataLength)
}

func ExtensionsEntryDataPrint(entryData clpi.ExtensionEntryData) {
	switch entryType := entryData.(type) {

	case *clpi.ExtensionLPCMDownMixCoefficient:
		//ExtensionLPCMDownMixCoefficientPrint(entryType)

	case *clpi.ExtensionExtentStartPoints:
		ExtensionExtentStartPointsPrint(entryType)

	case *clpi.ExtensionProgramInfoSS:
		ExtensionProgramInfoSSPrint(entryType)

	case *clpi.ExtensionCPISS:
		ExtensionCPISSPrint(entryType)
	}
}

func ExtensionExtentStartPointsPrint(ext *clpi.ExtensionExtentStartPoints) {
	PadPrintln(4, "ExtensionExtentStartPoints:")
	PadPrintf(6, "Length: %d\n", ext.Length)
	PadPrintf(6, "NumberOfPoints: %d\n", ext.NumberOfPoints)
	if ext.NumberOfPoints > 0 {
		PadPrintln(6, "Point Entries:")
		for i, pnt := range ext.PointEntries {
			PadPrintf(8, "Point[%d]: %d\n", i, pnt.Point)
		}
		PadPrintln(6, "---")
	}
	PadPrintln(4, "---")
}

func ExtensionProgramInfoSSPrint(ext *clpi.ExtensionProgramInfoSS) {
	PadPrintln(4, "ExtensionProgramInfoSS:")
	ProgramInfoPrint(ext)
}

func ExtensionCPISSPrint(ext *clpi.ExtensionCPISS) {
	PadPrintln(4, "ExtensionCPISS:")
	CPIPrint(ext)
}
