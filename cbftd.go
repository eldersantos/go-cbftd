package cbftd

import (
	"fmt"
	"io/ioutil"
	"log"
	"sort"
)

const (
	ARRAY_LIMIT = 256 // assuming english files 256 is enough
)

type ByteHistogram struct {
	Count [ARRAY_LIMIT]uint64
}

type byteCountPair struct {
	b byte
	c uint64
}

type byCountAsc []byteCountPair

func (a byCountAsc) Len() int      { return len(a) }
func (a byCountAsc) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byCountAsc) Less(i, j int) bool {
	return (a[i].c < a[j].c) || (a[i].c == a[j].c && a[i].b < a[j].b)
}

type byCountDesc []byteCountPair

func (a byCountDesc) Len() int      { return len(a) }
func (a byCountDesc) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byCountDesc) Less(i, j int) bool {
	return (a[i].c > a[j].c) || (a[i].c == a[j].c && a[i].b < a[j].b)
}

// NewByteHistogram creates a new ByteHistogram.
func NewByteHistogram() *ByteHistogram {
	return &ByteHistogram{}
}

// Update updates a ByteHistogram with an array of bytes.
func (bh *ByteHistogram) Update(bytes []byte) {
	for _, b := range bytes {
		if !isFmtChar(b) {
			bh.Count[b]++
		}
	}
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

func isFmtChar(c byte) bool {
	return c == ' ' || c == '\t' || c == '\r' || c == '\n'
}

// SortedByteList returns two values as the ByteList function does, but the
// resulting slices are sorted by the number of bytes.
// The sorting order is specified by ascOrder, that will be ascending if
// the param is true or descending if it is false.
func (bh *ByteHistogram) SortedByteList(ascOrder bool) ([]byte, []uint64) {
	pairs := make([]byteCountPair, 256)

	for i, count := range bh.Count {
		pairs[i] = byteCountPair{b: byte(i), c: count}
	}

	if ascOrder {
		sort.Sort(byCountAsc(pairs))
	} else {
		sort.Sort(byCountDesc(pairs))
	}

	bytelist := make([]byte, 256)
	bytecount := make([]uint64, 256)
	listlen := 0

	for _, pair := range pairs {
		if pair.c > 0 {
			bytelist[listlen] = pair.b
			bytecount[listlen] = pair.c
			listlen++
		}
	}

	return bytelist[0:listlen], bytecount[0:listlen]
}

func (bh *ByteHistogram) Train(samples string) {

	files, err := ioutil.ReadDir(samples)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if !file.IsDir() {
			content, err := ioutil.ReadFile(samples + file.Name())
			if err != nil {
				panic(err)
			}
			bh.Update(content)
		}
	}

	bh.norm()
}

func (bh *ByteHistogram) String() (s string) {
	bytelist, bytecount := bh.SortedByteList(false)
	for i := range bytelist {
		s += fmt.Sprintf("%s - %d\n", string(bytelist[i]), bytecount[i])
	}

	return fmt.Sprintf("\n%s", s)
}

// normalize the values
// assuming the slice is sorted already
func (bh *ByteHistogram) norm() {
	top := bh.Count[0]
	for i := range bh.Count {
		if bh.Count[i] != 0 {
			bh.Count[i] = bh.Count[i] / top
		}
	}
}
