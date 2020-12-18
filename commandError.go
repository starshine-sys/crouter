package crouter

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

// CommandError sends an error message and optionally returns an error for logging purposes
func (ctx *Ctx) CommandError(err error) error {
	switch err {
	case ErrorNoDMs:
		_, err := ctx.Sendf("%v This command cannot be used in DMs.", ErrorEmoji)
		return err
	case ErrorMissingBotOwner:
		_, err := ctx.Sendf("%v You need to be the bot owner to use this command.", ErrorEmoji)
		return err
	case ErrorMissingPerms:
		_, err := ctx.Sendf("%v You need the following permission to use this command:\n> %v", ErrorEmoji, PermToString(ctx.Cmd.Permissions))
		return err
	case ErrorNotEnoughArgs:
		ctx.React(WarnEmoji)
		_, msgErr := ctx.Send(&discordgo.MessageSend{
			Content: WarnEmoji + " Command call was missing arguments.",
			AllowedMentions: &discordgo.MessageAllowedMentions{
				Parse: []discordgo.AllowedMentionType{},
			},
		})
		return msgErr
	case ErrorTooManyArgs:
		ctx.React(WarnEmoji)
		_, msgErr := ctx.Send(&discordgo.MessageSend{
			Content: WarnEmoji + " Command call has too many arguments.",
			AllowedMentions: &discordgo.MessageAllowedMentions{
				Parse: []discordgo.AllowedMentionType{},
			},
		})
		return msgErr
	}
	switch err.(type) {
	case *discordgo.RESTError:
		e := err.(*discordgo.RESTError)
		if e.Message != nil {
			_, err = ctx.Send(&discordgo.MessageEmbed{
				Title:       "REST error occurred",
				Description: fmt.Sprintf("```%v ```", e.Message.Message),
				Fields: []*discordgo.MessageEmbedField{{
					Name:   "Raw",
					Value:  string(e.ResponseBody),
					Inline: false,
				}},
				Footer: &discordgo.MessageEmbedFooter{
					Text: fmt.Sprintf("Error code: %v", e.Message.Code),
				},
				Color:     0xbf1122,
				Timestamp: time.Now().UTC().Format(time.RFC3339),
			})
		} else {
			_, err = ctx.Send(&discordgo.MessageEmbed{
				Title:       "REST error occurred",
				Description: fmt.Sprintf("```%v```", e.ResponseBody),
				Color:       0xbf1122,
				Timestamp:   time.Now().UTC().Format(time.RFC3339),
			})
		}
		return err
	default:
		ctx.Session.MessageReactionAdd(ctx.Message.ChannelID, ctx.Message.ID, ErrorEmoji)

		embed := &discordgo.MessageEmbed{
			Title:       "Internal error occured",
			Description: fmt.Sprintf("```%v```\nIf this error persists, please contact the bot developer.", err.Error()),
			Color:       0xbf1122,
			Timestamp:   time.Now().UTC().Format(time.RFC3339),
		}
		_, err = ctx.Send(embed)
		return err
	}
}
