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
		if err := InitMorph(path.Join(dir, "russian.mdb")); err != nil {
			fmt.Printf("error calling InitMorph: %s", err)
		}
		defer FinalizeMorph()

		text := "в больничном дворе стоит небольшой флигель, окружённый целым лесом репейника, крапивы и дикой конопли."
		var matches []ParseMatch
		for i := 0; i < 100; i++ {
			matches = Parse(text, "sentence", 1)
		}
		str, _ := json.Marshal(matches)
		fmt.Printf("%s", str)
	} else {
		BuildMorphDb(path.Join(dir, "dict.opcorpora.txt"), path.Join(dir, "russian.mdb"))
	}
}
