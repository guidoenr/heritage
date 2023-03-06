package main

import (
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()
	playlistID := "7FfTSCdDior4RcVj02cXub"

	// initialize the beatsrot
	var beatSort BeatSort
	beatSort.init(ctx, playlistID)

	// find tracks BPM
	tracks, err := beatSort.GetTracks(true)
	if err != nil {
		fmt.Println(err)
	}
	for _, t := range tracks {
		fmt.Printf("[%dbpm] %s | Minutes: %.2f\n", t.BPM, t.Name, t.Minutes)
	}

}
