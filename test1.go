package main

import "compress/flate"
// import "log"
import "bytes"
// import "os"
// import "io/ioutil"
import "io"
import "fmt"


// 0x48 0x65 0x6c 0x6c 0x6f (contains "Hello")
// 0xf2 0x48 0xcd 0xc9 0xc9 0x07 0x00 (compressed)

func main() {
    // inData, _ := ioutil.ReadFile("stuff.dat")

// [72 101 108 108 111] ---- hello
// [242 72 205 201 201 7 4 0 0 255 255] ---- compressed
// 0xf2 0x48 0xcd 0xc9 0xc9 0x07

    inData := []byte{104,101,108,108,111,32,102,114,111,109,32,98,114,111,119,115,101,114}

    compressedData := new(bytes.Buffer) 

    // compressedData := make([]byte, 100)
    compress(inData, compressedData, 9)
    // fmt.Println(byte(compressedData))
    fmt.Println(compressedData.Bytes())

    // ioutil.WriteFile("compressed.dat", compressedData.Bytes(), os.ModeAppend)

// hello from browser
    inData1 := []byte{202,72,205,201,201,87,72,43,202,207,85,72,42,202,47,47,78,45,2,0}
    compressedData1 := new(bytes.Buffer) 
    compressedData1.Write(inData1)

    deCompressedData := new(bytes.Buffer)
    decompress(compressedData, deCompressedData)
    // log.Print(deCompressedData)
    fmt.Println(deCompressedData.Bytes())
    fmt.Println(string(deCompressedData.Bytes()))
    // fmt.Println(byte(deCompressedData))
}
func compress(src []byte, dest io.Writer, level int) {
    compressor, _ := flate.NewWriter(dest, level)
    compressor.Write(src)
    compressor.Close()
}
func decompress(src io.Reader, dest io.Writer) {
    decompressor := flate.NewReader(src)
    io.Copy(dest, decompressor)
    decompressor.Close()
}