package main

import (
	"log"

	"github.com/Starshine113/crouter"
)

// this function is called after every succesful command execution
func postFunc(ctx *crouter.Ctx) {
	log.Printf("[cmd] Command called: `%v`, arguments: %v\nUser: %v (%v), Channel: %v (%v), Guild: %v", ctx.Command, ctx.Args, ctx.Author.String(), ctx.Author.ID, ctx.Channel.Name, ctx.Channel.ID, ctx.Message.GuildID)
}
