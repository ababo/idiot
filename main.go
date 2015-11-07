package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
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

func test(nonterminal, text string) {
	if err := initParser(); err != nil {
		fmt.Printf("failed to initialize parser: %s\n", err)
	}
	defer finalizeParser()

	matches := Parse(nonterminal, text, 0)
	json, err := json.Marshal(matches)
	if err != nil {
		fmt.Printf("failed to marshal parse matches: %s\n", err)
		return
	}
	fmt.Printf("%s\n", json)
}

func printUsage() {
	fmt.Printf("usage: idiot (morph | corpus | parse | test)\n")
}

func main() {
	if len(os.Args) == 1 {
		printUsage()
		return
	}

	switch os.Args[1] {
	case "morph":
		morph()
		return
	case "corpus":
		corpus()
		return
	case "parse":
		parseFlags := flag.NewFlagSet("parse", flag.ExitOnError)
		from := parseFlags.Int("from", 0, "begin of sentence interval")
		to := parseFlags.Int("to", 1000000, "end of sentence interval")
		save := parseFlags.Bool("save", false, "save result changes")
		verbose := parseFlags.Bool("verbose", false, "verbose output")
		parseFlags.Parse(os.Args[2:])
		parse(*from, *to, *save, *verbose)
		return
	case "test":
		testFlags := flag.NewFlagSet("test", flag.ExitOnError)
		nonterm := testFlags.String("nonterm",
			"sentence", "non-terminal to parse against")
		text := testFlags.String("text", "", "text to parse")
		testFlags.Parse(os.Args[2:])
		test(*nonterm, *text)
	default:
		printUsage()
	}
}
