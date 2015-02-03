package main

import (
	"bufio"
	"errors"
	"io/ioutil"
	"os"
	"reflect"
	"sort"
	"strings"
	"unsafe"
)

const (
	morphMagic      = 0xFC1290C8
	morphBufferSize = 0x80000
)

var (
	morphPosValues = []string{
		"noun", "advb", "adjf", "adjs", "comp", "verb", "infn", "prtf",
		"prts", "grnd", "conj", "intj", "prcl", "prep", "pred", "numr",
		"npro",
	}

	morphNumberValues = []string{
		"sing", "plur",
	}

	morphCaseValues = []string{
		"nomn", "gent", "gen1", "gen2", "datv", "accs", "ablt", "loct",
		"loc1", "loc2", "voct",
	}
)

type morphHeader struct {
	magic        uint32
	text_size    uint32
	entries_size uint32
	reserved     uint32
}

type morphEntry struct {
	text     uint32
	pos      uint8
	number   uint8
	case_    uint8
	reserved uint8
}

func findString(strs []string, str string) int {
	for i, s := range strs {
		if s == str {
			return i
		}
	}
	return -1
}

func findMorphEntry(entries []morphEntry, entry morphEntry) int {
	for i, e := range entries {
		if e == entry {
			return i
		}
	}
	return -1
}

func updateMorphEntry(entry *morphEntry, value string) {
	if i := findString(morphPosValues, value); i != -1 {
		entry.pos = uint8(i + 1)
	} else if i := findString(morphNumberValues, value); i != -1 {
		entry.number = uint8(i + 1)
	} else if i := findString(morphCaseValues, value); i != -1 {
		entry.case_ = uint8(i + 1)
	}
}

func castMorphEntriesToBytes(entries []morphEntry) []byte {
	var bytes []byte
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&entries))
	sh2 := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	sh2.Data = sh.Data
	sh2.Len = sh.Len * int(unsafe.Sizeof(morphEntry{}))
	sh2.Cap = sh2.Len
	return bytes
}

func parseMorphLines(
	lines []string, writer *bufio.Writer) (uint32, uint32, error) {

	var prevText string
	var prevEntries []morphEntry
	var prevTsize, tsize, esize uint32
	entries := make([]morphEntry, 0, len(lines))

	for _, l := range lines {
		split := strings.Split(l, "\t")
		if len(split) < 2 {
			continue
		}

		text := split[0] + "\x00"
		if text != prevText {
			writer.WriteString(text)
			prevTsize = tsize
			tsize += uint32(len(text))
			prevText = text
			prevEntries = nil
		}

		entry := morphEntry{text: prevTsize}
		split = strings.FieldsFunc(split[1], func(r rune) bool {
			return r == ',' || r == ' '

		})
		for _, v := range split {
			updateMorphEntry(&entry, v)
		}

		if findMorphEntry(prevEntries, entry) == -1 {
			entries = append(entries, entry)
			esize += uint32(unsafe.Sizeof(entry))
			prevEntries = append(prevEntries, entry)
		}
	}

	_, err := writer.Write(castMorphEntriesToBytes(entries))
	return tsize, esize, err
}

func castMorphHeaderToBytes(header *morphHeader) []byte {
	var bytes []byte
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	sh.Data = uintptr(unsafe.Pointer(header))
	sh.Len = int(unsafe.Sizeof(*header))
	sh.Cap = sh.Len
	return bytes
}

func BuildMorph(txt_filename, morph_filename string) error {
	db, err := os.Create(morph_filename)
	if err != nil {
		return err
	}
	defer db.Close()

	content, err := ioutil.ReadFile(txt_filename)
	if err != nil {
		return err
	}

	lines := strings.Split(strings.ToLower(string(content)), "\n")
	sort.Strings(lines)

	var header morphHeader
	db.Seek(int64(unsafe.Sizeof(header)), 0)
	writer := bufio.NewWriterSize(db, morphBufferSize)
	tsize, esize, err := parseMorphLines(lines, writer)
	if err != nil {
		return err
	}
	if err := writer.Flush(); err != nil {
		return err
	}

	header = morphHeader{morphMagic, tsize, esize, 0}
	_, err = db.WriteAt(castMorphHeaderToBytes(&header), 0)
	return err
}

var (
	morphText    string
	morphEntries []morphEntry
)

func InitMorph(morph_filename string) error {
	FinalizeMorph()

	db, err := os.Open(morph_filename)
	if err != nil {
		return err
	}
	defer db.Close()

	var header morphHeader
	_, err = db.Read(castMorphHeaderToBytes(&header))
	if header.magic != morphMagic {
		return errors.New("bad file magic")
	}

	text := make([]byte, header.text_size)
	read, err := db.Read(text)
	if uint32(read) < header.text_size {
		return errors.New("unexpected end of file")
	}

	esize := uint32(unsafe.Sizeof(morphEntry{}))
	entries := make([]morphEntry, header.entries_size/esize)
	read, err = db.Read(castMorphEntriesToBytes(entries))
	if uint32(read) < header.entries_size {
		return errors.New("unexpected end of file")
	}

	morphText = string(text)
	morphEntries = entries
	return nil
}

func FinalizeMorph() {
	morphText = ""
	morphEntries = nil
}

func getMorphEntryText(i int) string {
	from := morphEntries[i].text
	len := strings.Index(morphText[from:], "\x00")
	return morphText[from : int(from)+len]
}

func getMorphEntryMatch(i int) ParseMatch {
	attrs := []Attribute{}
	if morphEntries[i].pos > 0 {
		value := morphPosValues[morphEntries[i].pos-1]
		attrs = append(attrs, Attribute{Name: "pos", Value: value})
	}
	if morphEntries[i].number > 0 {
		value := morphNumberValues[morphEntries[i].number-1]
		attrs = append(attrs, Attribute{Name: "number", Value: value})
	}
	if morphEntries[i].case_ > 0 {
		value := morphCaseValues[morphEntries[i].case_-1]
		attrs = append(attrs, Attribute{Name: "case", Value: value})
	}

	return ParseMatch{Text: getMorphEntryText(i), Attributes: attrs}
}

func FindTerminals(prefix, separator string) []ParseMatch {
	index := sort.Search(len(morphEntries), func(i int) bool {
		return getMorphEntryText(i) >= prefix
	})

	matches := []ParseMatch{}
	for i := index; i < len(morphEntries) &&
		getMorphEntryText(i) == prefix; i++ {
		matches = append(matches, getMorphEntryMatch(i))
	}

	if len(separator) == 0 {
		return matches
	}

	prefix += separator
	index += sort.Search(len(morphEntries)-index-1, func(i int) bool {
		return getMorphEntryText(index+i+1) >= prefix
	}) + 1

	for i := index; i < len(morphEntries) &&
		strings.HasPrefix(getMorphEntryText(i), prefix); i++ {
		matches = append(matches, getMorphEntryMatch(i))
	}

	return matches
}
