package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	text := "В больничном дворе стоит небольшой флигель, окруженный целым лесом репейника, крапивы и дикой конопли."
	//text := "стоит"
	matches := Parse(text, "sentence", 2)
	str, _ := json.MarshalIndent(matches, "", "  ")
	fmt.Printf("%s", str)
}
