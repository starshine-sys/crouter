package crouter

import (
	"errors"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var idRegex *regexp.Regexp

// Errors for parsing
var (
	ErrMemberNotFound  = errors.New("member not found")
	ErrChannelNotFound = errors.New("channel not found")
	ErrRoleNotFound    = errors.New("role not found")
	ErrNoID            = errors.New("input is not an ID")
	ErrBrokenMention   = errors.New("input is a broken/invalid mention")
)

func isID(id string) bool {
	if idRegex == nil {
		idRegex = regexp.MustCompile("([0-9]{15,})")
	}
	return idRegex.MatchString(id)
}

// ParseChannel takes a string and attempts to find a channel that matches that string
func (ctx *Ctx) ParseChannel(channel string) (*discordgo.Channel, error) {
	if isID(channel) {
		// ID Acting mode
		if strings.HasPrefix(channel, "<") {
			if !strings.HasPrefix(channel, "<#") || !strings.HasSuffix(channel, ">") {
				return nil, ErrBrokenMention
			}
			channel, _ = between(channel, "<#", ">")
		}
		c, err := ctx.Session.State.Channel(channel)
		if err == discordgo.ErrStateNotFound {
			c, err = ctx.Session.Channel(channel)
		}
		return c, err
	}

	channel = strings.ToLower(channel)

	// Try to find it by name
	g, err := ctx.Session.State.Guild(ctx.Message.GuildID)
	if err != nil {
		return nil, err
	}

	for _, x := range g.Channels {
		if strings.ToLower(x.Name) == channel {
			return x, nil
		}
	}

	return nil, ErrChannelNotFound
}

// ParseRole takes a string and attempts to find a role that matches that string
func (ctx *Ctx) ParseRole(role string) (*discordgo.Role, error) {
	if isID(role) {
		// ID Acting mode
		if strings.HasPrefix(role, "<") {
			if !strings.HasPrefix(role, "<@&") || !strings.HasSuffix(role, ">") {
				return nil, ErrBrokenMention
			}
			role, _ = between(role, "<@&", ">")
		}
		return ctx.Session.State.Role(ctx.Message.GuildID, role)
	}

	role = strings.ToLower(role)

	// Try to find it by name
	g, err := ctx.Session.State.Guild(ctx.Message.GuildID)
	if err != nil {
		return nil, err
	}

	for _, x := range g.Roles {
		if strings.ToLower(x.Name) == role {
			return x, nil
		}
	}

	return nil, ErrRoleNotFound
}

// ParseMember takes a string and attempts to find a member that matches that string
func (ctx *Ctx) ParseMember(member string) (*discordgo.Member, error) {
	if isID(member) {
		if strings.HasPrefix(member, "<") {
			if !strings.HasPrefix(member, "<@") || !strings.HasSuffix(member, ">") {
				return nil, errors.New("not a member mention or broken mention")
			}
			member, _ = between(member, "<@", ">")
			if member[0] == '!' {
				member = member[1:]
			}
		}
		m, err := ctx.Session.State.Member(ctx.Message.GuildID, member)
		if err == discordgo.ErrStateNotFound {
			m, err = ctx.Session.GuildMember(ctx.Message.GuildID, member)
		}
		return m, err
	}

	member = strings.ToLower(member)

	// Try to find it by name
	g, err := ctx.Session.State.Guild(ctx.Message.GuildID)
	if err != nil {
		return nil, err
	}

	for _, x := range g.Members {
		if strings.ToLower(x.User.String()) == member {
			return x, nil
		}
		if strings.ToLower(x.User.Username) == member {
			return x, nil
		}
		if strings.ToLower(x.Nick) == member {
			return x, nil
		}
	}

	return nil, ErrMemberNotFound
}

// ParseUser parses a user
func (ctx *Ctx) ParseUser(user string) (*discordgo.User, error) {
	if m, err := ctx.ParseMember(user); err == nil {
		return m.User, nil
	}

	// try parsing an off-server user
	if !isID(user) {
		return nil, ErrNoID
	}

	return ctx.Session.User(user)
}

// from: https://codeberg.org/evieDelta/drc/src/commit/c0facdfd1b017d5cbe181024a51f6da77bce7f61/detc/etc.go#L21
// between returns the contents of a string between before and after, yes its messy but it works
func between(s, after, before string) (string, bool) {
	if strings.Index(s, after) >= 0 && strings.Index(s, before) > 0 {
		return s[strings.Index(s, after)+len(after) : strings.Index(s, before)], true
	}
	return s, false
}
