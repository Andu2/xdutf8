package main

import (
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/andu2/xdutf8/streamvalidator"
)

func main() {
	testInput := []byte("The quick brown fox jumps over the lazy dog. 日本語テスト. Ünïcödé. 🔥🎉The quick brown fox jumps over the lazy dog. 日本語テスト. Ünïcödé. 🔥🎉The quick brown fox jumps over the lazy dog. 日本語テスト. Ünïcödé. 🔥🎉The quick brown fox jumps over the lazy dog. 日本語テスト. Ünïcödé. 🔥🎉The quick brown fox jumps over the lazy dog. 日本語テスト. Ünïcödé. 🔥🎉The quick brown fox jumps over the lazy dog. 日本語テスト. Ünïcödé. 🔥🎉The quick brown fox jumps over the lazy dog. 日本語テスト. Ünïcödé. 🔥🎉The quick brown fox jumps over the lazy dog. 日本語テスト. Ünïcödé. 🔥🎉The quick brown fox jumps over the lazy dog. 日本語テスト. Ünïcödé. 🔥🎉")

	const iterations = 100_000

	start := time.Now()
	for i := 0; i < iterations; i++ {
		validator := streamvalidator.NewUtf8StreamValidator()
		validator.Validate(testInput)
	}
	streamDuration := time.Since(start)

	start = time.Now()
	for i := 0; i < iterations; i++ {
		utf8.Valid(testInput)
	}
	stdlibDuration := time.Since(start)

	fmt.Printf("Input length:     %d bytes\n", len(testInput))
	fmt.Printf("Iterations:       %d\n", iterations)
	fmt.Printf("Stream validator: %v (%.0f ns/op)\n", streamDuration, float64(streamDuration.Nanoseconds())/iterations)
	fmt.Printf("Standard library: %v (%.0f ns/op)\n", stdlibDuration, float64(stdlibDuration.Nanoseconds())/iterations)
	fmt.Printf("Ratio:            %.1fx\n", float64(stdlibDuration)/float64(streamDuration))
}
