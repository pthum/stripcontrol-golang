package telegram

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type botCommand struct {
	Cmd         string
	Description string
	Action      func(*tgbotapi.Message) string
}

func GetCommands() []botCommand {
	return []botCommand{
		{
			Cmd:         "/ledon",
			Description: "Set an LED Strip to on, use: /ledon <id>",
			Action: func(inp *tgbotapi.Message) string {
				return "Turning on LED with id"
			},
		},
		{
			Cmd:         "/ledoff",
			Description: "Set an LED Strip to off, use: /ledoff <id>",
			Action: func(inp *tgbotapi.Message) string {
				return "Turning off LED with id"
			},
		},
		{
			Cmd:         "/getall",
			Description: "Returns all LED Strips, use: /getstrips",
			Action: func(inp *tgbotapi.Message) string {
				return "All LED Strips: \n"
			},
		},
		{
			Cmd:         "/help",
			Description: "Prints all commands",
			Action: func(inp *tgbotapi.Message) string {
				aCmds := GetCommands()
				msg := "All available commands:\n"
				for _, cmd := range aCmds {
					msg += cmd.Cmd + " " + cmd.Description + "\n"
				}
				return msg
			},
		},
	}
}

func callForMsg(msg *tgbotapi.Message) string {
	cmd := commandForMsg(msg)
	if cmd != nil {
		return cmd.Action(msg)
	}
	return ""
}

func commandForMsg(msg *tgbotapi.Message) *botCommand {
	if msg == nil || msg.Text == "" {
		return nil
	}
	for _, cmd := range GetCommands() {
		if strings.HasPrefix(msg.Text, cmd.Cmd) {
			return &cmd
		}
	}

	return nil
}
