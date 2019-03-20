package go_lz4_wrapper

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"

	"github.com/pierrec/lz4"
)

const (
	HEADER_SIZE = 4
)

func decompress(src []byte) ([]byte, error) {
	if len(src) < HEADER_SIZE {
		return nil, errors.New("input too short.")
	}
	originalDataLength := binary.LittleEndian.Uint32(src[0:4])
	if originalDataLength > math.MaxInt32 {
		return nil, fmt.Errorf("invalid size in header: 0x%x", src[0:4])
	}
	dst := make([]byte, originalDataLength)
	decompressLength, err := lz4.UncompressBlock(src[4:], dst)
	if err != nil {
		return nil, err
	}
	if decompressLength != originalDataLength {
		return nil, fmt.Errorf("coruppted size in header: 0x%x", src[0:4])
	}
	return dst, nil
}