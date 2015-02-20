package main

import (
	"encoding/json"
	"fmt"
	"path"
)

func main() {
	if true {
		if err := initParser(); err != nil {
			fmt.Println(err)
		}
		defer finalizeParser()

		text := "крыша на нём ржавая, труба наполовину обвалилась, ступеньки у крыльца сгнили и поросли травой, а от штукатурки остались одни только следы."
		matches := Parse(text, "sentence", 0)
		data, _ := json.Marshal(&matches)
		fmt.Println(string(data))
	} else {
		dir := getRootDir()

		if err := BuildMorph(
			path.Join(dir, "dict.opcorpora.txt"),
			path.Join(dir, "morph.bin")); err != nil {
			fmt.Println(err)
		}
	}
}
