package main

import (
	"strings"

	"github.com/gempir/go-twitch-irc/v4"
)

func handleChat(msg twitch.PrivateMessage) {

	qlog("Chat: %v: %v", msg.User.DisplayName, msg.Message)
	command, isCommand := strings.CutPrefix(msg.Message, userSettings.CmdPrefix)

	if isCommand {
		if handleModCommands(msg, command) {
			return
		} else if UserMsgDict.Voting {
			handleUserDictMsg(msg, command)
		}
	}
}
