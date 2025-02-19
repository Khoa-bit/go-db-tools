package main

import (
	"go-db-tools/internal"
	"log"
)

func main() {
	log.Print("<><><><><><><><><><> BuildNestedModelExample <><><><><><><><><><>\n")
	internal.BuildNestedModelExample()
	log.Print("\n\n<><><><><><><><><><> BuildINQueryExample <><><><><><><><><><>\n")
	internal.BuildINQueryExample()
}
