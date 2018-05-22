package main

import (
	"fmt"
	"strings"
)

func addRoutes() {
	addRoute("help", helpAction)
	addRoute("start", helpAction)
	addRoute("ping", pingAction)
}

func helpAction(req BotRequest) {
	req.QuickAnswer(fmt.Sprintf(strings.Join([]string{
		"Available commads:",
		"/help - Get this message",
	}, "\n")))
}

func pingAction(req BotRequest) {
	req.QuickAnswer("Pong!")
}
