package main

import (
	//"encoding/json"
	"fmt"
	"path"
	"runtime"
)

func getDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}

func main() {
	/*
		text := "В больничном дворе стоит небольшой флигель, окруженный целым лесом репейника, крапивы и дикой конопли."
		matches := Parse(text, "sentence", 0)
		str, _ := json.Marshal(matches)
		fmt.Printf("%s", str)
	*/

	dir := getDir()

	if err := InitData(path.Join(dir, "russian.db")); err != nil {
		fmt.Printf("failed to init data: %s\n", err)
		return
	}
	defer FinalizeData()
	/*
		skipped, err := buildTerminalData(path.Join(dir, "data.txt"))
		if err != nil {
			fmt.Printf("failed to build data: %s\n", err)
			return
		}

		fmt.Printf("Skipped values: %v\n", skipped)
	*/

	fmt.Printf("%v", FindTerminals("а", "."))
}
