package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"path"
)

func morph() {
	fmt.Println("building morphology base...")

	dir := getRootDir()
	if err := BuildMorph(
		path.Join(dir, "dict.opcorpora.txt"),
		path.Join(dir, "morph.bin")); err != nil {
		fmt.Printf("failed to build morphology base: %s\n", err)
		return
	}

	fmt.Println("done")
}

func corpus() {
	fmt.Println("building text corpus...")

	dir := getRootDir()
	if err := BuildCorpus(
		path.Join(dir, "annot.opcorpora.xml"),
		path.Join(dir, "corpus.txt")); err != nil {
		fmt.Printf("failed to build text corpus: %s\n", err)
		return
	}

	fmt.Println("done")
}

func nextSentence(index int, sentence string,
	matches []ParseMatch, verbose bool) {
	if len(matches) > 0 {
		fmt.Printf("+%d: %s\n", index, sentence)
		if verbose {
			json, err := json.Marshal(matches)
			if err != nil {
				fmt.Printf("failed to marshal parse matches: %s\n", err)
				return
			}
			fmt.Printf("%s\n\n", json)
		}
	} else {
		fmt.Printf("-%d: %s\n", index, sentence)
	}
}

func nextStats(succeeded, failed int) {
	fmt.Printf("succeeded: %d, failed: %d\n", succeeded, failed)
}

func parse(from, to int, save, verbose bool) {
	dir := getRootDir()
	if err := ParseCorpus(path.Join(dir, "corpus.txt"),
		from, to, path.Join(dir, "corparse.rec"), save, verbose,
		nextSentence, nextStats); err != nil {
		fmt.Println(err)
	}
}

func main() {
	prefix := "(for \"parse\" command) "
	command := flag.String("command", "parse",
		"can be \"morph\", \"corpus\" or \"parse\"")
	from := flag.Int("from", 0, prefix+"begin of sentence interval")
	to := flag.Int("to", 1000000, prefix+"end of sentence interval")
	save := flag.Bool("save", false, prefix+"save result changes")
	verbose := flag.Bool("verbose", false, prefix+"verbose output")
	flag.Parse()

	switch *command {
	case "morph":
		morph()
		return
	case "corpus":
		corpus()
		return
	case "parse":
		parse(*from, *to, *save, *verbose)
		return
	default:
		flag.PrintDefaults()
	}
}
