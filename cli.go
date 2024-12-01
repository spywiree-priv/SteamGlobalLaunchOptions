package main

var CLI struct {
	Value        string `arg:""`
	UserId       string `short:"U" default:"${default_user_id}"`
	Overrite     bool   `short:"O" default:"false"`
	RestoreSteam *bool  `short:"R"`
}
