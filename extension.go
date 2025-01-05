package goldmark_discord_mentions

import (
	"encoding/json"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type Extension struct {
	resolver MentionResolver
}

func (e *Extension) Extend(m goldmark.Markdown) {
	m.Renderer().
		AddOptions(renderer.WithNodeRenderers(util.Prioritized(&discordMentionRenderer{resolver: e.resolver}, 500)))
	m.Parser().
		AddOptions(parser.WithInlineParsers(util.Prioritized(&discordMentionParser{}, 500)))
}

var _ goldmark.Extender = &Extension{}

func New(session *discordgo.Session, cachePath string) *Extension {
	// TODO: Enable providing a guild to pre-fill the cache for
	cachedUsers := make(map[string]cachedMention)

	if cachePath != "" {
		f, err := os.Open(cachePath)
		if err != nil {
			// TODO: Log?
		}
		defer f.Close()

		err = json.NewDecoder(f).Decode(&cachedUsers)
		if err != nil {
			// TODO: Log?
		}
	}

	resolver := &cachedResolver{
		session:     session,
		cachedUsers: make(map[string]cachedMention),
		cachePath:   cachePath,
	}

	return &Extension{resolver: resolver}
}
