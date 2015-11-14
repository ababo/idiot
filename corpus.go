package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

const forwardPunctuators = "«"
const backwardPunctuators = ",.:»"

func normalizeSentence(sentence string) string {
	sentence = strings.Trim(sentence, " ")
	if len(sentence) == 0 {
		return ""
	}

	for _, chr := range forwardPunctuators {
		punct := string(chr)
		sentence = strings.Replace(sentence, punct+" ", punct, -1)
	}

	for _, chr := range backwardPunctuators {
		punct := string(chr)
		sentence = strings.Replace(sentence, " "+punct, punct, -1)
	}

	last := sentence[len(sentence)-1]
	if last != '.' && last != '?' && last != '!' {
		return ""
	}

	return sentence
}

func BuildCorpus(xmlFilename, corpusFilename string) error {
	xml_, err := os.Open(xmlFilename)
	if err != nil {
		return err
	}
	defer xml_.Close()

	txt, err := os.Create(corpusFilename)
	if err != nil {
		return err
	}
	defer txt.Close()

	decoder := xml.NewDecoder(xml_)
	for {
		token, err := decoder.Token()
		if token == nil {
			break
		}
		if err != nil {
			return err
		}

		switch element := token.(type) {
		case xml.StartElement:
			if element.Name.Local != "source" {
				continue
			}

			var sentence string
			decoder.DecodeElement(&sentence, &element)

			if sentence = normalizeSentence(sentence); sentence == "" {
				continue
			}

			if _, err = txt.WriteString(sentence + "\n"); err != nil {
				return err
			}
		}
	}

	return nil
}

func ReadCorpus(corpusFilename string, from, to int) ([]string, error) {
	file, err := os.Open(corpusFilename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	index := 0
	var sentences []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if index >= from && index < to {
			sentences = append(sentences, scanner.Text())
		}
		index += 1
	}

	return sentences, scanner.Err()
}

type ParseCorpusSentenceCallback func(
	index int, sentence string, matches []ParseMatch, verbose bool)

type ParseCorpusStatsCallback func(succeeded, failed int)

func readRecord(file *os.File, index int) (bool, error) {
	byte_ := []byte{0}
	offset := int64(index / 8)

	if _, err := file.ReadAt(byte_, offset); err != nil {
		if err == io.EOF {
			return false, nil
		}
		return false, err
	}

	return byte_[0]&(1<<uint(index%8)) != 0, nil
}

func writeRecord(file *os.File, index int, parsed bool) error {
	byte_ := []byte{0}
	offset := int64(index / 8)

	if _, err := file.ReadAt(byte_, offset); err != nil && err != io.EOF {
		return err
	}

	if parsed {
		byte_[0] |= 1 << uint(index%8)
	} else {
		byte_[0] &= ^(1 << uint(index%8))
	}

	_, err := file.WriteAt(byte_, offset)
	return err
}

func ParseCorpus(corpusFilename string, from, to int,
	recordFilename string, saveChanges, verbose bool,
	sentenceCallback ParseCorpusSentenceCallback,
	statsCallback ParseCorpusStatsCallback) error {

	sentences, err := ReadCorpus(corpusFilename, from, to)
	if err != nil {
		return fmt.Errorf("failed to read text corpus: %s\n", err)
	}

	record, err := os.OpenFile(recordFilename, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("failed to open record file: %s\n", err)
	}
	defer record.Close()

	if err := initParser(); err != nil {
		return fmt.Errorf("failed to initialize parser: %s\n", err)
	}
	defer finalizeParser()

	lastTime := time.Now()
	succeeded, failed, index := 0, 0, from
	for _, sentence := range sentences {
		matches := Parse(strings.ToLower(sentence), "sentence", 0)
		ClearCache()

		parsed := false
		for _, m := range matches {
			if len(m.Text) == len(sentence) {
				parsed = true
				break
			}
		}

		parsedBefore, err := readRecord(record, index)
		if err != nil {
			return fmt.Errorf("failed to read record file: %s\n", err)
		}

		if parsed != parsedBefore {
			sentenceCallback(index, sentence, matches, verbose)
			if saveChanges {
				if err := writeRecord(record, index, parsed); err != nil {
					return fmt.Errorf("failed to write record file: %s\n", err)
				}
			}
		}

		if parsed {
			succeeded += 1
		} else {
			failed += 1
		}
		index += 1

		if now := time.Now(); now.Sub(lastTime).Seconds() >= 1 {
			statsCallback(succeeded, failed)
			lastTime = now
		}
	}

	statsCallback(succeeded, failed)

	return nil
}
