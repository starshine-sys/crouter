package crouter

import (
	"regexp"
	"strings"
)

// MatchPrefix checks if the message matched any prefix
func (ctx *Ctx) MatchPrefix() bool {
	return HasAnyPrefix(strings.ToLower(ctx.Message.Content), ctx.Router.Prefixes...)
}

// Match checks if any of the given command aliases match
func (ctx *Ctx) Match(cmds ...string) bool {
	for _, cmd := range cmds {
		if strings.ToLower(ctx.Command) == strings.ToLower(cmd) {
			return true
		}
	}
	return false
}

// MatchRegexp checks if the command matches the given regex
func (ctx *Ctx) MatchRegexp(re *regexp.Regexp) bool {
	if re == nil {
		return false
	}
	return re.MatchString(strings.ToLower(ctx.Command))
}
