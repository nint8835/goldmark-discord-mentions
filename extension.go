package goldmark_discord_mentions

import "github.com/yuin/goldmark"

type Extension struct{}

func (e *Extension) Extend(m goldmark.Markdown) {}

var _ goldmark.Extender = &Extension{}
