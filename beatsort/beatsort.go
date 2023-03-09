package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type BeatSort struct {
	PlaylistPath string
	Playlist     []Song
}

func (bs *BeatSort) init(path string) error {
	// set the path
	bs.PlaylistPath = path

	// get all the files inside the dir
	filesDir, err := os.ReadDir(bs.PlaylistPath)
	if err != nil {
		return err
	}

	// for each file in the directory
	for index, file := range filesDir {
		var newSong Song
		fmt.Println(file.Name())
		newSong.init(file.Name(), index+1) //  because the index in the list starts from 0

		bs.Playlist = append(bs.Playlist, newSong)
	}

	err = bs.SavePlaylistState()
	if err != nil {
		return err
	}

	return nil
}

// SavePlaylistState just dumps the content in the bs.Playlist and dump it into order.json
func (bs *BeatSort) SavePlaylistState() error {
	// get the json data
	data, err := json.MarshalIndent(bs.Playlist, "", "   ")
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

// SortPlaylist changes the name of each file in the directory with the given order in the `order.json` file
func (bs *BeatSort) SortPlaylist() error {
	// read the file
	data, err := os.ReadFile("order.json")
	if err != nil {
		return err
	}

	// put the data in the playlist
	err = json.Unmarshal(data, &bs.Playlist)
	if err != nil {
		return err
	}

	// for each song in the playlist
	for _, song := range bs.Playlist {
		err = song.Rename(bs.PlaylistPath)
		if err != nil {
			return err
		}
	}

	return nil
}
