package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/zmb3/spotify"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
	"log"
	"os"
	"sort"
	"sync"
)

// BeatSort is the class to encapsulate all the spotify logic
type BeatSort struct {
	client     spotify.Client // the spotify client
	playlistID string
}

// init initializes the BeatSort class
func (bs *BeatSort) init(ctx context.Context, playlistID string) {
	// read the local env configs
	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     spotifyauth.TokenURL,
	}
	// create the token
	token, err := config.Token(ctx)
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}
	// create the http client
	httpClient := spotifyauth.New().Client(ctx, token)

	// create and set the fields
	bs.client = spotify.NewClient(httpClient)
	bs.playlistID = playlistID
}

// ShowTracksTerminal prints the tracks in terminal based
func (bs *BeatSort) ShowTracksTerminal() {
	tracks, _ := bs.GetTrackAudioFeatures()

	for _, t := range tracks {
		fmt.Printf("[%d bpm] %s - %s \n", int(t.BeatsPerMin), t.TrackName, t.ArtistName)
	}

}

// GetTrackAudioFeatures has the responsibility of launch the algorithm with go routines
// to find all the BPM in all the tracks
func (bs *BeatSort) GetTrackAudioFeatures() ([]TrackFeatures, error) {
	// get the spotifyTracks
	spotifyTracks, err := bs.getTracks()
	if err != nil {
		return nil, err
	}

	// create a channel to collect the audio features for each track
	trackFeaturesChan := make(chan TrackFeatures)

	// create a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	// for each track, launch a new goroutine to fetch the audio features
	for _, t := range spotifyTracks.Tracks {
		wg.Add(1)
		// ---------------------------------- GO ROUTINE
		go func(track spotify.FullTrack) {
			// TODO, output in a file
			/*log.Printf("[goroutine %s] with track '%s' \n", track.ID, track.Name)*/
			defer wg.Done()
			// get the audio features of the track
			audioFeatures, err := bs.client.GetAudioFeatures(track.ID)
			if err != nil {
				log.Printf("error fetching audio features for %s: %v", track.Name, err)
				return
			}

			// add the content of the audioFeatures to the features channel
			trackFeaturesChan <- TrackFeatures{
				TrackName:   track.Name,
				ArtistName:  track.Artists[0].Name,
				BeatsPerMin: audioFeatures[0].Tempo,
				Genre:       audioFeatures[0].Energy,
			}
			/*log.Printf("[goroutine %s] finished!\n", track.ID)*/

		}(t.Track)
		// ---------------------------------- GO ROUTINE
	}

	// wait for all goroutines to finish and close the channel
	go func() {
		wg.Wait()
		close(trackFeaturesChan)
	}()

	// read data from the channel and append it to the slice
	var tracks []TrackFeatures
	for track := range trackFeaturesChan {
		tracks = append(tracks, track)
	}

	// sort the slice based on the BeatsPerMin field
	// (quicksort)
	sort.Slice(tracks, func(i, j int) bool {
		return tracks[i].BeatsPerMin < tracks[j].BeatsPerMin
	})

	return tracks, nil
}

// getTracks returns the entire playlist tracks page
func (bs *BeatSort) getTracks() (*spotify.PlaylistTrackPage, error) {
	// getting the tracks
	tracks, err := bs.client.GetPlaylistTracks(spotify.ID(bs.playlistID))
	if err != nil {
		msg := fmt.Sprintf("getting playlist tracks: %v", err)
		return nil, errors.New(msg)
	}

	log.Printf("total tracks found in playlist: %d \n", tracks.Total)

	return tracks, nil
}
