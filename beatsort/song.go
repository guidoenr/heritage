package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

type Song struct {
	Index    int    `json:"index"`
	Filename string `json:"filename"`
	Name     string `json:"name"`
}

// init initializes the song with their correct format: "02 - Title - Artists"
func (s *Song) init(fileName string, index int) {
	s.Filename = fileName
	s.Name = s.AnalyzeSongName(fileName)

	if index < 10 {
		
	}
	s.Index = index
}

func (s *Song) AnalyzeSongName(fileName string) string {
	// get the index of the song and check if is a valid int
	index, _ := strconv.Atoi(fileName[0:2])
	if index != 0 {
		name := fileName[4:]
		return name
	}
	return fileName
}

func (s *Song) Rename(path string) error {
	// get the real path
	oldName := filepath.Join(path, s.Filename)
	fmt.Printf("old name: %s", oldName)

	// TODO handle logic of 01 and 1

	indexedName := fmt.Sprintf("%d - %s", s.Index, s.Name)
	newName := filepath.Join(path, indexedName)
	fmt.Printf("new name: %s", indexedName)

	// rename the song
	err := os.Rename(oldName, newName)

	if err != nil {
		msg := fmt.Sprintf("renaming song '%s': %v", s.Name, err)
		return errors.New(msg)
	}

	return nil
}
