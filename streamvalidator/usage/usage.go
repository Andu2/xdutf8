package main

import (
	"fmt"

	"github.com/andu2/xdutf8/streamvalidator"
)

func main() {

	testChunks := [][]byte{
		// 2-byte character missing the 2nd byte
		{0b11010000},

		// Continuation of previous byte, then an ASCII byte, then a full 3-byte character
		{0b10100000, 0b01000000, 0b11101000, 0b10100000, 0b10010000},

		// Overlong character - given as 3 bytes but could fit in 2
		{0b11100000, 0b10011111, 0b10111111},

		// Valid ASCII
		[]byte("aaaaa"),
	}

	validator := streamvalidator.NewUtf8StreamValidator()
	for _, chunk := range testChunks {
		valid, errMsg := validator.Validate(chunk)
		isComplete := validator.IsComplete()
		fmt.Printf("Bytes: %b, valid: %v, complete: %v, msg: %s\n", chunk, valid, isComplete, errMsg)
	}

	// output:
	// Bytes: [11010000], valid: true, complete: false, msg:
	// Bytes: [10100000 1000000 11101000 10100000 10010000], valid: true, complete: true, msg:
	// Bytes: [11100000 10011111 10111111], valid: false, complete: false, msg: overlong character
	// Bytes: [1100001 1100001 1100001 1100001 1100001], valid: false, complete: false, msg: previously found invalid

	fmt.Println("Resetting")
	validator.Reset()
	valid, errMsg := validator.Validate(testChunks[3])
	isComplete := validator.IsComplete()
	fmt.Printf("Bytes: %b, valid: %v, complete: %v, msg: %s\n", testChunks[3], valid, isComplete, errMsg)

	// output:
	// Bytes: [1100001 1100001 1100001 1100001 1100001], valid: true, complete: true, msg:
}
