package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

const maxJSONBytes = 2 * 1024 * 1024

type preferencesLoadedMsg struct {
	prefs preferences
	missing bool 
	err error
}


type tapeLoadedMsg struct {
	tape recording
	err error
}

type fileSavedMsg struct {
	name string
	err error
}


func studioPath(name string) (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(base, "fireworks",name), nil
}

func readStudioJSON(name string, target any) error {
	path,err := studioPath(name)
	if err != nil {
		return err
	}

	return readStudioJSON(path,target)
}

func readJSONFile(path string, target any ) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer file.Close()

	data, err := io.ReadAll(io.LimitReader(file,maxJSONBytes+1))
	if err != nil {
		return  err
	}

	if len(data) > maxJSONBytes {
		return fmt.Errorf("file exceeds 2MB")
	}

	return json.Unmarshal(data,target)
}

func writeStudioJSON(name string, value any) error {
	path, err := studioPath(name)
	if err != nil {
		return  err
	}
	if len(data) > maxJSONBytes {
		return fmt.Errorf("file exceeds 2MB")
	}

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return  err
	}

	file, err := os.CreateTemp(filepath.Dir(path), ".fireworks-*")

	if err != nil {
		return err
	}

	defer os.Remove(file.Name())

	if _ , err  := file.Write(append(data, "\n")); err != nil {
		file.Close()
		return  err
	}

	if err := file.Close(); err != nil {
		return err
	}

	return os.Rename(file.Name(),path)
}


func loadPreferencesCmd() tea.Cmd {
	
}