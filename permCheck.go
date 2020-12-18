package crouter

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

// Errors relating to missing permissions
var (
	ErrorNoDMs           = errors.New("this command cannot be run in DMs")
	ErrorMissingPerms    = errors.New("you are missing required permissions")
	ErrorMissingBotOwner = errors.New("you are not a bot owner")
)

// Check checks if the user has permissions to run a command
func (ctx *Ctx) Check() (err error) {
	if ctx.Cmd.GuildOnly && ctx.Message.GuildID == "" {
		return ErrorNoDMs
	}
	if ctx.Cmd.OwnerOnly {
		return checkOwner(ctx.Author.ID, ctx.Router.BotOwners)
	}

	if ctx.Cmd.Permissions != 0 {
		return ctx.checkPerms(ctx.Author.ID)
	}
	return nil
}

func checkOwner(userID string, owners []string) (err error) {
	for _, u := range owners {
		if userID == u {
			return nil
		}
	}
	return ErrorMissingBotOwner
}

func (ctx *Ctx) checkPerms(userID string) (err error) {
	// check if in DMs
	if ctx.Message.GuildID == "" {
		return ErrorNoDMs
	}

	// get the guild
	guild, err := ctx.Session.State.Guild(ctx.Message.GuildID)
	if err == discordgo.ErrStateNotFound {
		guild, err = ctx.Session.Guild(ctx.Message.GuildID)
	}
	if err != nil && err != discordgo.ErrStateNotFound {
		return err
	}

	// get the member
	member, err := ctx.Session.State.Member(ctx.Message.GuildID, ctx.Author.ID)
	if err == discordgo.ErrStateNotFound {
		member, err = ctx.Session.GuildMember(ctx.Message.GuildID, ctx.Author.ID)
	}
	if err != nil && err != discordgo.ErrStateNotFound {
		return err
	}

	// if the user is the guild owner, they have permission to use the command
	if member.User.ID == guild.OwnerID {
		return nil
	}

	// iterate through all guild roles
	for _, r := range guild.Roles {
		// iterate through member roles
		for _, u := range member.Roles {
			// if they have the role...
			if u == r.ID {
				// ...and the role has the required perms, return
				if r.Permissions&ctx.Cmd.Permissions == ctx.Cmd.Permissions {
					return nil
				}
			}
		}
	}

	return ErrorMissingPerms
}

func checkPerms(p int, c ...int) bool {
	for _, perm := range c {
		if p&perm == perm {
			return true
		}
	}
	return false
}
