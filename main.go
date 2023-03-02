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
	playlistID := "5yeMjFdKUWapyQ7854FtBT"

	// initialize the beatsrot
	var beatSort BeatSort
	beatSort.init(ctx, playlistID)

	// find tracks tempo
	beatSort.ShowTracksTerminal()

}
