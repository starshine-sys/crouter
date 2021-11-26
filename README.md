# This is outdated! You shouldn't use it--if you desperately wanna use our spaghetti, try [bcr](https://github.com/starshine-sys/bcr)

# crouter

[![Go Reference](https://pkg.go.dev/badge/github.com/Starshine113/crouter.svg)](https://pkg.go.dev/github.com/Starshine113/crouter)

A simple command handler for [discordgo](https://github.com/bwmarrin/discordgo).

## Usage

```go
import "github.com/Starshine113/crouter"

// ... session initialisation code ...

// create the router
r := crouter.NewRouter(dg, []string{botOwner}, []string{"?", "!"})

// add the message create handler
dg.AddHandler(r.MessageCreate)

// add a command
r.AddCommand(&crouter.Command{
    Name: "Ping",

    Summary: "Check if the bot is running",

    Command: func(ctx *crouter.Ctx) (err error) {
        ping := ctx.Session.HeartbeatLatency().Round(time.Millisecond)
        _, err = ctx.Sendf("Ping! Average latency: %s", ping)
        return err
    },
})

// add intents (or just add crouter.RequiredIntents to your existing intents)
dg.Identify.Intents = discordgo.MakeIntent(crouter.RequiredIntents)

// open the session
err = dg.Open()

// ...
```

A more complete example can be found in the `example/` directory.
