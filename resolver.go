package goldmark_discord_mentions

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

type MentionResolver interface {
	ResolveUser(id string) string
}

type cachedMention struct {
	Username  string
	CacheTime time.Time
}

type cachedResolver struct {
	session     *discordgo.Session
	cachedUsers map[string]cachedMention
	cachePath   string
}

func (r *cachedResolver) ResolveUser(id string) string {
	// TODO: Configurable cache time
	if cachedVal, ok := r.cachedUsers[id]; ok && time.Since(cachedVal.CacheTime) < time.Hour*24*7 {
		return cachedVal.Username
	}

	var username string
	user, err := r.session.User(id)
	if err != nil {
		username = fmt.Sprintf("<@%s>", id)
	} else {
		username = user.Username
	}

	r.cachedUsers[id] = cachedMention{
		Username:  username,
		CacheTime: time.Now(),
	}

	if r.cachePath != "" {
		f, err := os.Create(r.cachePath)
		if err != nil {
			// TODO: Better error handling
			return username
		}
		defer f.Close()

		err = json.NewEncoder(f).Encode(r.cachedUsers)
		if err != nil {
			// TODO: Better error handling
			return username
		}
	}

	return username
}

var _ MentionResolver = &cachedResolver{}
