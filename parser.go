package goldmark_discord_mentions

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

var discordMentionKind = ast.NewNodeKind("DiscordMention")

type discordMentionNode struct {
	ast.BaseInline

	ID string
}

func (d *discordMentionNode) Dump(source []byte, level int) {
	ast.DumpHelper(
		d,
		source,
		level,
		map[string]string{
			"ID": d.ID,
		},
		nil,
	)
}

func (d *discordMentionNode) Kind() ast.NodeKind {
	return discordMentionKind
}

var _ ast.Node = (*discordMentionNode)(nil)

type discordMentionParser struct{}

func (d discordMentionParser) Trigger() []byte {
	return []byte("<")
}

func (d discordMentionParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, _ := block.PeekLine()

	lineHead := 0
	userId := ""
	inMention := false

	for ; lineHead < len(line); lineHead++ {
		char := line[lineHead]

		if char == '@' {
			inMention = true
		} else if inMention && char == '>' {
			inMention = false
			break
		} else if inMention {
			userId += string(char)
		}
	}

	if inMention || userId == "" {
		return nil
	}

	block.Advance(lineHead + 1)

	return &discordMentionNode{
		ID: userId,
	}
}

var _ parser.InlineParser = (*discordMentionParser)(nil)
