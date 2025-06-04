

### Steam Coding types

| Hex  | Dec | Label                 | Type     | Remark                   |
| ---- | --: | --------------------- | -------- | -                        |
| 0x01 |   1 | MPEG1                 | Video    |                          |
| 0x02 |   2 | MPEG2                 | Video    |                          |
| 0x03 |   3 | MPEG1                 | Audio    |                          |
| 0x04 |   4 | MPEG2                 | Audio    |                          |
| 0x1b |  27 | H264                  | Video    |                          |
| 0x20 |  32 | MVC                   | Video    | **M**ulti **V**iew **C**oding is an extension of H264 to render 3D video   |
| 0x24 |  36 | HEVC                  | Video    | aka H265                 |
| 0x80 | 128 | LPCM                  | Audio    |                          |
| 0x81 | 129 | AC3                   | Audio    |                          |
| 0x82 | 130 | DTS                   | Audio    |                          |
| 0x83 | 131 | TRU HD                | Audio    |                          |
| 0x84 | 132 | AC3 PLUS              | Audio    |                          |
| 0x85 | 133 | DTS HD                | Audio    |                          |
| 0x86 | 134 | DTS HD MASTER         | Audio    |                          |
| 0x90 | 144 | PRESENTATION GRAPHICS | Subtitle |                          |
| 0x91 | 145 | INTERACTIVE GRAPHICS  | Subtitle |                          |
| 0x92 | 146 | TEXT                  | Subtitle | Found in PG streams      |
| 0xa1 | 161 | AC3 PLUS (SECONDARY)  | Audio    |                          |
| 0xa2 | 162 | DTS HD (SECONDARY)    | Audio    |                          |
| 0xea | 234 | VC1                   | Video    |                          |

---

### Language Code types

| Hex  | Label         | Remark                    |
| -    | -             | -                         |
| 0x01 | UTF8          | Unicode 8-bit             |
| 0x02 | UTF16BE       | Unicode 16-bit Big Endian | 
| 0x03 | SHIFT_JIS     | Japanese                  |
| 0x04 | EUD_KR        | Korean                    |
| 0x05 | GB18030_20001 | Chinese National Standard | 
| 0x06 | CN_GB         | Chinese                   |
| 0x07 | BIG5          | Chinese Traditional       |

---

### Still Mode Code types

| Hex  | Label    | Remark                     |
| -    | -        | -                          |
| 0x00 | NONE     | No Still (normal playback) |
| 0x01 | TIME     | Finite Still Time          |
| 0x02 | INFINITE | Infinite Still Time        |

---

### Stream Audio Sample Rate

| Hex  | Bin        | Dec | Label     | Remark                                                        |
| -    |  -         |  -: | -         | -                                                             |
| 0x01 | ``0b0001`` |   1 | 48        | -                                                             |
| 0x04 | ``0b0100`` |   4 | 96        | -                                                             |
| 0x05 | ``0b0101`` |   5 | 192       | -                                                             |
| 0x0c | ``0b1100`` |  12 | 192_COMBO | 4x 48 (ac3/dts) <br> 2x 96 (ac3/dts) <br> 1x 192 (mpl/dts-hd) |
| 0x0e | ``0b1110`` |  14 | 96_COMBO  | 2x 48 (ac3/dts) <br> 1x 96 (mpl/dts-hd)                       |

---

### Stream Video Frame Rate

| Hex      | Bin            | Dec   | Label      | Remark    |
| -        |  -             |  -:   | -          | -         |
|   0x01   | ``0b0001``     |   1   | 24000_1001 | 23.976 Hz |
|   0x02   | ``0b0010``     |   2   | 24         | 24 Hz     |
|   0x03   | ``0b0011``     |   3   | 48         | 25 Hz     |
|   0x04   | ``0b0100``     |   4   | 30000_1001 | 29.97 Hz  |
| ~~0x05~~ | ~~``0b0101``~~ | ~~5~~ | -          | n/a       |
|   0x06   | ``0b0110``     |   6   | 50         | 50 Hz     |
|   0x07   | ``0b0111``     |   7   | 60000_1001 | 59.94 Hz  |

Notes:
* The underscore indicate division, for example $\frac{24000}{1001}=23.976 \,\text{Hz}$
* Number 5 is missing, and no idea why.

---