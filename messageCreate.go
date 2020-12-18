package crouter

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

// MessageCreate handles message create events. If you want more control over what happens when the event is received, you can instead call the (*Router).Context() and (*Router).Execute() functions manually.
func (r *Router) MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	var err error

	// if message was sent by a bot return
	if m.Author.Bot {
		return
	}

	// get context
	ctx, err := r.Context(m)
	if err != nil {
		log.Println("Error creating context:", err)
		return
	}

	// check if the message might be a command
	if ctx.MatchPrefix() {
		r.Execute(ctx)
	}
}
