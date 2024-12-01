package main

import (
	"os"
	"strings"
)

const ConfigFilePath = "config/localconfig.vdf"

func parsePath(path string) string {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}

		path = home + path[1:]
	}
	return path
}

func getDefaultUserId() string {
	dir, err := os.ReadDir(parsePath(SteamUserdata))
	if err != nil {
		return ""
	}

	if len(dir) != 1 {
		return ""
	}
	if !dir[0].IsDir() {
		return ""
	}

	return dir[0].Name()
}
