package main

import "fmt"

func main() {

	// initialize the beatsrot
	var beatSort BeatSort
	err := beatSort.init("/home/guido/Desktop/heritage/beatsort/playlist/")
	fmt.Println(err)
	err = beatSort.SortPlaylist()
	fmt.Println(err)

	err = beatSort.SavePlaylistState()
	err = beatSort.SortPlaylist()
}
