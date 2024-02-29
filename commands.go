package main

import "github.com/Adeithe/go-twitch/irc"

func adminCommands(msg irc.ChatMessage) bool {

	if msg.Sender.IsBroadcaster || msg.Sender.IsModerator {
		//Eat this

		//return true
	}

	return false
}
