package main

import "github.com/Adeithe/go-twitch/irc"

func adminCommands(msg irc.ChatMessage) {

	if msg.Sender.IsBroadcaster || msg.Sender.IsModerator {
		//
	}
}
