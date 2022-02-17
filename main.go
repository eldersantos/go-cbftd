package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	ARRAY_LIMIT = 256
)

type ByteHistogram struct {
	Count    [ARRAY_LIMIT]uint64
	DataSize uint64
}

// NewByteHistogram creates a new ByteHistogram.
func NewByteHistogram() *ByteHistogram {
	return &ByteHistogram{}
}

func (bh *ByteHistogram) Init() {
	for i := 0; i < ARRAY_LIMIT; i++ {
		bh.Count[i] = 0
	}
	bh.DataSize = 0
}

// Update updates a ByteHistogram with an array of bytes.
func (bh *ByteHistogram) Update(bytes []byte) {
	for _, b := range bytes {
		if !isSpace(b) {
			bh.Count[b]++
		}
	}
	bh.DataSize += uint64(len(bytes))
}

// ByteList returns two values: a slice of the bytes that have been counted
// once at least and a slice with the actual number of times that every byte
// appears on the processed data.
func (bh *ByteHistogram) ByteList() ([]byte, []uint64) {
	bytelist := make([]byte, ARRAY_LIMIT)
	bytecount := make([]uint64, ARRAY_LIMIT)
	listlen := 0

	for i, c := range bh.Count {
		if c > 0 {
			bytelist[listlen] = byte(i)
			bytecount[listlen] = uint64(c)
			listlen++
		}
	}

	return bytelist[0:listlen], bytecount[0:listlen]
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\r' || c == '\n'
}

func main() {

	content, err := ioutil.ReadFile("./testdata/17.json")
	if err != nil {
		panic(err)
	}

	bh := NewByteHistogram()

	bh.Update(content)

	bytelist, bytecount := bh.ByteList()

	for i := range bytelist {
		fmt.Printf("%s - %d\n", string(bytelist[i]), bytecount[i])
	}

	fmt.Println(http.DetectContentType(content[:512]))
}
