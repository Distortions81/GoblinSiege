package main

import (
	"goTwitchGame/sclean"
	"strings"

	"github.com/gempir/go-twitch-irc/v4"
)

func handleChat(msg twitch.PrivateMessage) {

	qlog("%v: %v", msg.User.DisplayName, msg.Message)

	message := sclean.StripControlAndSpecial(msg.Message)
	command, isCommand := strings.CutPrefix(message, "!")

	if isCommand {
		if handleModCommands(msg, command) {
			return
		} else if UserMsgDict.Voting {
			handleUserDictMsg(msg, command)
		}
	}
}
