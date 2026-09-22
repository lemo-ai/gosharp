package gbinary_test

import (
	"testing"

	"github.com/lemo-ai/gosharp/encoding/gbinary"
	"github.com/stretchr/testify/assert"
)

func TestEncodeIntNegative(t *testing.T) {
	b := gbinary.LeEncodeInt(-1000)
	assert.Equal(t, 2, len(b), " -1000 should encode as int16")
	assert.Equal(t, int16(-1000), gbinary.LeDecodeToInt16(b))

	b64 := gbinary.LeEncodeInt(int(1)<<40 + 1)
	assert.Equal(t, 8, len(b64))

	be := gbinary.BeEncodeInt(-1000)
	assert.Equal(t, 2, len(be))
	assert.Equal(t, int16(-1000), gbinary.BeDecodeToInt16(be))
}

func TestEncodeBitsToBytesPadding(t *testing.T) {
	bits := []gbinary.Bit{1}
	b := gbinary.EncodeBitsToBytes(bits)
	assert.Equal(t, 1, len(b))
	assert.Equal(t, byte(0b10000000), b[0])

	bits = []gbinary.Bit{1, 0, 1}
	b = gbinary.EncodeBitsToBytes(bits)
	assert.Equal(t, 1, len(b))
	assert.Equal(t, byte(0b10100000), b[0])
}

func TestDecodeEmpty(t *testing.T) {
	assert.Equal(t, uint8(0), gbinary.LeDecodeToUint8(nil))
	assert.Equal(t, int8(0), gbinary.LeDecodeToInt8([]byte{}))
	assert.Equal(t, uint8(0), gbinary.BeDecodeToUint8(nil))
	assert.Equal(t, int8(0), gbinary.BeDecodeToInt8([]byte{}))
}

func TestLeRoundTrip(t *testing.T) {
	assert.Equal(t, int32(123456), gbinary.LeDecodeToInt32(gbinary.LeEncodeInt32(123456)))
	assert.Equal(t, float64(3.14), gbinary.LeDecodeToFloat64(gbinary.LeEncodeFloat64(3.14)))
}
