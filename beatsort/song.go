package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

type Song struct {
	Name     string `json:"name"`
	Index    int    `json:"index"`
	Filename string `json:"filename"`
}

func (s *Song) init(fileName string, index int) {
	s.Filename = fileName
	s.Name = s.GetSongName(fileName)
	s.Index = index
}

func (s *Song) GetSongName(fileName string) string {
	index, _ := strconv.Atoi(fileName[0:2])
	fmt.Println(index)
	if index != 0 {
		return fileName[3:]
	}
	return fileName
}

func (s *Song) Rename(path string) error {
	// get the real path
	oldName := filepath.Join(path, s.Filename)
	newName := filepath.Join(path, s.GetSongName(s.Filename))

	// rename the song
	err := os.Rename(oldName, newName)

	if err != nil {
		msg := fmt.Sprintf("renaming song '%s': %v", s.Name, err)
		return errors.New(msg)
	}

	return nil
}
