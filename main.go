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

	// find tracks tempo
	tracks, err := beatSort.GetTracksAudioFeatures()
	if err != nil {
		fmt.Println(err)
	}
	for _, t := range tracks {
		fmt.Printf("[%dbpm] %s | duration: %.2f \n", t.tempo, t.name, t.duration)
	}

}
