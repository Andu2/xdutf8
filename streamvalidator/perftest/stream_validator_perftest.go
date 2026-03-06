package main

import (
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/andu2/xdutf8/streamvalidator"
)

func main() {
	tests := []struct {
		name string
		data []byte
	}{{
		name: "ASCII",
		data: []byte("mK#9vQ!zR@2pL$wX&8nJ%4tY^6fH*0dG~3sB+7cA=1eM<5uN>9oP|2iC-6kT/4lV.7bW,8gZ?3hD;0jF:5rE"),
	}, {
		name: "Mixed Length Unicode",
		data: []byte("㸽㢐╧̨ƺǙ𐁝𐊘ęƋ𐋤᫃ᣭơỹ𐌊ÙΝĲᶊ࿽ሷηᗇɰ𐂒ࡪ᷆𐁁˾𐎠ᇳÛĨˮ㮒𐂎𐈢𐇕𐊆̛⩑𐎲ɭʹ𐁲𐅎𐂣˅⌌𐋫𐃎ഋ㯍㮧Ŏ𐈹ᧈ͎⺣㾢à𐌉㚳ຊ෯𐈤𐅍ƮƟ̜Ƥζ𐇁⭔𐅞Ύ㣝ะ𐋶ḃ𐇜𐇵𐆉𐇒ʋÞǜğે𐋗⁐‹𐆭𐌵㒮𐈨𐉐͙β"),
	}, {
		name: "Mixed ASCII and Unicode",
		data: []byte("𐄽#~p6T-=)VH;2ɠ7{𐏤<$0Ū𐍢>E㽖fGᖮv*ó[R:'bhoσ↪y(UAzLC𐌟ீOX}|t_✔?suS˩%gÏ⛎K.⁁]nƲPd9lMDğ͛Zm@e኉!&/Q𐎣Y𐀄a𐉈I𐈟Fcixr"),
	}}

	const iterations = 100_000

	for _, test := range tests {
		start := time.Now()
		for i := 0; i < iterations; i++ {
			validator := streamvalidator.NewUtf8StreamValidator()
			validator.Validate(test.data)
		}
		streamDuration := time.Since(start)

		start = time.Now()
		for i := 0; i < iterations; i++ {
			utf8.Valid(test.data)
		}
		stdlibDuration := time.Since(start)

		fmt.Printf("---%s---\n", test.name)
		fmt.Printf("Input length:     %d bytes\n", len(test.data))
		fmt.Printf("Iterations:       %d\n", iterations)
		fmt.Printf("Stream validator: %v (%.0f ns/op)\n", streamDuration, float64(streamDuration.Nanoseconds())/iterations)
		fmt.Printf("Standard library: %v (%.0f ns/op)\n", stdlibDuration, float64(stdlibDuration.Nanoseconds())/iterations)
		fmt.Printf("Ratio:            %.2fx\n", float64(stdlibDuration)/float64(streamDuration))
	}
}
