package streamvalidator

type streamValidatorState struct {
	bytesLeft           uint8
	checkReserved       bool
	checkMaxUnicode     bool
	isInvalid           bool
	overlongBitsToCheck uint8
}

const (
	bytesLeftMask           = 3 << 0
	overlongBitsToCheckMask = 3 << 2
	checkReservedMask       = 1 << 4
	checkMaxUnicodeMask     = 1 << 5
	isInvalidMask           = 1 << 6
)

func NewUtf8StreamValidator() *streamValidatorState {
	return &streamValidatorState{}
}

func (validator *streamValidatorState) Reset() {
	validator.bytesLeft = 0
	validator.checkReserved = false
	validator.checkMaxUnicode = false
	validator.isInvalid = false
	validator.overlongBitsToCheck = 0
}

func (validator *streamValidatorState) IsComplete() bool {
	return validator.bytesLeft == 0
}

func (validator *streamValidatorState) Validate(octets []byte) (bool, string) {
	if validator.isInvalid {
		return false, "previously found invalid"
	}

	for i := 0; i < len(octets); i++ {
		octet := octets[i]
		if validator.bytesLeft == 0 {
			if octet&0b10000000 == 0 {
				// ASCII
				continue
			} else if octet&0b11000000 == 0b10000000 {
				validator.isInvalid = true
				return false, "continuation octet without character start"
			} else {
				// start of char
				if octet&0b11100000 == 0b11000000 { // length 2
					if octet&0b00011110 == 0 {
						// This is the only overlong case that we can determine from the first byte
						// (also, the last bit may be 1 and it is still overlong)
						validator.isInvalid = true
						return false, "overlong ascii character"
					}

					validator.bytesLeft = 1
				} else if octet&0b11110000 == 0b11100000 { // length 3
					validator.bytesLeft = 2
					if octet&0b00001111 == 0 {
						validator.overlongBitsToCheck = 1
					}
					if octet == 0b11101101 {
						// The unicode values U+D800 through U+DFFF are reserved for utf-16 and invalid in utf-8.
						// If we have a 3 byte char and the starting four bits are 1101 (D), we need to check the next bit.
						// If the next bit of the char is a 1, it's invalid.
						validator.checkReserved = true
					}
				} else if octet&0b11111000 == 0b11110000 { // length 4
					if octet&0b00000100 == 0b00000100 {
						// If any of the next 4 bits are 1, it's invalid because it's too large
						// We can check the first two here, but have to set a flag to check the next two in the next byte
						if octet&0b00000011 > 0 {
							validator.isInvalid = true
							return false, "exceeds max unicode value"
						}
						validator.checkMaxUnicode = true
					}

					validator.bytesLeft = 3
					if octet&0b00000111 == 0 {
						validator.overlongBitsToCheck = 2
					}
				} else {
					// utf-8 was designed to allow up to 6 byte characters,
					// but it was later decided that it is only allowed to go up to U+10FFFF
					// which only requires 4 bytes
					validator.isInvalid = true
					return false, "character is more than 4 bytes"
				}
			}
		} else {
			// continuation
			if octet&0b11000000 != 0b10000000 {
				validator.isInvalid = true
				return false, "character started before previous finished"
			}
			if validator.checkReserved {
				if octet&0b00100000 == 0b00100000 {
					// This is invalid because it is a reserved value. See above where checkReserved is set to true.
					validator.isInvalid = true
					return false, "utf16 surrogate half"
				}
				validator.checkReserved = false
			}
			if validator.checkMaxUnicode {
				if octet&0b00110000 > 0 {
					validator.isInvalid = true
					return false, "exceeds max unicode value"
				}
				validator.checkMaxUnicode = false
			}
			if validator.overlongBitsToCheck > 0 {
				checkBits := octet << 1 >> (7 - validator.overlongBitsToCheck)
				if checkBits == 0 {
					validator.isInvalid = true
					return false, "overlong character"
				}
				validator.overlongBitsToCheck = 0
			}
			validator.bytesLeft -= 1
		}
	}

	return true, ""
}
