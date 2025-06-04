package main

//
// Just constant data tables of various information
// Some of this was discovered in libbluray source code.
//

// NOTE: Go doesn't allow for good constants

const (
	TEXT_CHAR_CODE_UTF8          = 0x01 // Unicode 8-bit
	TEXT_CHAR_CODE_UTF16BE       = 0x02 // Unicode 16-bit Big Endian
	TEXT_CHAR_CODE_SHIFT_JIS     = 0x03 // Japanese
	TEXT_CHAR_CODE_EUC_KR        = 0x04 // Korean
	TEXT_CHAR_CODE_GB18030_20001 = 0x05 // Chinese National Standard
	TEXT_CHAR_CODE_CN_GB         = 0x06 // Chinese
	TEXT_CHAR_CODE_BIG5          = 0x07 // Traditional Chinese
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

/** Stream video coding type */
const (
	STREAM_TYPE_VIDEO_MPEG1             = 0x01 // 1
	STREAM_TYPE_VIDEO_MPEG2             = 0x02 // 2
	STREAM_TYPE_AUDIO_MPEG1             = 0x03 // 3
	STREAM_TYPE_AUDIO_MPEG2             = 0x04 // 4
	STREAM_TYPE_VIDEO_H264              = 0x1b // 27
	STREAM_TYPE_VIDEO_H264_MVC          = 0x20 // 32
	STREAM_TYPE_VIDEO_HEVC              = 0x24 // 36
	STREAM_TYPE_AUDIO_LPCM              = 0x80 // 128
	STREAM_TYPE_AUDIO_AC3               = 0x81 // 129
	STREAM_TYPE_AUDIO_DTS               = 0x82 // 130
	STREAM_TYPE_AUDIO_TRUHD             = 0x83 // 131
	STREAM_TYPE_AUDIO_AC3PLUS           = 0x84 // 132
	STREAM_TYPE_AUDIO_DTSHD             = 0x85 // 133
	STREAM_TYPE_AUDIO_DTSHD_MASTER      = 0x86 // 134
	STREAM_TYPE_SUB_PG                  = 0x90 // 144
	STREAM_TYPE_SUB_IG                  = 0x91 // 145
	STREAM_TYPE_SUB_TEXT                = 0x92 // 146
	STREAM_TYPE_AUDIO_AC3PLUS_SECONDARY = 0xa1 // 161
	STREAM_TYPE_AUDIO_DTSHD_SECONDARY   = 0xa2 // 162
	STREAM_TYPE_VIDEO_VC1               = 0xea // 234
)

func StreamCodec(code uint8) string {
	switch code {
	case STREAM_TYPE_VIDEO_MPEG1:
		return "MPEG1 VIDEO"
	case STREAM_TYPE_VIDEO_MPEG2:
		return "MPEG2 VIDEO"
	case STREAM_TYPE_AUDIO_MPEG1:
		return "MPEG1 AUDIO"
	case STREAM_TYPE_AUDIO_MPEG2:
		return "MPEG2 AUDIO"
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
	case STREAM_TYPE_VIDEO_VC1:
		return "VC1 VIDEO"
	case STREAM_TYPE_VIDEO_H264:
		return "H264 VIDEO"
	case STREAM_TYPE_VIDEO_H264_MVC:
		return "H264 MULTI VIDEO CODING (STEREOSCOPIC 3D) VIDEO"
	case STREAM_TYPE_VIDEO_HEVC:
		return "HEVC VIDEO"
	case STREAM_TYPE_SUB_PG:
		return "PRESENTATION GRAPHICS SUBTITLE"
	case STREAM_TYPE_SUB_IG:
		return "INTERACTIVE GRAPHICS SUBTITLE"
	case STREAM_TYPE_SUB_TEXT:
		return "TEXT SUBTITLE"
	case 161:
		return "AC3PLUS SECONDARY AUDIO"
	case 162:
		return "DTSHD SECONDARY AUDIO"
	default:
		return ""
	}
}

const (
	sub_path_pabs      = 0x02 /* Primary audio of the Browsable slideshow */
	sub_path_ig_menu   = 0x03 /* Interactive Graphics presentation menu */
	sub_path_textst    = 0x04 /* Text Subtitle */
	sub_path_async_es  = 0x05 /* Out-of-mux Synchronous elementary streams */
	sub_path_async_pip = 0x06 /* Out-of-mux Asynchronous Picture-in-Picture presentation */
	sub_path_sync_pip  = 0x07 /* In-mux Synchronous Picture-in-Picture presentation */
	sub_path_ss_video  = 0x08 /* SS Video */
	sub_path_dv_el     = 0x0a /* Dolby Vision Enhancement Layer */
)

func SubPathType(code uint8) string {
	switch code {
	case sub_path_pabs:
		return "PRIMARY AUDIO"
	case sub_path_ig_menu:
		return ""
	case sub_path_textst:
		return ""
	case sub_path_async_es:
		return ""
	case sub_path_async_pip:
		return ""
	case sub_path_sync_pip:
		return ""
	case sub_path_ss_video:
		return ""
	case sub_path_dv_el:
		return ""
	default:
		return ""
	}
}
