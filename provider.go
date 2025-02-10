package main

import (
	"log/slog"
	"os"
)

const ID = "cm3pbolfc0023h6k4d9t2j93o"

func test() {
	raw, error := os.Open("New York.csv")
	if err != nil {
		slog.Error(`os.Open("New York.csv")`, err)
		return
	}
}
