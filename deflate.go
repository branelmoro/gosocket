package main

import (
    "fmt"
	"compress/flate"
	"bytes"
	"io"
    "reflect"
    // "strings"
)

func compress(src []byte, level int) []byte {
    compressed_buffer := new(bytes.Buffer)
    compressor, _ := flate.NewWriter(compressed_buffer, level)
    compressor.Write(src)
    compressor.Close()
    return compressed_buffer.Bytes()
}

func decompress(src []byte) []byte {

    compressed_buffer := new(bytes.Buffer) 
    compressed_buffer.Write(src)

    decompressed_buffer := new(bytes.Buffer)

    decompressor := flate.NewReader(compressed_buffer)
    // decompressor.(flate.Limiter).SetReadLimit(compressed_buffer.Len())
    _, err := io.Copy(decompressed_buffer, decompressor)
    fmt.Println("Fmt err is---------", err)
    decompressor.Close()
    return decompressed_buffer.Bytes()
}

func main() {
    if false {
        brokenReads()
        return
    }
    var(
        originCompressed []byte
        compressedData []byte
        data []byte
    )

    // fmt.Println("My decomressed data - ", decompress([]byte{34,77,53,0,0,0,255,255}))

    // originCompressed = []byte("5sfsdfs wef eg e4 g23r h3 er")
    originCompressed = []byte{50,45,78,43,78,73,43,86,40,79,77,83,72,77,87,72,53,81,72,55,50,46,82,200,48,86,139,83,85,81,118,80,84,86,81,141,83,85,118,80,1,51,212,0,0,0,0,255,255,50,37,73,53,0,0,0,255,255}
    fmt.Println("original comressed data - ", originCompressed)

    data = decompress(originCompressed)

    fmt.Println("original decomressed data", data, string(data))

    // return

    i := -2
    for false {
        compressedData = compress(data, i)
        if reflect.DeepEqual(originCompressed, compressedData) {
            fmt.Println("Match found at ---------", i)
        }
        fmt.Println("original comressed---", originCompressed)
        fmt.Println("My comressed---------", compressedData, i, "\r\n")
        if i == 3 {
            break
        }
        i += 1
    }







    // var decompressorReader io.Reader

    compressed_buffer := new(bytes.Buffer)
    decompressed_buffer := new(bytes.Buffer)
    compressor, _ := flate.NewWriter(compressed_buffer, 9)
    decompressor := flate.NewReader(compressed_buffer)

    src := []byte("5sfsdfs wef eg e4 g23r h3&^%$#@!#$%^%#@$!#$%^&")
    fmt.Println("This comressed data - ", compress(src, 9))
    i = 0
    for {
        fmt.Println("\r\n\r\nString To Compress---", string(src))
        if i == 3 {
            compressed_buffer.Write([]byte{50,37,73,53,0,0,0,255,255})
        } else {
            compressor.Write(src)
            compressor.Flush()
        }
        fmt.Println("comressed data is---", compressed_buffer.Bytes(), compressed_buffer.Len())


        bf := new(bytes.Buffer)
        bf.Write(compressed_buffer.Bytes())
        // decompressor.(flate.Resetter).Reset(bf, nil)
        // decompressor.(flate.Limiter).SetReadLimit(bf.Len())
        _, err := io.Copy(decompressed_buffer, decompressor)
        fmt.Println("decompressed data is -----------", string(decompressed_buffer.Bytes()), err)

        decompressed_buffer.Reset()
        compressed_buffer.Reset()


        if i == 5 {
            break
        }
        i += 1

    }
    // decompressor.Close()
    compressor.Close()








}

func brokenReads() {


    // broken reads test
    compressed_buffer := new(bytes.Buffer)
    decompressed_buffer := new(bytes.Buffer)
    // compressor, _ := flate.NewWriter(compressed_buffer, 9)
    decompressor := flate.NewReader(compressed_buffer)


    compressedData := []byte{50,45,78,43,78,73,43,86,40,79,77,83,72,77,87,72,53,81,72,55,50,46,82,200,48,86,139,83,85,81,118,80,84,86,81,141,83,85,118,80,1,51,212,0,1,0,0,255,255}


    fmt.Println("comressed data is---", compressedData)

    uncompressedData := []byte{}

    compressed_buffer.Write(compressedData[:5])
    // compressed_buffer.Write([]byte{0,0,255,255})
    // decompressor.(flate.Resetter).ResetReader(compressed_buffer)
    // decompressor.(flate.Limiter).SetReadLimit(compressed_buffer.Len())
    _, err := io.Copy(decompressed_buffer, decompressor)
    fmt.Println("decompressed data  is -err ----------", err)
    uncompressedData = append(uncompressedData, decompressed_buffer.Bytes()...)
    decompressed_buffer.Reset()
    compressed_buffer.Reset()



    compressed_buffer.Write(compressedData[5:15])
    // compressed_buffer.Write([]byte{0,0,255,255})
    // decompressor.(flate.Resetter).ResetReader(compressed_buffer)
    // decompressor.(flate.Limiter).AddReadLimit(compressed_buffer.Len())
    _, err = io.Copy(decompressed_buffer, decompressor)
    fmt.Println("decompressed data  is ---err ----------", err)
    uncompressedData = append(uncompressedData, decompressed_buffer.Bytes()...)
    decompressed_buffer.Reset()
    compressed_buffer.Reset()



    compressed_buffer.Write(compressedData[15:])
    // compressed_buffer.Write([]byte{0,0,255,255})
    // decompressor.(flate.Resetter).ResetReader(compressed_buffer)
    // decompressor.(flate.Limiter).AddReadLimit(compressed_buffer.Len())
    _, err = io.Copy(decompressed_buffer, decompressor)
    fmt.Println("decompressed data  is -------err ----------", err)
    uncompressedData = append(uncompressedData, decompressed_buffer.Bytes()...)
    decompressed_buffer.Reset()
    compressed_buffer.Reset()


    fmt.Println("decompressed data is -----------", string(uncompressedData), err)
}