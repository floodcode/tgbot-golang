package main

import (
	"fmt"
	"strings"
)

func addRoutes() {
	addRoute("help", helpAction)
	addRoute("start", helpAction)
	addRoute("ping", pingAction)
	addRoute("compile", compileAction)
}

func helpAction(req BotRequest) {
	req.QuickAnswer(fmt.Sprintf(strings.Join([]string{
		"Available commads:",
		"/help - Get this message",
		"/compile - Compile code (from newline)",
	}, "\n")))
}

func pingAction(req BotRequest) {
	req.QuickAnswer("Pong!")
}

func compileAction(req BotRequest) {
	events, err := executeSource(req.args)
	if err != nil {
		req.QuickAnswer(err.Error())
		return
	}

	output := ""
	for _, event := range events {
		output += event.Message
	}

	if len(output) == 0 {
		req.QuickAnswer("Program was executed with no output")
		return
	}

	req.QuickAnswer(fmt.Sprintf("*Output:*\n```\n%s\n```", output))
}
