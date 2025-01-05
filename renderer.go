package goldmark_discord_mentions

import (
	"fmt"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type discordMentionRenderer struct {
	resolver MentionResolver
}

func (r *discordMentionRenderer) render(
	w util.BufWriter, source []byte, n ast.Node, entering bool,
) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	username := fmt.Sprintf("<@%s>", n.(*discordMentionNode).ID)

	if r.resolver != nil {
		username = r.resolver.ResolveUser(n.(*discordMentionNode).ID)
	}

	// TODO: Is this vulnerable to XSS?
	// TODO: Wrap in span & allow customizing class?
	_, err := w.WriteString(fmt.Sprintf("@%s", username))
	if err != nil {
		return ast.WalkStop, fmt.Errorf("failed to write string: %w", err)
	}

	return ast.WalkContinue, nil
}

func (r *discordMentionRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(discordMentionKind, r.render)
}

var _ renderer.NodeRenderer = (*discordMentionRenderer)(nil)
