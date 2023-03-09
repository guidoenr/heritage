package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type BeatSort struct {
	Path     string
	Playlist []Song
}

func (bs *BeatSort) init(path string) error {
	// get all the files inside the dir
	filesDir, err := os.ReadDir(path)
	if err != nil {
		fmt.Println(err)
	}

	for index, file := range filesDir {
		var newSong Song
		newSong.init(file.Name(), index)

		bs.Playlist = append(bs.Playlist, newSong)
	}

	err = bs.SavePlaylistState()
	if err != nil {
		return err
	}

	return nil
}

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
