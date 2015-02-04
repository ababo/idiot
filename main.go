package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"path"
	"runtime"
)

func getRootDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}

func main() {
	dir := getRootDir()
	fmt.Printf("dir: %s\n", dir)
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

		text := "крыша на нем ржавая, труба наполовину обвалилась, ступеньки у крыльца сгнили и поросли  травой, а от штукатурки остались одни только следы."

		matches := Parse(text, "sentence", 0)

		str, _ := yaml.Marshal(matches)
		fmt.Printf("%s\n", str)
	} else {
		BuildMorph(path.Join(dir, "dict.opcorpora.txt"),
			path.Join(dir, "russian.morph"))
	}
}
