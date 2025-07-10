package clpi

import "fmt"

// CharacterCodeType is the text stream character (symbol/glyph) encoding standard.
type CharacterCodeType uint8

const (
	TEXT_CHAR_CODE_UTF8          uint8 = 0x01 // Unicode 8-bit
	TEXT_CHAR_CODE_UTF16BE       uint8 = 0x02 // Unicode 16-bit Big Endian
	TEXT_CHAR_CODE_SHIFT_JIS     uint8 = 0x03 // Japanese
	TEXT_CHAR_CODE_EUC_KR        uint8 = 0x04 // Korean
	TEXT_CHAR_CODE_GB18030_20001 uint8 = 0x05 // Chinese National Standard
	TEXT_CHAR_CODE_CN_GB         uint8 = 0x06 // Chinese
	TEXT_CHAR_CODE_BIG5          uint8 = 0x07 // Traditional Chinese
)

func CharacterCode(code uint8) string {
	switch code {
	case TEXT_CHAR_CODE_UTF8:
		return "UTF8"
	case TEXT_CHAR_CODE_UTF16BE:
		return "UTF16BE"
	case TEXT_CHAR_CODE_SHIFT_JIS:
		return "SHIFT JIS"
	case TEXT_CHAR_CODE_EUC_KR:
		return "EUC KR"
	case TEXT_CHAR_CODE_GB18030_20001:
		return "GB18030-2000"
	case TEXT_CHAR_CODE_CN_GB:
		return "GB2312"
	case TEXT_CHAR_CODE_BIG5:
		return "BIG5"
	default:
		return ""
	}
}

type StreamCodingType uint8

/** Stream coding type */
const (
	STREAM_TYPE_VIDEO_MPEG1             StreamCodingType = 0x01 // 1
	STREAM_TYPE_VIDEO_MPEG2             StreamCodingType = 0x02 // 2
	STREAM_TYPE_AUDIO_MPEG1             StreamCodingType = 0x03 // 3
	STREAM_TYPE_AUDIO_MPEG2             StreamCodingType = 0x04 // 4
	STREAM_TYPE_VIDEO_H264              StreamCodingType = 0x1b // 27
	STREAM_TYPE_VIDEO_H264_MVC          StreamCodingType = 0x20 // 32
	STREAM_TYPE_VIDEO_HEVC              StreamCodingType = 0x24 // 36
	STREAM_TYPE_AUDIO_LPCM              StreamCodingType = 0x80 // 128
	STREAM_TYPE_AUDIO_AC3               StreamCodingType = 0x81 // 129
	STREAM_TYPE_AUDIO_DTS               StreamCodingType = 0x82 // 130
	STREAM_TYPE_AUDIO_TRUHD             StreamCodingType = 0x83 // 131
	STREAM_TYPE_AUDIO_AC3PLUS           StreamCodingType = 0x84 // 132
	STREAM_TYPE_AUDIO_DTSHD             StreamCodingType = 0x85 // 133
	STREAM_TYPE_AUDIO_DTSHD_MASTER      StreamCodingType = 0x86 // 134
	STREAM_TYPE_SUB_PG                  StreamCodingType = 0x90 // 144
	STREAM_TYPE_SUB_IG                  StreamCodingType = 0x91 // 145
	STREAM_TYPE_SUB_TEXT                StreamCodingType = 0x92 // 146
	STREAM_TYPE_AUDIO_AC3PLUS_SECONDARY StreamCodingType = 0xa1 // 161
	STREAM_TYPE_AUDIO_DTSHD_SECONDARY   StreamCodingType = 0xa2 // 162
	STREAM_TYPE_VIDEO_VC1               StreamCodingType = 0xea // 234
)

