package main

import (
	"encoding/json"
	"fmt"
	"path"
	"runtime"
)

func getDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}

func main() {
	dir := getDir()
	if true {
		err := InitMorph(path.Join(dir, "russian.morph"))
		if err != nil {
			fmt.Printf("%s\n", err)
		}
		defer FinalizeMorph()

		err = InitRules(path.Join(dir, "russian.rules"))
		if err != nil {
			fmt.Printf("%s\n", err)
		}
		defer FinalizeRules()

		text := "в больничном дворе стоит небольшой флигель, окружённый целым лесом репейника, крапивы и дикой конопли."
		var matches []ParseMatch
		for i := 0; i < 100; i++ {
			matches = Parse(text, "sentence", 1)
		}
		str, _ := json.Marshal(matches)
		fmt.Printf("%s\n", str)
	} else {
		BuildMorph(path.Join(dir, "dict.opcorpora.txt"),
			path.Join(dir, "russian.morph"))
	}
}
