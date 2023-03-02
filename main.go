package main

import (
	"context"
)

type TrackFeatures struct {
	TrackName   string
	ArtistName  string
	BeatsPerMin float32
	Genre       float32
}

func main() {
	ctx := context.Background()
	playlistID := "7FfTSCdDior4RcVj02cXub"

	// initialize the beatsrot
	var beatSort BeatSort
	beatSort.init(ctx, playlistID)

	// find tracks tempo
	beatSort.ShowTracksTerminal()

}
