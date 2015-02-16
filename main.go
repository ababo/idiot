package main

import (
	"encoding/json"
	"fmt"
	"path"
	"runtime"
	"sort"
	"strings"
)

func getRootDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}

func getParseMatchesJson(matches []ParseMatch) string {
	jsons := make([]string, len(matches))
	for i, m := range matches {
		data, _ := json.Marshal(&m)
		jsons[i] = string(data)
	}

	sort.Strings(jsons)
	return "[" + strings.Join(jsons, ",") + "]"
}

func main() {
	dir := getRootDir()
	if true {
		err := InitMorph(path.Join(dir, "morph.bin"))
		if err != nil {
			fmt.Printf("%s\n", err)
		}
		defer FinalizeMorph()

		err = InitRules(path.Join(dir, "rules.yaml"))
		if err != nil {
			fmt.Printf("%s\n", err)
		}
		defer FinalizeRules()

		InitCache(256)
		defer FinalizeCache()

		//text := "крыша на нем ржавая, труба наполовину обвалилась, ступеньки у крыльца сгнили и поросли  травой, а от штукатурки остались одни только следы."
		text := "в больничном дворе стоит небольшой флигель, окружённый целым лесом репейника, крапивы и дикой конопли."

		fmt.Print(getParseMatchesJson(Parse(text, "sentence", 0)))
	} else {
		BuildMorph(path.Join(dir, "dict.opcorpora.txt"),
			path.Join(dir, "morph.bin"))
	}
}
