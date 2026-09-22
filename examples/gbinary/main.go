package main

import (
	"fmt"

	"github.com/lemo-ai/gosharp/encoding/gbinary"
)

func main() {
	b := gbinary.LeEncodeInt32(123456)
	fmt.Println("LeEncodeInt32:", b)
	fmt.Println("LeDecodeToInt32:", gbinary.LeDecodeToInt32(b))

	neg := gbinary.LeEncodeInt(-1000)
	fmt.Println("LeEncodeInt(-1000) len:", len(neg), "value:", gbinary.LeDecodeToInt16(neg))

	bits := gbinary.EncodeBits(nil, 0b1010, 4)
	fmt.Println("EncodeBitsToBytes:", gbinary.EncodeBitsToBytes(bits))
}
