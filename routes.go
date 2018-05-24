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
	addRoute("main", mainAction)
}

func helpAction(req BotRequest) {
	req.QuickAnswer(fmt.Sprintf(strings.Join([]string{
		"Available commads:",
		"/help - Get this message",
		"/compile - Compile code",
		"/main - Compile code in main function",
	}, "\n")))
}

func pingAction(req BotRequest) {
	req.QuickAnswer("Pong!")
}

func compileAction(req BotRequest) {
	req.SendTyping()

	if len(strings.TrimSpace(req.args)) == 0 {
		req.QuickError("no input specified")
		return
	}

	output, err := compileAndRun(req.args)
	if err != nil {
		req.QuickError(err.Error())
		return
	}

	req.QuickAnswer(fmt.Sprintf("```\n%s\n```", output))
}

func mainAction(req BotRequest) {
	if len(strings.TrimSpace(req.args)) == 0 {
		req.QuickError("no input specified")
		return
	}

	codeTemplate := `
	package main

	func main() {
		%s
	}`

	req.args = fmt.Sprintf(codeTemplate, req.args)
	compileAction(req)
}

func compileAndRun(src string) (string, error) {
	events, err := executeSource(strings.TrimSpace(src))
	if err != nil {
		return "", err
	}

	var output string
	for _, event := range events {
		output += event.Message
	}

	output = strings.TrimSpace(output)

	if len(output) == 0 {
		return "[no output]", nil
	}

	return output, nil
}
