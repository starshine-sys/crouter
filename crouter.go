package crouter

import (
	"regexp"
	"strings"
	"time"

	"github.com/ReneKroon/ttlcache/v2"
	"github.com/bwmarrin/discordgo"
)

// RequiredIntents are the intents required for the router
const RequiredIntents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildMessageReactions | discordgo.IntentsDirectMessages | discordgo.IntentsDirectMessageReactions | discordgo.IntentsGuilds

// Router is the command router
type Router struct {
	Commands []*Command
	Groups   []*Group

	BotOwners []string
	Prefixes  []string

	Session *discordgo.Session

	Cooldowns *ttlcache.Cache
	Handlers  *ttlcache.Cache

	prefixUsersSet bool
}

// Command is a single command
type Command struct {
	Name    string
	Aliases []string
	Regex   *regexp.Regexp

	// Summary is used in the command list
	Summary string
	// Description is used in the help command
	Description string
	// Usage is appended to the command name in help commands
	Usage string

	// Command is the actual command function
	Command func(*Ctx) error

	Permissions int
	GuildOnly   bool
	OwnerOnly   bool
	Cooldown    time.Duration

	Router *Router
}

// NewRouter creates a Router object
func NewRouter(s *discordgo.Session, owners, prefixes []string) *Router {
	cache := ttlcache.NewCache()
	cache.SkipTTLExtensionOnHit(true)

	h := ttlcache.NewCache()
	h.SetCacheSizeLimit(10000)
	h.SetTTL(15 * time.Minute)
	h.SetExpirationCallback(func(key string, value interface{}) {
		value.(func())()
	})

	router := &Router{
		BotOwners: owners,
		Session:   s,
		Cooldowns: cache,
		Handlers:  h,
		Prefixes:  prefixes,
	}

	router.AddCommand(&Command{
		Name:        "Commands",
		Description: "Show a list of commands",
		Usage:       "[command]",
		Command:     router.dummy,
	})

	return router
}

// dummy is used when a command isn't handled with the normal process
func (r *Router) dummy(ctx *Ctx) error {
	return nil
}

// AddCommand adds a command to the router
func (r *Router) AddCommand(cmd *Command) {
	cmd.Router = r
	if cmd.Cooldown == 0 {
		cmd.Cooldown = 500 * time.Millisecond
	}
	r.Commands = append(r.Commands, cmd)
}

// GetCommand gets a command by name
func (r *Router) GetCommand(name string) (c *Command) {
	for _, cmd := range r.Commands {
		if strings.ToLower(cmd.Name) == strings.ToLower(name) {
			return cmd
		}
		for _, a := range cmd.Aliases {
			if strings.ToLower(a) == strings.ToLower(name) {
				return cmd
			}
		}
	}
	return nil
}
