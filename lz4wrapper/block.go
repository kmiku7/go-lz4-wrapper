/*
Package go_lz4_wrapper implements Decompress/Compress function to interoperate
with python-lz4 library, with the data format which has the size of
uncompressed data stored at the start.
*/
package lz4wrapper

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

// Decompress decompress data which compressed using lz4 block format and
// include the size of uncompressed data at the start.
func Decompress(src []byte) ([]byte, error) {
	if len(src) <= headerSize {
		return nil, errors.New("input too short")
	}
	originalDataLength := binary.LittleEndian.Uint32(src[0:4])
	if originalDataLength > uint32(math.MaxInt32) {
		return nil, fmt.Errorf("invalid size in header: 0x%x", src[0:4])
	}
	if originalDataLength == 0 {
		if len(src) != 5 {
			return nil, fmt.Errorf("coruppted size in header: 0x%x", src[0:4])
		}
		if src[4] != byte(0) {
			return nil, fmt.Errorf("coruppted literal size: 0x%x", src[4])
		}
		return nil, nil
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

// Compress compress data using lz4 block format and append the size of
// uncompressed data at head.
func Compress(src []byte) ([]byte, error) {
	var ht [1 << 16]int
	dstBuffer := make([]byte, headerSize+lz4.CompressBlockBound(len(src)))

	binary.LittleEndian.PutUint32(dstBuffer, uint32(len(src)))
	comLen, err := lz4.CompressBlock(src, dstBuffer[4:], ht[:])
	if err != nil {
		return nil, err
	}
	if comLen == 0 {
		dstOffset := headerSize
		srcLen := len(src)
		if srcLen < 0xF {
			dstBuffer[dstOffset] = byte(srcLen << 4)
			dstOffset++
		} else {
			dstBuffer[dstOffset] = 0xF0
			dstOffset++
			for srcLen -= 0xF; srcLen >= 0xFF; srcLen -= 0xFF {
				dstBuffer[dstOffset] = 0xFF
				dstOffset++
			}
			dstBuffer[dstOffset] = byte(srcLen)
			dstOffset += 1
		}
		copy(dstBuffer[dstOffset:], src)
		comLen = dstOffset + len(src) - headerSize
	}
	if (headerSize + comLen) < (len(dstBuffer) * 2 / 3) {
		newDstBuffer := make([]byte, headerSize+comLen)
		copy(newDstBuffer, dstBuffer)
		dstBuffer = newDstBuffer
	}
	return dstBuffer[0 : headerSize+comLen], nil
}
