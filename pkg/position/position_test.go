package position_test

import (
	"testing"
	"text/template/parse"

	"github.com/walteh/go-tmpl-typer/pkg/position"
)

func TestGetLineAndColumn(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		pos      parse.Pos
		wantLine int
		wantCol  int
	}{
		{
			name:     "empty text",
			text:     "",
			pos:      parse.Pos(0),
			wantLine: 1,
			wantCol:  1,
		},
		{
			name:     "single line, first position",
			text:     "Hello, World! ",
			pos:      parse.Pos(2),
			wantLine: 1,
			wantCol:  3,
		},
		{
			name:     "single line, middle position",
			text:     "Hello, World!",
			pos:      parse.Pos(7),
			wantLine: 1,
			wantCol:  8,
		},
		{
			name:     "multiple lines, first line",
			text:     "Hello\nWorld\nTest",
			pos:      parse.Pos(3),
			wantLine: 1,
			wantCol:  4,
		},
		{
			name:     "multiple lines, second line",
			text:     "Hello\nWorld\nTest zzz",
			pos:      parse.Pos(8),
			wantLine: 2,
			wantCol:  3,
		},
		{
			name:     "multiple lines with varying lengths",
			text:     "Hello, World!\nThis is a test\nShort\nLonger line here zzz",
			pos:      parse.Pos(16),
			wantLine: 2,
			wantCol:  3,
		},
		{
			name:     "broken example",
			text:     "{{- /*gotype: test.Person*/ -}}\nAddress:\n  Street: {{.Address.Street}}",
			pos:      parse.Pos(61),
			wantLine: 3,
			wantCol:  21,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLine, gotCol := position.GetLineAndColumn(tt.text, tt.pos)
			if gotLine != tt.wantLine || gotCol != tt.wantCol {
				t.Errorf("GetLineAndColumn() = (%v, %v), want (%v, %v)", gotLine, gotCol, tt.wantLine, tt.wantCol)
			}
		})
	}
}

func TestRawPosition(t *testing.T) {
	tests := []struct {
		name     string
		pos      position.RawPosition
		wantText string
		wantID   string
	}{
		{
			name: "basic position",
			pos: position.RawPosition{
				Text:   "test",
				Offset: 10,
			},
			wantText: "test",
			wantID:   "test@10",
		},
		{
			name: "empty text",
			pos: position.RawPosition{
				Text:   "",
				Offset: 0,
			},
			wantText: "",
			wantID:   "@0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pos.Text; got != tt.wantText {
				t.Errorf("RawPosition.Text = %v, want %v", got, tt.wantText)
			}
			if got := tt.pos.ID(); got != tt.wantID {
				t.Errorf("RawPosition.ID() = %v, want %v", got, tt.wantID)
			}
		})
	}
}

func TestHasRangeOverlap(t *testing.T) {
	tests := []struct {
		name      string
		pos       position.RawPosition
		start     position.RawPosition
		wantMatch bool
	}{
		{
			name: "overlapping ranges",
			pos: position.RawPosition{
				Text:   "test",
				Offset: 5,
			},
			start: position.RawPosition{
				Text:   "testing",
				Offset: 3,
			},
			wantMatch: true,
		},
		{
			name: "non-overlapping ranges",
			pos: position.RawPosition{
				Text:   "test",
				Offset: 10,
			},
			start: position.RawPosition{
				Text:   "test",
				Offset: 0,
			},
			wantMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pos.HasRangeOverlapWith(tt.start); got != tt.wantMatch {
				t.Errorf("HasRangeOverlapWith() = %v, want %v", got, tt.wantMatch)
			}
		})
	}
}
