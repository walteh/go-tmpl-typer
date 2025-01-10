package lsp

import (
	"context"
	"encoding/json"

	"github.com/sourcegraph/jsonrpc2"
	"github.com/walteh/go-tmpl-typer/pkg/parser"
	"github.com/walteh/go-tmpl-typer/pkg/types"
	"gitlab.com/tozd/go/errors"
)

func (s *Server) handleTextDocumentHover(ctx context.Context, req *jsonrpc2.Request) (interface{}, error) {
	if s.debug {
		s.debugf(ctx, "handling textDocument/hover")
	}

	var params HoverParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, errors.Errorf("failed to unmarshal hover params: %w", err)
	}

	// Get document content
	content, ok := s.getDocument(s.normalizeURI(params.TextDocument.URI))
	if !ok {
		return nil, errors.Errorf("document not found: %s", params.TextDocument.URI)
	}

	// Parse the template
	info, err := s.server.parser.Parse(ctx, []byte(content), params.TextDocument.URI)
	if err != nil {
		return nil, errors.Errorf("parsing template for hover: %w", err)
	}

	// Find type hint for the position
	hint := findTypeHintForPosition(info, params.Position.Line+1, params.Position.Character+1)
	if hint == nil {
		return nil, nil
	}

	// Get type info
	typeInfo, err := s.server.validator.ValidateType(ctx, hint.TypePath)
	if err != nil {
		return nil, errors.Errorf("validating type for hover: %w", err)
	}

	// Create hover response
	hover := &Hover{
		Contents: MarkupContent{
			Kind:  "markdown",
			Value: "Type: " + typeInfo.Name + "\n\nFields:\n" + formatFieldsMarkdown(typeInfo.Fields),
		},
	}

	return hover, nil
}

func findTypeHintForPosition(info *parser.TemplateInfo, line, character int) *parser.TypeHint {
	if len(info.TypeHints) == 0 {
		return nil
	}

	// For now, just return the first type hint
	// TODO: Implement proper position-based type hint lookup
	return &info.TypeHints[0]
}

func formatFieldsMarkdown(fields map[string]types.FieldInfo) string {
	var result string
	for name, field := range fields {
		result += "- " + name + ": " + field.Type + "\n"
	}
	return result
}
