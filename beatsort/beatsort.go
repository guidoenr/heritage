package main

import (
	"fmt"
	"os"
	"strconv"
)

type Song struct {
	Name     string `json:"name"`
	Index    int    `json:"int"`
	Filename string `json:"filename"`
}

func (s *Song) init(name string, index int) {
	s.Name = name
	s.Index = index + 1

	// first get the old index of the song
	oldIndex, _ := strconv.Atoi(name[0:2])
	if oldIndex != 0 {
		s.RemoveIndex()
	}

	s.Filename = fmt.Sprintf("%d- %s", s.Index, s.Name)
}

func (s *Song) RemoveIndex() {
	s.Name = s.Name[3:]
}

func (s *Song) SetIndex(index int) {
	// set the new index
	s.Index = index
	s.Filename = string(rune(index)) + "-" + s.Name
}

type BeatSort struct {
	Path     string
	Playlist []Song
}

func (bs *BeatSort) init(path string) {
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

}
