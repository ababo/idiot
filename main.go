package main

import (
	"encoding/json"
	"fmt"
	"path"
	"runtime"
)

func getRootDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
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
		matches := Parse(text, "sentence", 0)

		str, _ := json.Marshal(matches)
		fmt.Printf("%s\n", str)
	} else {
		BuildMorph(path.Join(dir, "dict.opcorpora.txt"),
			path.Join(dir, "morph.bin"))
	}
}
