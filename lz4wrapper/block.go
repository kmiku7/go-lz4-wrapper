/*
Package go_lz4_wrapper implements Decompress/Compress function to interoperate
with python-lz4 library, with the data format which has the size of
uncompressed data stored at the start.
*/
package go_lz4_wrapper

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"

	"github.com/pierrec/lz4"
)

const (
	headerSize = 4
)

// Decompress decompress data which compressed with lz4 block format and
// include the size of uncompressed data at the start.
func Decompress(src []byte) ([]byte, error) {
	if len(src) < headerSize {
		return nil, errors.New("input too short")
	}
	originalDataLength := binary.LittleEndian.Uint32(src[0:4])
	if originalDataLength > uint32(math.MaxInt32) {
		return nil, fmt.Errorf("invalid size in header: 0x%x", src[0:4])
	}
	dst := make([]byte, originalDataLength)
	decompressLength, err := lz4.UncompressBlock(src[4:], dst)
	if err != nil {
		return nil, err
	}
	if uint32(decompressLength) != originalDataLength {
		return nil, fmt.Errorf("coruppted size in header: 0x%x", src[0:4])
	}
	return dst, nil
}