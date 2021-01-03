package crouter

import (
	"time"

	"github.com/ReneKroon/ttlcache/v2"
	"github.com/bwmarrin/discordgo"
)

// PermCache is a cache of channel permissions.
// Expires after 2 minutes by default, other timeouts can be set by using the SetTimeout() function
type PermCache struct {
	session *discordgo.Session
	cache   *ttlcache.Cache
}

// NewPermCache returns a new initialised permission cache
func NewPermCache(s *discordgo.Session) *PermCache {
	p := &PermCache{session: s}

	cache := ttlcache.NewCache()
	cache.SetTTL(2 * time.Minute)
	cache.SetCacheSizeLimit(10000)
	cache.SkipTTLExtensionOnHit(true)

	p.cache = cache
	return p
}

// Get gets the permissions in a channel
func (p *PermCache) Get(user, channel string) (perms int, err error) {
	if v, ok := p.cache.Get(user + "-" + channel); ok == nil {
		return v.(int), nil
	}

	perms, err = p.session.State.UserChannelPermissions(user, channel)
	if err == discordgo.ErrStateNotFound {
		perms, err = p.session.UserChannelPermissions(user, channel)
		if err == nil {
			p.cache.Set(user+"-"+channel, perms)
		}
		return perms, err
	}
	if err != nil {
		return perms, err
	}
	p.cache.Set(user+"-"+channel, perms)
	return perms, err
}

// HasPermissions checks if the user has the given permissions.
// If Get() returns an error, the user is assumed to not have the permissions.
func (p *PermCache) HasPermissions(user, channel string, required int) bool {
	perms, err := p.Get(user, channel)
	if err != nil {
		return false
	}
	return perms&required != 0
}
