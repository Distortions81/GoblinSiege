package main

import (
	"strings"

	"github.com/gempir/go-twitch-irc/v4"
)

func handleChat(msg twitch.PrivateMessage) {

	gameLock.Lock()
	defer gameLock.Unlock()

	//Ignore messages that are not for us
	if msg.Channel != userSettings.UserName {
		return
	}

	qlog("Chat: %v: %v", msg.User.DisplayName, msg.Message)
	command, isCommand := strings.CutPrefix(msg.Message, userSettings.CmdPrefix)

	if isCommand {
		if handleModCommands(msg, command) {
			return

		} else {
			//If a vote is active, handle votes
			handleVoteMsg(msg, command)
		}
	}
}
