package main

import (
	"errors"
	"log"
	"os/exec"
	"path/filepath"

	"github.com/alecthomas/kong"
	"github.com/shirou/gopsutil/v4/process"
)

var ErrProcessNotFound = errors.New("process not found")

func KillProcess(name string) ([]string, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, err
	}
	for _, p := range processes {
		n, err := p.Name()
		if err != nil {
			return nil, err
		}
		if n == name {
			args, err := p.CmdlineSlice()
			if err != nil {
				return nil, err
			}

			return args, p.Kill()
		}
	}
	return nil, ErrProcessNotFound
}

var CLI struct {
	Value string `arg:""`

	SteamUserdata string `type:"existingdir" default:"${default_steam_userdata}"`
	UserId        string `default:"${default_user_id}"`

	Overwrite        bool `short:"O"`
	DontRestartSteam bool `short:"R"`
}

func main() {
	_ = kong.Parse(&CLI,
		kong.Name("SteamGlobalLaunchOptions"),
		kong.Description("A CLI tool to apply launch options for all Steam games at once."),
		kong.UsageOnError(),
		kong.ConfigureHelp(
			kong.HelpOptions{
				Compact: true,
				Summary: true,
			},
		),
		kong.Vars{
			"default_steam_userdata": GetDefaultSteamUserdata(),
			"default_user_id":        GetDefaultUserId(),
		},
	)
	if CLI.UserId == "" {
		log.Fatalln("Default user id have not been found. You must provide user id.")
	} else if CLI.SteamUserdata == "" {
		log.Fatalln("Default Steam userdata folder path have not been found. You must provide Steam userdata folder path.")
	}

	args, err := KillProcess("steam")
	if errors.Is(err, ErrProcessNotFound) {
		CLI.DontRestartSteam = true
	} else if err != nil {
		log.Fatalln(err)
	}

	err = ApplyLaunchOptions(
		CLI.Value,
		filepath.Join(CLI.SteamUserdata, CLI.UserId, ConfigFilePath),
		CLI.Overwrite,
	)
	if err != nil {
		log.Fatalln(err)
	}

	if !CLI.DontRestartSteam {
		if len(args) <= 1 {
			_ = exec.Command("steam").Start()
		} else {
			_ = exec.Command("steam", args[1:]...).Start()
		}
	}
}
