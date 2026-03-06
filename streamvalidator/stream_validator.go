package streamvalidator

type streamValidatorState struct {
	bytesLeft           uint8
	checkReserved       bool
	checkMaxUnicode     bool
	isInvalid           bool
	overlongBitsToCheck uint8
}

func NewUtf8StreamValidator() *streamValidatorState {
	return &streamValidatorState{}
}

func (state *streamValidatorState) Reset() {
	state.bytesLeft = 0
	state.checkReserved = false
	state.checkMaxUnicode = false
	state.isInvalid = false
	state.overlongBitsToCheck = 0
}

func (state *streamValidatorState) IsComplete() bool {
	return state.bytesLeft == 0
}

func (state *streamValidatorState) Validate(octets []byte) (bool, string) {
	if state.isInvalid {
		return false, "previously found invalid"
	}

	i := 0
	for i < len(octets) {
		if state.bytesLeft == 0 {
			// Check ASCII first
			// Fast path: Try to skip runs of 8 bytes at a time
			for i+8 <= len(octets) {
				// The following is copied from the standard library at unicode/utf8/utf8.go

				// Combining two 32 bit loads allows the same code to be used
				// for 32 and 64 bit platforms.
				// The compiler can generate a 32bit load for first32 and second32
				// on many platforms. See test/codegen/memcombine.go.
				first32 := uint32(octets[i+0]) | uint32(octets[i+1])<<8 | uint32(octets[i+2])<<16 | uint32(octets[i+3])<<24
				second32 := uint32(octets[i+4]) | uint32(octets[i+5])<<8 | uint32(octets[i+6])<<16 | uint32(octets[i+7])<<24
				if (first32|second32)&0x80808080 != 0 {
					break
				}

				i += 8
			}

			// Continue to try to skip through ASCII without having to check bytesLeft again
			for i < len(octets) && octets[i]&0b10000000 == 0 {
				i++
			}
			if i >= len(octets) {
				break
			}

			octet := octets[i]
			if octet&0b11000000 == 0b10000000 {
				state.isInvalid = true
				return false, "continuation octet without character start"
			} else {
				// start of char
				if octet&0b11100000 == 0b11000000 { // length 2
					if octet&0b00011110 == 0 {
						// This is the only overlong case that we can determine from the first byte
						// (also, the last bit may be 1 and it is still overlong)
						state.isInvalid = true
						return false, "overlong ascii character"
					}

					state.bytesLeft = 1
				} else if octet&0b11110000 == 0b11100000 { // length 3
					state.bytesLeft = 2
					if octet&0b00001111 == 0 {
						state.overlongBitsToCheck = 1
					}
					if octet == 0b11101101 {
						// The unicode values U+D800 through U+DFFF are reserved for utf-16 and invalid in utf-8.
						// If we have a 3 byte char and the starting four bits are 1101 (D), we need to check the next bit.
						// If the next bit of the char is a 1, it's invalid.
						state.checkReserved = true
					}
				} else if octet&0b11111000 == 0b11110000 { // length 4
					if octet&0b00000100 == 0b00000100 {
						// If any of the next 4 bits are 1, it's invalid because it's too large
						// We can check the first two here, but have to set a flag to check the next two in the next byte
						if octet&0b00000011 > 0 {
							state.isInvalid = true
							return false, "exceeds max unicode value"
						}
						state.checkMaxUnicode = true
					}

					state.bytesLeft = 3
					if octet&0b00000111 == 0 {
						state.overlongBitsToCheck = 2
					}
				} else {
					// utf-8 was designed to allow up to 6 byte characters,
					// but it was later decided that it is only allowed to go up to U+10FFFF
					// which only requires 4 bytes
					state.isInvalid = true
					return false, "character is more than 4 bytes"
				}
			}
		} else {
			octet := octets[i]
			// continuation
			if octet&0b11000000 != 0b10000000 {
				state.isInvalid = true
				return false, "character started before previous finished"
			}
			// Doing this makes branch prediction better, I guess. It makes it slightly faster
			if state.checkReserved || state.checkMaxUnicode || state.overlongBitsToCheck > 0 {
				if state.checkReserved {
					if octet&0b00100000 == 0b00100000 {
						// This is invalid because it is a reserved value. See above where checkReserved is set to true.
						state.isInvalid = true
						return false, "utf16 surrogate half"
					}
					state.checkReserved = false
				}
				if state.checkMaxUnicode {
					if octet&0b00110000 > 0 {
						state.isInvalid = true
						return false, "exceeds max unicode value"
					}
					state.checkMaxUnicode = false
				}
				if state.overlongBitsToCheck > 0 {
					var checkMask byte
					if state.overlongBitsToCheck == 1 {
						checkMask = 0b00100000
					} else if state.overlongBitsToCheck == 2 {
						checkMask = 0b00110000
					}
					if octet&checkMask == 0 {
						state.isInvalid = true
						return false, "overlong character"
					}
					state.overlongBitsToCheck = 0
				}
			}
			state.bytesLeft -= 1
		}
		i++
	}

	return true, ""
}
