package main

import (
	"fmt"
	"league-of-legends-fight-tactics/internal/champion_reader"
)

func main() {
	championReader := champion_reader.NewReader()

	champion, err := championReader.LoadChampion("Jhin")
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	fmt.Printf("%+v\n", champion)
}