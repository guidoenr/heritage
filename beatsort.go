package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type BeatSort struct {
	songs []string // the songs list [url] to download
}

// init initializes the beatsort
func (bs *BeatSort) init(songs []string) {
	bs.songs = songs
}

func (bs *BeatSort) getTempo(song string) {
	url := "https://developer.echonest.com/api/v4/track/profile?api_key=YOUR_API_KEY&format=json&id=spotify:track:0eGsygTp906u18L0Oimnem&bucket=audio_summary"
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error retrieving tempo:", err)
		return
	}
	defer response.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(response.Body).Decode(&result)
	tempo := result["response"].(map[string]interface{})["track"].(map[string]interface{})["tempo"].(float64)
	fmt.Println("Tempo:", tempo)
}