func StreamCodec(code StreamCodingType) string {
	switch code {
	case STREAM_TYPE_VIDEO_MPEG1:
		return "MPEG1 VIDEO"
	case STREAM_TYPE_VIDEO_MPEG2:
		return "MPEG2 VIDEO"
	case STREAM_TYPE_AUDIO_MPEG1:
		return "MPEG1 AUDIO"
	case STREAM_TYPE_AUDIO_MPEG2:
		return "MPEG2 AUDIO"
	case STREAM_TYPE_VIDEO_H264:
		return "H264 VIDEO"
	case STREAM_TYPE_VIDEO_H264_MVC:
		return "H264 MULTI VIDEO CODING (STEREOSCOPIC 3D) VIDEO"
	case STREAM_TYPE_VIDEO_HEVC:
		return "HEVC VIDEO"
	case STREAM_TYPE_AUDIO_LPCM:
		return "LPCM AUDIO"
	case STREAM_TYPE_AUDIO_AC3:
		return "AC3 AUDIO"
	case STREAM_TYPE_AUDIO_DTS:
		return "DTS AUDIO"
	case STREAM_TYPE_AUDIO_TRUHD:
		return "TRUEHD AUDIO"
	case STREAM_TYPE_AUDIO_AC3PLUS:
		return "AC3PLUS AUDIO"
	case STREAM_TYPE_AUDIO_DTSHD:
		return "DTSHD AUDIO"
	case STREAM_TYPE_AUDIO_DTSHD_MASTER:
		return "DTSHD MASTER AUDIO"
	case STREAM_TYPE_SUB_PG:
		return "PRESENTATION GRAPHICS SUBTITLE"
	case STREAM_TYPE_SUB_IG:
		return "INTERACTIVE GRAPHICS SUBTITLE"
	case STREAM_TYPE_SUB_TEXT:
		return "TEXT SUBTITLE"
	case STREAM_TYPE_AUDIO_AC3PLUS_SECONDARY:
		return "AC3PLUS SECONDARY AUDIO"
	case STREAM_TYPE_AUDIO_DTSHD_SECONDARY:
		return "DTSHD SECONDARY AUDIO"
	case STREAM_TYPE_VIDEO_VC1:
		return "VC1 VIDEO"
	default:
		return ""
	}
}

// VideoFormatType defines the video format.
type VideoFormatType uint8

const (
	VIDEO_FORMAT_480I  VideoFormatType = 1 // ITU-R BT.601-5
	VIDEO_FORMAT_576I  VideoFormatType = 2 // ITU-R BT.601-4
	VIDEO_FORMAT_480P  VideoFormatType = 3 // SMPTE 293M
	VIDEO_FORMAT_1080I VideoFormatType = 4 // SMPTE 274M
	VIDEO_FORMAT_720P  VideoFormatType = 5 // SMPTE 296M
	VIDEO_FORMAT_1080P VideoFormatType = 6 // SMPTE 274M
	VIDEO_FORMAT_576P  VideoFormatType = 7 // ITU-R BT.1358
	VIDEO_FORMAT_2160P VideoFormatType = 8 // BT.2020
)

func VideoFormat(code VideoFormatType) string {
	switch code {
	case VIDEO_FORMAT_480I:
		return "480I"
	case VIDEO_FORMAT_576I:
		return "576I"
	case VIDEO_FORMAT_480P:
		return "480P"
	case VIDEO_FORMAT_1080I:
		return "1080I"
	case VIDEO_FORMAT_720P:
		return "720P"
	case VIDEO_FORMAT_1080P:
		return "1080P"
	case VIDEO_FORMAT_576P:
		return "576P"
	case VIDEO_FORMAT_2160P:
		return "2160P"
	default:
		return ""
	}
}

// VideoRateType defines the video refresh rate.
type VideoRateType uint8

const (
	VIDEO_RATE_24000_1001 VideoRateType = 1 // 23.976 Hz
	VIDEO_RATE_24000_1000 VideoRateType = 2 // 24 Hz
	VIDEO_RATE_25000_1000 VideoRateType = 3 // 25 Hz
	VIDEO_RATE_30000_1001 VideoRateType = 4 // 29.97 Hz
	VIDEO_RATE_50000_1000 VideoRateType = 6 // 50 Hz
	VIDEO_RATE_60000_1001 VideoRateType = 7 // 59.94 Hz
)

func VideoRate(code VideoRateType) string {
	switch code {
	case VIDEO_RATE_24000_1001:
		return fmt.Sprintf("%.3f Hz", float32(24000)/1001)
	case VIDEO_RATE_24000_1000:
		return fmt.Sprintf("%d Hz", 24)
	case VIDEO_RATE_25000_1000:
		return fmt.Sprintf("%d Hz", 25)
	case VIDEO_RATE_30000_1001:
		return fmt.Sprintf("%.3f Hz", float32(30000)/1001)
	case VIDEO_RATE_50000_1000:
		return fmt.Sprintf("%d Hz", 50)
	case VIDEO_RATE_60000_1001:
		return fmt.Sprintf("%.3f Hz", float32(60000)/1001)
	default:
		return ""
	}
}

