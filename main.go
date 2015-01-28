package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	text := "В больничном дворе стоит небольшой флигель, окруженный целым лесом репейника, крапивы и дикой конопли."
	matches := Parse(text, "sentence", 0)
	str, _ := json.Marshal(matches)
	fmt.Printf("%s", str)
}
