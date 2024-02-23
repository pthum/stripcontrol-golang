package telegram

import (
	"fmt"
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	alog "github.com/pthum/stripcontrol-golang/internal/log"
	"github.com/pthum/stripcontrol-golang/internal/service"
	"github.com/samber/do"
)

var rx_digits, _ = regexp.Compile(`\d+`)

type botCommand struct {
	Cmd         string
	Description string
	Action      func(*tgbotapi.Message) string
}
type cmdHandler struct {
	lsvc service.LEDService
	l    alog.Logger
}

func NewCmdHandler(i *do.Injector) *cmdHandler {
	lsvc := do.MustInvoke[service.LEDService](i)
	l := alog.NewLogger("cmdhandler")
	return &cmdHandler{
		lsvc: lsvc,
		l:    l,
	}
}
func (c *cmdHandler) GetCommands() []botCommand {
	return []botCommand{
		{
			Cmd:         "/ledon",
			Description: "Set an LED Strip to on, use: /ledon <id>",
			Action:      c.actionLedOn,
		},
		{
			Cmd:         "/ledoff",
			Description: "Set an LED Strip to off, use: /ledoff <id>",
			Action:      c.actionLedOff,
		},
		{
			Cmd:         "/getstrips",
			Description: "Returns all LED Strips, use: /getstrips",
			Action:      c.actionGetAll,
		},
		{
			Cmd:         "/help",
			Description: "Prints all commands, use: /help",
			Action:      c.actionGetCommands,
		},
	}
}

func (c *cmdHandler) actionLedOn(inp *tgbotapi.Message) string {
	return c.setLEDState(true, inp)
}

func (c *cmdHandler) actionLedOff(inp *tgbotapi.Message) string {
	return c.setLEDState(false, inp)
}

func (c *cmdHandler) actionGetAll(inp *tgbotapi.Message) string {
	strips, err := c.lsvc.GetAll()
	if err != nil {
		emsg := "Error getting all LED strips"
		c.l.Error(emsg+": %v", err)
		return emsg
	}
	msg := "All LED Strips: \n"
	for _, s := range strips {
		msg += fmt.Sprintf("ID: %s - '%s'\n", s.GetStringID(), s.Name)
	}
	return msg
}

func (c *cmdHandler) actionGetCommands(inp *tgbotapi.Message) string {
	aCmds := c.GetCommands()
	msg := "All available commands:\n"
	for _, cmd := range aCmds {
		msg += cmd.Cmd + " " + cmd.Description + "\n"
	}
	return msg
}

func (c *cmdHandler) setLEDState(enable bool, inp *tgbotapi.Message) string {
	procid := c.stripIdForMsg(inp.Text)
	action := "off"
	state := "Disabled"
	if enable {
		action = "on"
		state = "Enabled"
	}
	msg := fmt.Sprintf("Turning %s LED(s) with id %s\n", action, procid)
	for _, id := range procid {
		s, err := c.lsvc.GetLEDStrip(id)
		if err != nil {
			msg += fmt.Sprintf("Error getting ID %v\n", id)
			continue
		}
		s.Enabled = enable
		err = c.lsvc.UpdateLEDStrip(id, *s)
		if err != nil {
			msg += fmt.Sprintf("Error updating ID %v\n", id)
		} else {
			msg += fmt.Sprintf("%s ID %s\n", state, id)
		}
	}
	return msg
}

func (c *cmdHandler) stripIdForMsg(msg string) []string {
	procId := strings.TrimSpace(msg)
	if procId == "" {
		return []string{}
	}
	return rx_digits.FindAllString(procId, -1)
}

func (c *cmdHandler) callForMsg(msg *tgbotapi.Message) string {
	cmd := c.commandForMsg(msg)
	if cmd != nil {
		return cmd.Action(msg)
	}
	return ""
}

func (c *cmdHandler) commandForMsg(msg *tgbotapi.Message) *botCommand {
	if msg == nil || msg.Text == "" {
		return nil
	}
	for _, cmd := range c.GetCommands() {
		if strings.HasPrefix(msg.Text, cmd.Cmd) {
			return &cmd
		}
	}

	return nil
}