// AudioFormatType defines the audio codec format.
type AudioFormatType uint8

const (
	AUDIO_FORMAT_MONO         AudioFormatType = 0x01 // Mono audio
	AUDIO_FORMAT_STEREO       AudioFormatType = 0x02 // Stereo audio
	AUDIO_FORMAT_MULTICHANNEL AudioFormatType = 0x03 // Multi-channel audio
	AUDIO_FORMAT_COMBO        AudioFormatType = 0x06 // Stereo: ac3/dts; Multi-Channel: mlp/dts-hd
)

func AudioFormat(code AudioFormatType) string {
	switch code {
	case AUDIO_FORMAT_MONO:
		return "Mono"
	case AUDIO_FORMAT_STEREO:
		return "Stereo"
	case AUDIO_FORMAT_MULTICHANNEL:
		return "Multi-Channel"
	case AUDIO_FORMAT_COMBO:
		return "Stereo OR Multi-channel"
	default:
		return ""
	}
}

// AudioRateType defines the audio sampling rate.
type AudioRateType uint8

const (
	AUDIO_RATE_48kHZ  AudioRateType = 0x01 // 48 kHz
	AUDIO_RATE_96kHZ  AudioRateType = 0x02 // 96 kHz
	AUDIO_RATE_192kHZ AudioRateType = 0x03 // 192 kHz
)

func AudioRate(code AudioRateType) string {
	switch code {
	case AUDIO_RATE_48kHZ:
		return "48 kHz"
	case AUDIO_RATE_96kHZ:
		return "96 kHz"
	case AUDIO_RATE_192kHZ:
		return "192 kHz"
	default:
		return ""
	}
}

type VideoAspectRatioType uint8

const (
	VIDEO_ASPECT_RATIO_4_3  VideoAspectRatioType = 2 //  4:3 legacy
	VIDEO_ASPECT_RATIO_16_9 VideoAspectRatioType = 3 // 16:0 modern
)

func AspectRatio(code VideoAspectRatioType) string {
	switch code {
	case VIDEO_ASPECT_RATIO_4_3:
		return "4:3"
	case VIDEO_ASPECT_RATIO_16_9:
		return "16:9"
	default:
		return ""
	}
}

type ClipApplicationType uint8

const (
	CLIP_APP_TYPE_1 ClipApplicationType = 1 // "Main TS for a main-path of Movie"
	CLIP_APP_TYPE_2 ClipApplicationType = 2 // "Main TS for a main-path of Time based slide show"
	CLIP_APP_TYPE_3 ClipApplicationType = 3 // "Main TS for a main-path of Browsable slide show"
	CLIP_APP_TYPE_4 ClipApplicationType = 4 // "Sub TS for a sub-path of Browsable slide show"
	CLIP_APP_TYPE_5 ClipApplicationType = 5 // "Sub TS for a sub-path of Interactive Graphics menu"
	CLIP_APP_TYPE_6 ClipApplicationType = 6 // "Sub TS for a sub-path of Text subtitle"
	CLIP_APP_TYPE_7 ClipApplicationType = 7 // "Sub TS for a sub-path of one or more elementary streams path"
	CLIP_APP_TYPE_8 ClipApplicationType = 8 // "Sub TS for a main-path of Enhanced LR View"
)

func ClipApplication(code ClipApplicationType) string {
	switch code {
	case CLIP_APP_TYPE_1:
		return "Main TS for a main-path of Movie"
	case CLIP_APP_TYPE_2:
		return "Main TS for a main-path of Time based slide show"
	case CLIP_APP_TYPE_3:
		return "Main TS for a main-path of Browsable slide show"
	case CLIP_APP_TYPE_4:
		return "Sub TS for a sub-path of Browsable slide show"
	case CLIP_APP_TYPE_5:
		return "Sub TS for a sub-path of Interactive Graphics menu"
	case CLIP_APP_TYPE_6:
		return "Sub TS for a sub-path of Text subtitle"
	case CLIP_APP_TYPE_7:
		return "Sub TS for a sub-path of one or more elementary streams path"
	case CLIP_APP_TYPE_8:
		return "Sub TS for a main-path of Enhanced LR View"
	default:
		return ""
	}
}
