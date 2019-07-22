package main

import "compress/flate"
// import "log"
import "bytes"
// import "os"
// import "io/ioutil"
import "io"
import "fmt"
import "sync"
import "time"


// 0x48 0x65 0x6c 0x6c 0x6f (contains "Hello")
// 0xf2 0x48 0xcd 0xc9 0xc9 0x07 0x00 (compressed)


type serverConfig struct {
    a int
    b interface{}
}

func (p *serverConfig) AAA() int {
    return p.a
}


type pqrs struct {
    *serverConfig
    mux sync.Mutex
}

func (p *pqrs) getA() int {
    return p.a
}

func (p *pqrs) mutexTest(t time.Duration) {
    p.mux.Lock()

    fmt.Println("locked for", t)
    time.Sleep(t)

    p.mux.Unlock()
    fmt.Println("unlocked after", t)
}



func gr(l *chan bool) {
    lock := *l
    fmt.Println("lock----------", <- lock)
}

func switchCase() {

    // var pqr byte
    abc := "hahsikjakjfakfnk"

    sc := serverConfig{
        a:3452,
    }

    pqr := pqrs{
        serverConfig:&sc,
    }

    t1 := time.Now()
    t2 := t1.Add(3 * time.Second)
    t3 := t1.Add(1 * time.Second)

    go pqr.mutexTest(t2.Sub(t1))
    go pqr.mutexTest(t3.Sub(t1))

    fmt.Println("getA----------", pqr.a, pqr.getA(), pqr.AAA())

    lock := make(chan bool)
    go gr(&lock)

    lock <- true

    fmt.Println("here1----------", lock)

    // z := <- lock

    // fmt.Println("here1----------", z)
    // lock <- true

    // fmt.Println("here2----------")

    // lock <- true

    // fmt.Println("here3----------")

    c := "1"

    switch(c) {
        case "1":
            abc = "1sdfsdf"
            break
        case "2":
            abc = "2sdfsdf"
            break
        case "3":
            abc = "hgkuhkj"
            break
        default:
            abc = "default"
            break
    }

    fmt.Println(abc)
}




type Mutatable struct {
    a int
    b int
}

func (m Mutatable) StayTheSame() {
    m.a = 1
    m.b = 6
}

func (m *Mutatable) Mutate() {
    m.a = 5
    m.b = 7
}



func zyx(a int) interface{} {
    if a == 0 {
        return "sdgsdgfsdg"
    } else if a == 1 {
        return 19909
    } else {
        return false
    }
}



func main() {


    

    fmt.Println("zyx(0)--", zyx(0), ", zyx(1)---", zyx(1), ", zyx(2)---", zyx(2))

    switchCase()

    c := serverConfig{
        a:1,
        b:1,
    }

    fmt.Println(c)

    c.b = []int{45,44,67}

    fmt.Println(c)
    c.b ="sdsdsd"

    fmt.Println(c)

    cdc := serverConfig{
        a:1,
    }

    fmt.Println(cdc)


    // inData, _ := ioutil.ReadFile("stuff.dat")

// [72 101 108 108 111] ---- hello
// [242 72 205 201 201 7 4 0 0 255 255] ---- compressed
// 0xf2 0x48 0xcd 0xc9 0xc9 0x07

    // inData := []byte{104,101,108,108,111,32,102,114,111,109,32,98,114,111,119,115,101,114}

    // compressedData := new(bytes.Buffer) 

    // // compressedData := make([]byte, 100)
    // compress(inData, compressedData, 9)
    // // fmt.Println(byte(compressedData))
    // fmt.Println(compressedData.Bytes())

    // ioutil.WriteFile("compressed.dat", compressedData.Bytes(), os.ModeAppend)

// hello from browser
    inData1 := []byte{202,72,205,201,201,87,72,43,202,207,85,72,42,202,47,47,78,45,2,0}

    fmt.Println("originl comressed - ", inData1)

    a := decompress(&inData1) 
    fmt.Println(string(a))

    b := compress(&a, 9)
    fmt.Println(b)


    time.Sleep(10 * time.Second)
    // compressedData1 := new(bytes.Buffer) 
    // compressedData1.Write(inData1)

    // deCompressedData := new(bytes.Buffer)
    // decompress1(compressedData1, deCompressedData)
    // // log.Print(deCompressedData)
    // fmt.Println(string(deCompressedData.Bytes()))


    m := Mutatable{0, 0}
    fmt.Println(m)
    m.Mutate()
    fmt.Println(m)
    m.StayTheSame()
    fmt.Println(m)


}
// func compress(src []byte, dest io.Writer, level int) {
//     compressor, _ := flate.NewWriter(dest, level)
//     compressor.Write(src)
//     compressor.Close()
// }
// func decompress1(src io.Reader, dest io.Writer) {
//     decompressor := flate.NewReader(src)
//     io.Copy(dest, decompressor)
//     decompressor.Close()
// }




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