package main

import (
	"strings"

	"github.com/Adeithe/go-twitch/irc"
)

type commandData struct {
	Name   string
	Handle func()
}

var modCommands []commandData = []commandData{
	{
		Name:   "start",
		Handle: startVote,
	},
	{
		Name:   "end",
		Handle: endVote,
	},
	{
		Name:   "clear",
		Handle: endVote,
	},
}

func handleModCommands(msg irc.ChatMessage, command string) bool {

	if msg.Sender.IsBroadcaster || msg.Sender.IsModerator {
		for _, item := range modCommands {
			if strings.EqualFold(item.Name, command) {
				item.Handle()
			}
		}
	}

	return false
}
