package main

import (
	"encoding/xml"
	"os"
	"strings"
)

const forward_punctuators = "«"
const backward_punctuators = ",.:»"

func normalizeSentence(sentence string) string {
	sentence = strings.Trim(sentence, " ")
	if len(sentence) == 0 {
		return ""
	}

	for _, chr := range forward_punctuators {
		punct := string(chr)
		sentence = strings.Replace(sentence, punct+" ", punct, -1)
	}

	for _, chr := range backward_punctuators {
		punct := string(chr)
		sentence = strings.Replace(sentence, " "+punct, punct, -1)
	}

	last := sentence[len(sentence)-1]
	if last != '.' && last != '?' && last != '!' {
		return ""
	}

	return sentence
}

func BuildCorpus(xml_filename, corpus_filename string) error {
	xml_, err := os.Open(xml_filename)
	if err != nil {
		return err
	}
	defer xml_.Close()

	txt, err := os.Create(corpus_filename)
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
