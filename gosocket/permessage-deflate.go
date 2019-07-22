package gosocket

import (
	"gosocket/flate"
	"bytes"
	"io"
)

type perMessageDeflate struct {
    c *flate.Writer
    d io.ReadCloser

}

func (p *perMessageDeflate) compress(src []byte) ([]byte, error) {
    var(
        n int
        err error
    )
    compressedWriter := new(bytes.Buffer)
    p.c.ResetWriter(compressedWriter)
    for {
        n, err = p.c.Write(src)
        if err != nil {
            return nil, err
        }
        if n == len(src) {
            break
        }
        src = src[n:]
    }
    err = p.c.Flush()
    if err != nil {
        return nil, err
    }
    p.c.ResetWriter(nil)
    return compressedWriter.Bytes(), nil
}

func (p *perMessageDeflate) decompress(src []byte) ([]byte, error) {
    compressedReader := bytes.NewBuffer(src)
    compressedReader.Write([]byte{0x0, 0x0, 0xff, 0xff}) // append dataBlock finished bytes
    p.d.(flate.Resetter).Reset(compressedReader, nil)
    p.d.(flate.Limiter).SetReadLimit(compressedReader.Len())
    decompressedReader := new(bytes.Buffer)
    _, err := io.Copy(decompressedReader, p.d)
    if err != nil {
        if _, ok := err.(flate.ReadLimitReachedError); !ok {
            return nil, err
        }
    }
    p.d.(flate.Resetter).Reset(nil, nil)
    return decompressedReader.Bytes(), nil
}

func (p *perMessageDeflate) close() error {
    p.c.Close()
    return p.d.Close()
}

func newCompressor(level int) (*flate.Writer, error) {
    compressor, err := flate.NewWriter(nil, level)
    if err != nil {
        return nil, err
    }
    return compressor, err
}

func newDecompressor() io.ReadCloser {
    return flate.NewReader(nil)
}

func newPerMessageDeflate(level int) (*perMessageDeflate, error) {
    compressor, err := newCompressor(level)
    if err != nil {
        return nil, err
    }
    return &perMessageDeflate{
        c:compressor,
        d:newDecompressor(),
    }, nil
}
