package main

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
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

func parse(from, to int) {
	dir := getRootDir()
	sentences, err := ReadCorpus(path.Join(dir, "corpus.txt"), from, to)
	if err != nil {
		fmt.Println("failed to read text corpus: %s\n", err)
		return
	}

	if err := initParser(); err != nil {
		fmt.Printf("failed to initialize parser: %s\n", err)
		return
	}
	defer finalizeParser()

	ok, failed := 0, 0
	lastTime := time.Now()
	for _, sentence := range sentences {
		matches := Parse(strings.ToLower(sentence), "sentence", 0)
		ClearCache()

		if len(matches) > 0 {
			ok += 1
		} else {
			failed += 1
		}

		if now := time.Now(); now.Sub(lastTime).Seconds() >= 1 {
			fmt.Printf("ok: %d, failed: %d\n", ok, failed)
			lastTime = now
		}
	}
	fmt.Printf("ok: %d, failed: %d\n", ok, failed)
}

func readIntArg(index, default_ int) int {
	if len(os.Args) <= index {
		return default_
	}

	if value, err := strconv.Atoi(os.Args[index]); err == nil {
		return value
	}

	return default_
}

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "morph":
			morph()
			return
		case "corpus":
			corpus()
			return
		case "parse":
			parse(readIntArg(2, 0), readIntArg(3, MaxInt))
			return
		}
	}
	fmt.Println("usage: idiot <command>\n" +
		"commands:\n" +
		"       morph           - build morphology base\n" +
		"       corpus          - build text corpus\n" +
		"       parse [from to] - parse text corpus\n")
}
