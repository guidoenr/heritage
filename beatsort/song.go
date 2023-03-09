package main

import (
	"fmt"
	"strconv"
)

type Song struct {
	Name     string `json:"name"`
	Index    int    `json:"index"`
	Filename string `json:"filename"`
}

func (s *Song) init(name string, index int) {
	s.Name = name
	s.Index = index + 1

	// first get the old index of the song
	oldIndex, _ := strconv.Atoi(name[0:2])
	newName := s.Name
	if oldIndex != 0 {
		newName = s.RemoveIndex()
	}

	s.Filename = fmt.Sprintf("%d - %s", s.Index, newName)
}

func (s *Song) RemoveIndex() string {
	return s.Name[3:]
}
