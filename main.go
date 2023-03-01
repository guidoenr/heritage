package main

import (
	"context"
	"fmt"
	"github.com/zmb3/spotify"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
	"log"
	"math"
	"os"
	"sync"
)

type TrackFeatures struct {
	TrackName   string
	ArtistName  string
	BeatsPerMin float32
	Genre       float32
}

func main() {
	ctx := context.Background()
	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     spotifyauth.TokenURL,
	}
	token, err := config.Token(ctx)
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	httpClient := spotifyauth.New().Client(ctx, token)
	client := spotify.NewClient(httpClient)

	// Public playlist owned by noah.stride:
	// "Long playlist for testing pagination"
	/*playlistID := "be88a05423a943e9"
	if id := os.Getenv("SPOTIFY_PLAYLIST"); id != "" {
		playlistID = id
	}
	*/
	playlistID := "7FfTSCdDior4RcVj02cXub"

	// getting the tracks
	tracks, err := client.GetPlaylistTracks(spotify.ID(playlistID))
	if err != nil {
		log.Fatal(err)
	}

	// create a channel to collect the audio features for each track
	trackFeaturesChan := make(chan TrackFeatures)

	// create a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	// for each track, launch a new goroutine to fetch the audio features
	for _, t := range tracks.Tracks {
		wg.Add(1)
		go func(track spotify.FullTrack) {
			defer wg.Done()
			audioFeatures, err := client.GetAudioFeatures(track.ID)
			if err != nil {
				log.Printf("error fetching audio features for %s: %v", track.Name, err)
				return
			}
			trackFeaturesChan <- TrackFeatures{
				TrackName:   track.Name,
				ArtistName:  track.Artists[0].Name,
				BeatsPerMin: audioFeatures[0].Tempo,
				Genre:       audioFeatures[0].Energy,
			}
		}(t.Track)
	}

	// wait for all goroutines to finish and close the channel
	go func() {
		wg.Wait()
		close(trackFeaturesChan)
	}()

	// now print the entire list
	for track := range trackFeaturesChan {
		fmt.Printf("[%fBPM]- %s - %s  | %s \n", math.Round(float64(track.BeatsPerMin)), track.TrackName, track.ArtistName, track.Genre)
	}

}
