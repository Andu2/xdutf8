package streamvalidator

import "testing"

func TestUtf8StreamValidator(t *testing.T) {
	type testChunk struct {
		octets             []byte
		expectedInvalid    bool
		expectedIncomplete bool
		expectedMsg        string
	}
	type testCase struct {
		name   string
		chunks []testChunk
	}
	testCases := []testCase{
		{
			name: "valid",
			chunks: []testChunk{{
				octets:             []byte{0b11010000},
				expectedIncomplete: true,
			}, {
				octets: []byte{0b10100000, 0b01000000, 0b11101000, 0b10100000, 0b10010000},
			}, {
				octets: []byte("毧큰㴖탺㞹茦㑾쬘쟞臏棄䔲巋儤塹륟තⳈ䍉᱿▇ㆦ㍂矖๚䓂໘྽ꉌ㰃遀Ⰱ⑦죛ō⚣➻鍏쪈䵤빠꞊蛕캝ᩆⷺ᷒෪ꉥⰥ駽去씥既偩伌௡쾐촍纒ᚢ톍耟砍룄ࣤೡ栬Ɥ龵卾逭釭運㙣鳍䏬䄾À⪶ো吟㩜⪸ఌ鏘䓧鏈횴㒄ꀨ栺난荨솥踌᭍퀰乴裀"),
			}},
		},
		{
			name: "invalid",
			chunks: []testChunk{{
				octets:             []byte{0b11010000},
				expectedIncomplete: true,
			}, {
				octets:          []byte{0b01000000},
				expectedInvalid: true,
				expectedMsg:     "character started before previous finished",
			}, {
				octets:          []byte{0b10100000, 0b01000000},
				expectedInvalid: true,
				expectedMsg:     "previously found invalid",
			}},
		},
		{
			name: "invalidAutobahn631",
			chunks: []testChunk{{
				octets: []byte{0xce, 0xba, 0xe1, 0xbd, 0xb9, 0xcf, 0x83, 0xce, 0xbc, 0xce, 0xb5},
			}, {
				octets:          []byte{0xed, 0xa0, 0x80},
				expectedInvalid: true,
				expectedMsg:     "utf16 surrogate half",
			}},
		},
		{
			name: "validAutobahn624",
			chunks: []testChunk{{
				octets: []byte{0xce, 0xba, 0xe1, 0xbd, 0xb9, 0xcf, 0x83, 0xce, 0xbc, 0xce, 0xb5},
			}},
		},
		{
			name: "invalidOverlongAscii",
			chunks: []testChunk{{
				octets:          []byte{0b11000001, 0b10010000},
				expectedInvalid: true,
				expectedMsg:     "overlong ascii character",
			}},
		},
		{
			name: "invalidOverlong",
			chunks: []testChunk{{
				octets:          []byte{0b11100000, 0b10011111, 0b10111111},
				expectedInvalid: true,
				expectedMsg:     "overlong character",
			}},
		},
		// AI generated tests from here on
		{
			name: "validAsciiOnly",
			chunks: []testChunk{{
				octets: []byte("Hello, World! 0123456789"),
			}},
		},
		{
			name: "validEmptyInput",
			chunks: []testChunk{{
				octets: []byte{},
			}},
		},
		{
			name: "validSingleByteAtATime2Byte",
			chunks: []testChunk{{
				octets:             []byte{0xc3}, // first byte of ã (U+00E3) -> 0xC3 0xA3
				expectedIncomplete: true,
			}, {
				octets: []byte{0xa3},
			}},
		},
		{
			name: "validSingleByteAtATime3Byte",
			chunks: []testChunk{{
				octets:             []byte{0xe2}, // first byte of € (U+20AC) -> 0xE2 0x82 0xAC
				expectedIncomplete: true,
			}, {
				octets:             []byte{0x82},
				expectedIncomplete: true,
			}, {
				octets: []byte{0xac},
			}},
		},
		{
			name: "validSingleByteAtATime4Byte",
			chunks: []testChunk{{
				octets:             []byte{0xf0}, // first byte of 𐍈 (U+10348) -> 0xF0 0x90 0x8D 0x88
				expectedIncomplete: true,
			}, {
				octets:             []byte{0x90},
				expectedIncomplete: true,
			}, {
				octets:             []byte{0x8d},
				expectedIncomplete: true,
			}, {
				octets: []byte{0x88},
			}},
		},
		{
			name: "valid4ByteMaxUnicode",
			chunks: []testChunk{{
				// U+10FFFF — the maximum valid Unicode code point
				octets: []byte{0xf4, 0x8f, 0xbf, 0xbf},
			}},
		},
		{
			name: "invalidExceedsMaxUnicode",
			chunks: []testChunk{{
				// U+110000 — one beyond max
				octets:          []byte{0xf4, 0x90, 0x80, 0x80},
				expectedInvalid: true,
				expectedMsg:     "exceeds max unicode value",
			}},
		},
		{
			name: "invalidExceedsMaxUnicodeHighBits",
			chunks: []testChunk{{
				// F5 xx xx xx would encode U+140000+, well above max
				octets:          []byte{0xf5, 0x80, 0x80, 0x80},
				expectedInvalid: true,
				expectedMsg:     "exceeds max unicode value",
			}},
		},
		{
			name: "invalid5ByteSequence",
			chunks: []testChunk{{
				// 5-byte sequence (0b11111000)
				octets:          []byte{0xf8, 0x80, 0x80, 0x80, 0x80},
				expectedInvalid: true,
				expectedMsg:     "character is more than 4 bytes",
			}},
		},
		{
			name: "invalid6ByteSequence",
			chunks: []testChunk{{
				// 6-byte sequence (0b11111100)
				octets:          []byte{0xfc, 0x80, 0x80, 0x80, 0x80, 0x80},
				expectedInvalid: true,
				expectedMsg:     "character is more than 4 bytes",
			}},
		},
		{
			name: "invalidContinuationWithoutStart",
			chunks: []testChunk{{
				octets:          []byte{0x80},
				expectedInvalid: true,
				expectedMsg:     "continuation octet without character start",
			}},
		},
		{
			name: "invalidMultipleContinuationsWithoutStart",
			chunks: []testChunk{{
				octets:          []byte{0xbf},
				expectedInvalid: true,
				expectedMsg:     "continuation octet without character start",
			}},
		},
		{
			name: "invalidSurrogateMin",
			chunks: []testChunk{{
				// U+D800 (minimum surrogate) -> 0xED 0xA0 0x80
				octets:          []byte{0xed, 0xa0, 0x80},
				expectedInvalid: true,
				expectedMsg:     "utf16 surrogate half",
			}},
		},
		{
			name: "invalidSurrogateMax",
			chunks: []testChunk{{
				// U+DFFF (maximum surrogate) -> 0xED 0xBF 0xBF
				octets:          []byte{0xed, 0xbf, 0xbf},
				expectedInvalid: true,
				expectedMsg:     "utf16 surrogate half",
			}},
		},
		{
			name: "validJustBelowSurrogateRange",
			chunks: []testChunk{{
				// U+D7FF -> 0xED 0x9F 0xBF
				octets: []byte{0xed, 0x9f, 0xbf},
			}},
		},
		{
			name: "validJustAboveSurrogateRange",
			chunks: []testChunk{{
				// U+E000 -> 0xEE 0x80 0x80
				octets: []byte{0xee, 0x80, 0x80},
			}},
		},
		{
			name: "invalidOverlongAsciiNull",
			chunks: []testChunk{{
				// Overlong encoding of U+0000 as 2 bytes: 0xC0 0x80
				octets:          []byte{0xc0, 0x80},
				expectedInvalid: true,
				expectedMsg:     "overlong ascii character",
			}},
		},
		{
			name: "invalidOverlong3ByteSlash",
			chunks: []testChunk{{
				// Overlong encoding of '/' (U+002F) as 3 bytes: 0xE0 0x80 0xAF
				octets:          []byte{0xe0, 0x80, 0xaf},
				expectedInvalid: true,
				expectedMsg:     "overlong character",
			}},
		},
		{
			name: "invalidOverlong4Byte",
			chunks: []testChunk{{
				// Overlong 4-byte encoding: 0xF0 0x80 0x80 0xAF
				octets:          []byte{0xf0, 0x80, 0x80, 0xaf},
				expectedInvalid: true,
				expectedMsg:     "overlong character",
			}},
		},
		{
			name: "validMinimal2Byte",
			chunks: []testChunk{{
				// U+0080 -> 0xC2 0x80 (smallest valid 2-byte char)
				octets: []byte{0xc2, 0x80},
			}},
		},
		{
			name: "validMinimal3Byte",
			chunks: []testChunk{{
				// U+0800 -> 0xE0 0xA0 0x80 (smallest valid 3-byte char)
				octets: []byte{0xe0, 0xa0, 0x80},
			}},
		},
		{
			name: "validMinimal4Byte",
			chunks: []testChunk{{
				// U+10000 -> 0xF0 0x90 0x80 0x80 (smallest valid 4-byte char)
				octets: []byte{0xf0, 0x90, 0x80, 0x80},
			}},
		},
		{
			name: "invalidTruncated3ByteMiddleOfStream",
			chunks: []testChunk{{
				// Start a 3-byte char, send only 2 bytes, then start ASCII
				octets:             []byte{0xe2, 0x82},
				expectedIncomplete: true,
			}, {
				octets:          []byte{0x41}, // 'A' — not a continuation byte
				expectedInvalid: true,
				expectedMsg:     "character started before previous finished",
			}},
		},
		{
			name: "invalidTruncated4ByteThenNewChar",
			chunks: []testChunk{{
				octets:             []byte{0xf0, 0x90},
				expectedIncomplete: true,
			}, {
				octets:          []byte{0xc3, 0xa3}, // Start of new 2-byte char
				expectedInvalid: true,
				expectedMsg:     "character started before previous finished",
			}},
		},
		{
			name: "validMultipleCompleteCharsInOneChunk",
			chunks: []testChunk{{
				// "Ωé𐍈" — mix of 2, 2, and 4 byte characters
				octets: []byte{0xce, 0xa9, 0xc3, 0xa9, 0xf0, 0x90, 0x8d, 0x88},
			}},
		},
		{
			name: "validSplitAcrossManyChunks",
			chunks: []testChunk{{
				// "€" (U+20AC) = E2 82 AC, split one byte per chunk
				octets:             []byte{0xe2},
				expectedIncomplete: true,
			}, {
				octets:             []byte{0x82},
				expectedIncomplete: true,
			}, {
				// AC followed by ASCII 'X'
				octets: []byte{0xac, 0x58},
			}},
		},
		{
			name: "invalidSurrogateSplitAcrossChunks",
			chunks: []testChunk{{
				octets:             []byte{0xed},
				expectedIncomplete: true,
			}, {
				octets:          []byte{0xa0},
				expectedInvalid: true,
				expectedMsg:     "utf16 surrogate half",
			}},
		},
		{
			name: "validNonSurrogateSplitAcrossChunks",
			chunks: []testChunk{{
				// U+D7FF = ED 9F BF, split across chunks
				octets:             []byte{0xed},
				expectedIncomplete: true,
			}, {
				octets:             []byte{0x9f},
				expectedIncomplete: true,
			}, {
				octets: []byte{0xbf},
			}},
		},
		{
			name: "invalidMaxUnicodeSplitAcrossChunks",
			chunks: []testChunk{{
				// U+110000 = F4 90 80 80, split across chunks
				octets:             []byte{0xf4},
				expectedIncomplete: true,
			}, {
				octets:          []byte{0x90},
				expectedInvalid: true,
				expectedMsg:     "exceeds max unicode value",
			}},
		},
		{
			name: "validBoundaryU+FFFF",
			chunks: []testChunk{{
				// U+FFFF -> 0xEF 0xBF 0xBF (max 3-byte code point)
				octets: []byte{0xef, 0xbf, 0xbf},
			}},
		},
		{
			name: "validReplacementCharacter",
			chunks: []testChunk{{
				// U+FFFD (replacement character) -> 0xEF 0xBF 0xBD
				octets: []byte{0xef, 0xbf, 0xbd},
			}},
		},
		{
			name: "validNullByte",
			chunks: []testChunk{{
				octets: []byte{0x00},
			}},
		},
		{
			name: "validMaxAscii",
			chunks: []testChunk{{
				octets: []byte{0x7f},
			}},
		},
		{
			name: "incompleteAtEnd2Byte",
			chunks: []testChunk{{
				octets:             []byte{0xc2},
				expectedIncomplete: true,
			}},
		},
		{
			name: "incompleteAtEnd3Byte",
			chunks: []testChunk{{
				octets:             []byte{0xe0, 0xa0},
				expectedIncomplete: true,
			}},
		},
		{
			name: "incompleteAtEnd4Byte",
			chunks: []testChunk{{
				octets:             []byte{0xf0, 0x90, 0x80},
				expectedIncomplete: true,
			}},
		},
		{
			name: "invalidOverlongAsciiC0",
			chunks: []testChunk{{
				// 0xC0 0xAF — overlong encoding of '/'
				octets:          []byte{0xc0, 0xaf},
				expectedInvalid: true,
				expectedMsg:     "overlong ascii character",
			}},
		},
		{
			name: "invalidFEByte",
			chunks: []testChunk{{
				// 0xFE is never valid in UTF-8
				octets:          []byte{0xfe},
				expectedInvalid: true,
				expectedMsg:     "character is more than 4 bytes",
			}},
		},
		{
			name: "invalidFFByte",
			chunks: []testChunk{{
				// 0xFF is never valid in UTF-8
				octets:          []byte{0xff},
				expectedInvalid: true,
				expectedMsg:     "character is more than 4 bytes",
			}},
		},
	}

	validator := NewUtf8StreamValidator()

	for _, test := range testCases {
		for _, chunk := range test.chunks {
			valid, errMsg := validator.Validate(chunk.octets)
			expectedValid := !chunk.expectedInvalid
			if valid != expectedValid {
				t.Errorf("test '%s': expected valid = %v, got valid = %v", test.name, expectedValid, valid)
			}
			if valid == true {
				complete := validator.IsComplete()
				expectedComplete := !chunk.expectedIncomplete
				if complete != expectedComplete {
					t.Errorf("test '%s': expected IsComplete = %v, got IsComplete = %v", test.name, expectedComplete, complete)
				}
			} else {
				if errMsg != chunk.expectedMsg {
					t.Errorf("test '%s': expected msg = %v, got msg = %v", test.name, chunk.expectedMsg, errMsg)
				}
			}
		}
		validator.Reset()
	}
}
