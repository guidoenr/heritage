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

type HTrack struct {
	Id      spotify.ID
	Name    string
	BPM     int
	Minutes float64
	URI     string
}

// GetTracks returns all tracks in the playlist handling pagination if necessary
func (bs *BeatSort) GetTracks(sortedBPM bool) ([]HTrack, error) {
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

	// now fill the htracks slice
	var htracks []HTrack
	for _, t := range tracks {
		trackName := t.Track.Name + " - "
		for _, artist := range t.Track.Artists {
			trackName += artist.Name + " - "
		}

		htracks = append(htracks, HTrack{
			Id:      t.Track.ID,
			Name:    trackName[0:],
			BPM:     0,
			Minutes: float64(t.Track.Duration) / 60000,
			URI:     string(t.Track.URI),
		})
	}

	if sortedBPM {
		// set the track Id's
		var tracksId []spotify.ID
		for _, track := range htracks {
			tracksId = append(tracksId, track.Id)
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
			return nil, errors.New("NOT IMPLEMENTED YET: playlist has more than 100 htracks")
		}

		// get all the audio features
		audioFeatures, err := bs.client.GetAudioFeatures(tracksId...)
		if err != nil {
			msg := fmt.Sprintf("error fetching audio features: %v", err)
			return nil, errors.New(msg)
		}

		for i, af := range audioFeatures {
			htracks[i].BPM = int(af.Tempo)
		}

		// sort the slice based on the BeatsPerMin field
		// (quicksort)
		sort.Slice(htracks, func(i, j int) bool {
			return htracks[i].BPM < htracks[j].BPM
		})

		return htracks, nil
	}

	log.Printf("total tracks found in playlist: %d\n", len(htracks))
	return htracks, nil
}
