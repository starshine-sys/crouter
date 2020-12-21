package crouter

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Execute actually executes the router.
// You shouldn't have to use this most of the time, as the (*Router).MessageCreate function calls
// this, but if you want more control over what happens, you probably want to call this yourself.
func (r *Router) Execute(ctx *Ctx) (err error) {
	// catch panics in commands
	defer func() {
		if r := recover(); r != nil {
			ctx.Embed("Panic", fmt.Sprintf("```%v```", r), 0xbf392f)
			log.Printf("Caught panic in %v (channel ID %v, guild %v): %v", ctx.Command, ctx.Message.ChannelID, ctx.Message.GuildID, r)
		}
	}()
	help := r.GetCommand("commands")
	if ctx.Match(append([]string{help.Name}, help.Aliases...)...) {
		err = r.Help(ctx)
		return
	}
	for _, g := range r.Groups {
		if ctx.Match(append([]string{g.Name}, g.Aliases...)...) || ctx.MatchRegexp(g.Regex) {
			if len(ctx.Args) == 0 {
				ctx.Command = ""
			} else {
				ctx.Command = ctx.Args[0]
			}
			if len(ctx.Args) > 1 {
				ctx.Args = ctx.Args[1:]
			} else {
				ctx.Args = []string{}
			}
			return g.execute(ctx)
		}
	}
	for _, cmd := range r.Commands {
		if ctx.Match(append([]string{cmd.Name}, cmd.Aliases...)...) || ctx.MatchRegexp(cmd.Regex) {
			if len(ctx.Args) > 0 {
				if ctx.Args[0] == "help" || ctx.Args[0] == "usage" {
					ctx.Args[0] = ctx.Command
					return r.Help(ctx)
				}
			}
			ctx.Cmd = cmd
			if cmd.Blacklistable {
				if !r.blacklist(ctx) {
					return nil
				}
			}
			if perms := ctx.Check(); perms != nil {
				return ctx.CommandError(perms)
			}
			for _, c := range cmd.CustomPermissions {
				if p, b := c(ctx); !b {
					_, err = ctx.Sendf("You are not allowed to use this command: you are missing the `%v` permission.", p)
					return err
				}
			}
			if cmd.Cooldown != time.Duration(0) {
				if _, e := r.Cooldowns.Get(fmt.Sprintf("%v-%v-%v", ctx.Channel.ID, ctx.Author.ID, cmd.Name)); e == nil {
					_, err = ctx.Sendf("The command `%v` can only be run once every **%v**.", cmd.Name, cmd.Cooldown.Round(time.Millisecond).String())
					return err
				}
				err = r.Cooldowns.SetWithTTL(fmt.Sprintf("%v-%v-%v", ctx.Channel.ID, ctx.Author.ID, cmd.Name), true, cmd.Cooldown)
				if err != nil {
					return err
				}
			}
			return cmd.Command(ctx)
		}
	}

	_, err = ctx.Send(&discordgo.MessageSend{
		Content: fmt.Sprintf("%v Unknown command `%v`. For a list of commands, try `%v%v`.", ErrorEmoji, ctx.Command, r.Prefixes[0], help.Name),
		AllowedMentions: &discordgo.MessageAllowedMentions{
			Parse: []discordgo.AllowedMentionType{
				discordgo.AllowedMentionTypeUsers,
			},
		},
	})
	return
}
