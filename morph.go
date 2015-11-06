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

	morphGenderValues = []string{
		"masc", "femn", "neut",
	}

	morphTenseValues = []string{
		"past", "pres", "futr",
	}

	morphPersonValues = []string{
		"1per", "2per", "3per",
	}

	morphPersonMap = map[string]int{
		"impr": 1, /* 2per */
	}
)

type morphHeader struct {
	magic       uint32
	textSize    uint32
	entriesSize uint32
	reserved    uint32
}

type morphEntry struct {
	text  uint32
	attrs uint32
}

func (entry *morphEntry) getAttr(mask, offset uint32) uint32 {
	return (entry.attrs & mask) >> offset
}

func (entry *morphEntry) setAttr(imask, offset, value uint32) {
	entry.attrs = (entry.attrs & imask) | (value << offset)
}

func (entry *morphEntry) getPos() uint32 {
	return entry.getAttr(0x1F, 0)
}

func (entry *morphEntry) setPos(pos uint32) {
	entry.setAttr(^uint32(0x1F), 0, pos)
}

func (entry *morphEntry) getNumber() uint32 {
	return entry.getAttr(0x60, 5)
}

func (entry *morphEntry) setNumber(number uint32) {
	entry.setAttr(^uint32(0x60), 5, number)
}

func (entry *morphEntry) getCase() uint32 {
	return entry.getAttr(0x780, 7)
}

func (entry *morphEntry) setCase(case_ uint32) {
	entry.setAttr(^uint32(0x780), 7, case_)
}

func (entry *morphEntry) getGender() uint32 {
	return entry.getAttr(0x1800, 11)
}

func (entry *morphEntry) setGender(gender uint32) {
	entry.setAttr(^uint32(0x1800), 11, gender)
}

func (entry *morphEntry) getTense() uint32 {
	return entry.getAttr(0x6000, 13)
}

func (entry *morphEntry) setTense(tense uint32) {
	entry.setAttr(^uint32(0x6000), 13, tense)
}

func (entry *morphEntry) getPerson() uint32 {
	return entry.getAttr(0x18000, 15)
}

func (entry *morphEntry) setPerson(person uint32) {
	entry.setAttr(^uint32(0x18000), 15, person)
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
		entry.setPos(uint32(i + 1))
	} else if i := findString(morphNumberValues, value); i != -1 {
		entry.setNumber(uint32(i + 1))
	} else if i := findString(morphCaseValues, value); i != -1 {
		entry.setCase(uint32(i + 1))
	} else if i := findString(morphGenderValues, value); i != -1 {
		entry.setGender(uint32(i + 1))
	} else if i := findString(morphTenseValues, value); i != -1 {
		entry.setTense(uint32(i + 1))
	} else if i := findString(morphPersonValues, value); i != -1 {
		entry.setPerson(uint32(i + 1))
	} else if i, ok := morphPersonMap[value]; ok {
		entry.setPerson(uint32(i + 1))
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
		if strings.HasPrefix(l, "#") || len(split) < 2 {
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

func BuildMorph(txtFilename, morphFilename string) error {
	db, err := os.Create(morphFilename)
	if err != nil {
		return err
	}
	defer db.Close()

	content, err := ioutil.ReadFile(txtFilename)
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

func InitMorph(morphFilename string) error {
	FinalizeMorph()

	db, err := os.Open(morphFilename)
	if err != nil {
		return err
	}
	defer db.Close()

	var header morphHeader
	_, err = db.Read(castMorphHeaderToBytes(&header))
	if header.magic != morphMagic {
		return errors.New("bad file magic")
	}

	text := make([]byte, header.textSize)
	read, err := db.Read(text)
	if uint32(read) < header.textSize {
		return errors.New("unexpected end of file")
	}

	esize := uint32(unsafe.Sizeof(morphEntry{}))
	entries := make([]morphEntry, header.entriesSize/esize)
	read, err = db.Read(castMorphEntriesToBytes(entries))
	if uint32(read) < header.entriesSize {
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
	if pos := morphEntries[i].getPos(); pos > 0 {
		value := morphPosValues[pos-1]
		attrs = append(attrs, Attribute{Name: "pos", Value: value})
	}
	if number := morphEntries[i].getNumber(); number > 0 {
		value := morphNumberValues[number-1]
		attrs = append(attrs, Attribute{Name: "number", Value: value})
	}
	if case_ := morphEntries[i].getCase(); case_ > 0 {
		value := morphCaseValues[case_-1]
		attrs = append(attrs, Attribute{Name: "case", Value: value})
	}
	if gender := morphEntries[i].getGender(); gender > 0 {
		value := morphGenderValues[gender-1]
		attrs = append(attrs, Attribute{Name: "gender", Value: value})
	}
	if tense := morphEntries[i].getTense(); tense > 0 {
		value := morphTenseValues[tense-1]
		attrs = append(attrs, Attribute{Name: "tense", Value: value})
	}
	if person := morphEntries[i].getPerson(); person > 0 {
		value := morphPersonValues[person-1]
		attrs = append(attrs, Attribute{Name: "person", Value: value})
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
