package lz4wrapper

import (
	"bytes"
	"testing"
)

func operationPair(t *testing.T, caseName string, originalData []byte, expectCompressedData []byte) {

	compressedData, err := Compress(originalData)
	if err != nil {
		t.Errorf("Case name: %v, compress failed, err: %v.", caseName, err)
	}
	if expectCompressedData != nil && !bytes.Equal(compressedData, expectCompressedData) {
		t.Errorf("Case name: %v, compressed data not match.", caseName)
	}
	decompressedData, err := Decompress(compressedData)
	if err != nil {
		t.Errorf("Case name: %v, decompress failed, err: %v.", caseName, err)
	}
	if !bytes.Equal(decompressedData, originalData) {
		t.Errorf("Case name: %v, decompressed data not match.", caseName)
	}
}

func TestEmptyString(t *testing.T) {
	operationPair(t, "CaseEmptyString", nil, []byte{0x0, 0x0, 0x0, 0x0, 0x0})
}

func TestNormalString(t *testing.T) {
	operationPair(t, "CaseNormalString",
		[]byte("123456789012345678901234567890123456789012345678901234567890"), nil)
}

func TestNotCompressibleString(t *testing.T) {
	operationPair(t, "CaseNotCompressibleStringLen1", []byte("1"),
		[]byte{0x1, 0x0, 0x0, 0x0, 0x10, '1'})

	operationPair(t, "CaseNotCompressibleStringLen10", []byte("1234567890"),
		[]byte{0x0a, 0x0, 0x0, 0x0, 0xa0, '1', '2', '3', '4', '5', '6', '7', '8', '9', '0'})
}
