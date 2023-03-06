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
	tracks, err := bs.GetTracksAudioFeatures()
	if err != nil {
		fmt.Println("nothing")
	} else {
		for _, t := range tracks {
			fmt.Println(t.Tempo)
		}
	}
}

// GetTracksAudioFeatures returns the audio features of each track
func (bs *BeatSort) GetTracksAudioFeatures() ([]*spotify.AudioFeatures, error) {
	// get the totalTracks
	totalTracks, err := bs.GetTracks()
	if err != nil {
		return nil, err
	}

	// set the track id's
	var tracksId []spotify.ID
	for _, track := range totalTracks {
		tracksId = append(tracksId, track.Track.ID)
	}

	// 170 ? ==> 100 + 70
	// the request limit for getAudioFeatures
	//limit := 100
	//totalRoutines := 1
	//
	//// if the playlist contains more than 100 songs to analyze
	//if len(tracksId) > 100 {
	//	lengthTracksId := len(tracksId)
	//	// while the length is greater than 100
	//	for lengthTracksId > 100 {
	//		// reduce the length (for example: 601 songs - 100 limit = 501 songs)
	//		lengthTracksId = lengthTracksId - limit
	//		totalRoutines += 1
	//	}
	//}
	//log.Printf("total routines to launch: %d", totalRoutines)

	if len(tracksId) > 100 {
		return nil, errors.New("NOT IMPLEMENTED YET: playlist has more than 100 tracks")
	}

	// get all the audio features
	audioFeatures, err := bs.client.GetAudioFeatures(tracksId...)
	if err != nil {
		msg := fmt.Sprintf("error fetching audio features: %v", err)
		return nil, errors.New(msg)
	}

	// sort the slice based on the BeatsPerMin field
	// (quicksort)
	sort.Slice(audioFeatures, func(i, j int) bool {
		return audioFeatures[i].Tempo < audioFeatures[j].Tempo
	})

	return audioFeatures, nil
}

// GetTracks returns all tracks in the playlist, handling pagination if necessary
func (bs *BeatSort) GetTracks() ([]spotify.PlaylistTrack, error) {
	log.Println("getting tracks..")
	// create the slice of playlistTrack
	var tracks []spotify.PlaylistTrack

	// setting the options
	limit := 100
	opt := &spotify.Options{
		Limit: &limit,
	}

	// [pagination]
	for {
		// get the page
		page, err := bs.client.GetPlaylistTracksOpt(spotify.ID(bs.playlistID), opt, "")
		if err != nil {
			msg := fmt.Sprintf("getting playlist tracks: %v", err)
			return nil, errors.New(msg)
		}

		// append the tracks in that page
		tracks = append(tracks, page.Tracks...)
		if page.Next == "" {
			break
		}

		// creating the new offset
		newOffset := page.Offset + page.Limit
		opt.Offset = &newOffset
	}

	log.Printf("total tracks found in playlist: %d\n", len(tracks))
	return tracks, nil
}
