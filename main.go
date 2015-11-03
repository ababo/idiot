package main

import (
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

func parse() {
	if err := initParser(); err != nil {
		fmt.Printf("failed to initialize parser: %s\n", err)
		return
	}
	defer finalizeParser()

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
			parse()
			return
		}
	}
	fmt.Println("usage: idiot <command>\n" +
		"commands:\n" +
		"       morph  - build morphology base\n" +
		"       corpus - build text corpus\n" +
		"       parse  - parse text corpus\n")
}
