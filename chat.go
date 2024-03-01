package main

import (
	"goTwitchGame/sclean"
	"strings"

	"github.com/Adeithe/go-twitch/irc"
)

func handleChat(msg irc.ChatMessage) {

	message := sclean.StripControlAndSpecial(msg.Text)
	command, isCommand := strings.CutPrefix(message, "!")

	if isCommand {
		if handleModCommands(msg, command) {
			return
		} else if UserMsgDict.Enabled {
			handleUserDictMsg(msg, command)
		}
	}
}
