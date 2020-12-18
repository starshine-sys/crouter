package crouter

import (
	"errors"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	// SuccessEmoji is the emoji used to designate a successful action
	SuccessEmoji = "✅"
	// ErrorEmoji is the emoji used for errors
	ErrorEmoji = "❌"
	// WarnEmoji is the emoji used to warn that a command failed
	WarnEmoji = "⚠️"
)

// Ctx is the context for a command
type Ctx struct {
	Command string
	Args    []string
	RawArgs string

	Session *discordgo.Session
	BotUser *discordgo.User

	Message *discordgo.MessageCreate
	Channel *discordgo.Channel
	Author  *discordgo.User

	Cmd    *Command
	Router *Router

	AdditionalParams map[string]interface{}
}

// Errors when creating Context
var (
	ErrorNoBotUser = errors.New("bot user not found in state cache")
)

// Context creates a new Ctx
func (r *Router) Context(m *discordgo.MessageCreate) (ctx *Ctx, err error) {
	if r.Session.State.User == nil {
		return nil, ErrorNoBotUser
	}

	if !r.prefixUsersSet {
		r.Prefixes = append(r.Prefixes, "<@"+r.Session.State.User.ID+">", "<@!"+r.Session.State.User.ID+">")
		r.prefixUsersSet = true
	}
	messageContent := TrimPrefixesSpace(m.Content, r.Prefixes...)
	message := strings.Split(messageContent, " ")
	command := message[0]
	args := []string{}
	if len(message) > 1 {
		args = message[1:]
	}

	ctx = &Ctx{Command: command, Args: args, Message: m, Author: m.Author, RawArgs: TrimPrefixesSpace(m.Content, append(r.Prefixes, command)...), Router: r, Session: r.Session}

	channel, err := r.Session.Channel(m.ChannelID)
	if err != nil {
		return ctx, err
	}
	ctx.Channel = channel
	ctx.AdditionalParams = make(map[string]interface{})
	ctx.BotUser = r.Session.State.User

	return ctx, nil
}
