package main

import (
	"os"
	"runtime"
)

const ConfigFilePath = "config/localconfig.vdf"

func GetDefaultSteamUserdata() string {
	// https://help.steampowered.com/en/faqs/view/68D2-35AB-09A9-7678
	switch runtime.GOOS {
	case "windows":
		return `C:\Program Files (x86)\Steam\userdata`
	case "linux":
		// ~/.local/share/Steam/userdata
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		return home + `/.local/share/Steam/userdata`
	case "darwin":
		// ~/Library/Application Support/Steam/userdata
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		return home + `/Library/Application Support/Steam/userdata`
	default:
		return ""
	}
}

func GetDefaultUserId() string {
	dirpath := GetDefaultSteamUserdata()
	if dirpath == "" {
		return ""
	}

	dir, err := os.ReadDir(dirpath)
	if err != nil {
		return ""
	}

	if len(dir) == 1 && dir[0].IsDir() {
		return dir[0].Name()
	} else {
		return ""
	}
}
