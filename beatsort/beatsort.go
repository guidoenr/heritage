package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Song struct {
	NewName string `json:"newname"`
	OldName string `json:"oldname"`
	Index   int    `json:"index"`
}

type BeatSort struct {
	Playlist []Song
}

// init initializes the BeatSort module
func (bs *BeatSort) init(path string) error {
	// get all the files inside the dir
	filesDir, err := os.ReadDir(path)
	if err != nil {
		fmt.Println(err)
	}

	// fill the playlist
	for index, file := range filesDir {
		bs.Playlist = append(bs.Playlist, Song{
			OldName: file.Name(),
			NewName: fmt.Sprintf("[%d]%s", index, file.Name()),
			Index:   index,
		})
	}

	return nil
}

// SavePlaylist just write the content in the order.json file
func (bs *BeatSort) SavePlaylist() error {
	// get the json data
	data, err := json.MarshalIndent(bs.Playlist, "", "  ")
	if err != nil {
		return err
	}

	// write the file
	err = os.WriteFile("order.json", data, 0644)
	if err != nil {
		return err
	}

	return nil
}

// LoadPlaylist just write the content in the order.json file
func (bs *BeatSort) LoadPlaylist() error {
	data, err := os.ReadFile("order.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &bs.Playlist)
	if err != nil {
		return err
	}
	return nil
}
