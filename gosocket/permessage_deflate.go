package gosocket

import (
	"compress/flate"
	"bytes"
	"io"
)

func compress(src *[]byte, level int) []byte {
    compressed_buffer := new(bytes.Buffer)
    compressor, _ := flate.NewWriter(compressed_buffer, level)
    compressor.Write(*src)
    compressor.Close()
    return compressed_buffer.Bytes()
}

func decompress(src *[]byte) []byte {

    compressed_buffer := new(bytes.Buffer) 
    compressed_buffer.Write(*src)

    decompressed_buffer := new(bytes.Buffer)

    decompressor := flate.NewReader(compressed_buffer)
    io.Copy(decompressed_buffer, decompressor)
    decompressor.Close()
    return decompressed_buffer.Bytes()
}
