package crouter

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type cmdList []*Command

func (c cmdList) Len() int      { return len(c) }
func (c cmdList) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c cmdList) Less(i, j int) bool {
	return sort.StringsAreSorted([]string{c[i].Name, c[j].Name})
}

// Help is the help command
func (r *Router) Help(ctx *Ctx) (err error) {
	err = ctx.TriggerTyping()
	if err != nil {
		return
	}

	if len(ctx.Args) == 0 {
		return r.details(ctx)
	}

	var cmd *Command
	g := r.GetGroup(ctx.Args[0])
	if g != nil {
		if len(ctx.Args) == 1 {
			_, err = ctx.Send(ctx.groupEmbed(g))
			return
		}
		cmd = g.GetCommand(ctx.Args[1])
		if cmd != nil {
			_, err = ctx.Send(ctx.groupCmdEmbed(g, cmd))
			return
		}
	}
	cmd = r.GetCommand(ctx.Args[0])
	if cmd != nil {
		_, err = ctx.Send(ctx.cmdEmbed(cmd))
		return
	}

	_, err = ctx.Send(fmt.Sprintf("%v Invalid command or group provided:\n> `%v` is not a known command, group or alias.", ErrorEmoji, ctx.Args[0]))

	return
}

func (ctx *Ctx) groupEmbed(g *Group) *discordgo.MessageEmbed {
	var aliases string
	if g.Aliases == nil {
		aliases = "N/A"
	} else {
		aliases = strings.Join(g.Aliases, ", ")
	}

	var subCmds []string
	for _, cmd := range g.Subcommands {
		subCmds = append(subCmds, fmt.Sprintf("[%d] %s", cmd.Permissions, cmd.Name))
	}

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("```%v```", strings.ToUpper(g.Name)),
		Description: g.Description,
		Color:       0x21a1a8,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Subcommands",
				Value:  fmt.Sprintf("```%v```", strings.Join(subCmds, "\n")),
				Inline: false,
			},
			{
				Name:   "Aliases",
				Value:  fmt.Sprintf("```%v```\n** **", aliases),
				Inline: false,
			},
			{
				Name:   "Default command",
				Value:  g.Command.Description,
				Inline: false,
			},
			{
				Name:   "Usage",
				Value:  fmt.Sprintf("```%v%v %v```", ctx.Router.Prefixes[0], strings.ToLower(g.Name), g.Command.Usage),
				Inline: false,
			},
			{
				Name:   "Required permissions",
				Value:  "```" + PermToString(g.Command.Permissions) + "```",
				Inline: false,
			},
		},
	}

	return embed
}

func (ctx *Ctx) groupCmdEmbed(g *Group, cmd *Command) *discordgo.MessageEmbed {
	var aliases string

	if cmd.Aliases == nil {
		aliases = "N/A"
	} else {
		aliases = strings.Join(cmd.Aliases, ", ")
	}

	fields := make([]*discordgo.MessageEmbedField, 0)

	if cmd.Description != "" {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Description",
			Value:  cmd.Description,
			Inline: false,
		})
	}

	fields = append(fields, []*discordgo.MessageEmbedField{
		{
			Name:   "Usage",
			Value:  fmt.Sprintf("```%v%v %v %v```", ctx.Router.Prefixes[0], strings.ToLower(g.Name), strings.ToLower(cmd.Name), cmd.Usage),
			Inline: false,
		},
		{
			Name:   "Aliases",
			Value:  fmt.Sprintf("```%v```", aliases),
			Inline: false,
		},
		{
			Name:   "Required permissions",
			Value:  "```" + PermToString(cmd.Permissions) + "```",
			Inline: false,
		},
	}...)

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("```%v %v```", strings.ToUpper(g.Name), strings.ToUpper(cmd.Name)),
		Description: cmd.Description,
		Color:       0x21a1a8,
		Fields:      fields,
	}

	return embed
}

func (ctx *Ctx) cmdEmbed(cmd *Command) *discordgo.MessageEmbed {
	var aliases string

	if cmd.Aliases == nil {
		aliases = "N/A"
	} else {
		aliases = strings.Join(cmd.Aliases, ", ")
	}

	fields := make([]*discordgo.MessageEmbedField, 0)

	if cmd.Description != "" {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Description",
			Value:  cmd.Description,
			Inline: false,
		})
	}

	fields = append(fields, []*discordgo.MessageEmbedField{
		{
			Name:   "Usage",
			Value:  fmt.Sprintf("```%v%v %v```", ctx.Router.Prefixes[0], strings.ToLower(cmd.Name), cmd.Usage),
			Inline: false,
		},
		{
			Name:   "Aliases",
			Value:  fmt.Sprintf("```%v```", aliases),
			Inline: false,
		},
		{
			Name:   "Required permissions",
			Value:  "```" + PermToString(cmd.Permissions) + "```",
			Inline: false,
		},
	}...)

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("```%v```", strings.ToUpper(cmd.Name)),
		Description: cmd.Description,
		Color:       0x21a1a8,
		Fields:      fields,
	}

	return embed
}

func (r *Router) details(ctx *Ctx) (err error) {
	if err = ctx.TriggerTyping(); err != nil {
		return err
	}

	var cmds cmdList
	for _, c := range r.Commands {
		cmds = append(cmds, c)
	}

	for _, g := range r.Groups {
		cmds = append(cmds, &Command{
			Name:        g.Name,
			Permissions: g.Command.Permissions,
			Description: g.Command.Description,
		})
	}

	sort.Sort(cmds)
	cmdSlices := make([][]*Command, 0)

	for i := 0; i < len(cmds); i += 10 {
		end := i + 10

		if end > len(cmds) {
			end = len(cmds)
		}

		cmdSlices = append(cmdSlices, cmds[i:end])
	}

	embeds := make([]*discordgo.MessageEmbed, 0)

	for i, c := range cmdSlices {
		x := make([]string, 0)
		for _, cmd := range c {
			x = append(x, fmt.Sprintf("`[%d] %v`: %v", cmd.Permissions, cmd.Name, cmd.Description))
		}
		embeds = append(embeds, &discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{
				Name:    ctx.BotUser.Username + " help",
				IconURL: ctx.BotUser.AvatarURL("128"),
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("Page %v/%v", i+1, len(cmdSlices)),
			},
			Timestamp:   time.Now().Format(time.RFC3339),
			Description: strings.Join(x, "\n"),
			Fields: []*discordgo.MessageEmbedField{{
				Name:   "Usage",
				Value:  "Use ⬅️ ➡️ to navigate between pages, and use ❌ to delete this message.",
				Inline: false,
			}},
			Color: 0x21a1a8,
		})
	}

	_, err = ctx.PagedEmbed(embeds)
	return
}

// PrettyDurationString turns a duration into a string representation (maximum of days)
func PrettyDurationString(duration time.Duration) (out string) {
	var days, hours, hoursFrac, minutes float64

	hours = duration.Hours()
	hours, hoursFrac = math.Modf(hours)
	minutes = hoursFrac * 60

	hoursFrac = math.Mod(hours, 24)
	days = (hours - hoursFrac) / 24
	hours = hours - (days * 24)
	minutes = minutes - math.Mod(minutes, 1)

	if days != 0 {
		out += fmt.Sprintf("%v days, ", days)
	}
	if hours != 0 {
		out += fmt.Sprintf("%v hours, ", hours)
	}
	out += fmt.Sprintf("%v minutes", minutes)

	return
}
