package main

import (
	"log"
	"os"
)

var logger *log.Logger

func main() {
	logger = log.New(os.Stdout, "", 0)
	text := "В больничном дворе стоит небольшой флигель, окруженный целым лесом репейника, крапивы и дикой конопли."
	matches := Parse(text, "sentence", 1)
	logger.Printf("%v", matches)
}
