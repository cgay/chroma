package html

import (
	"errors"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

func TestCompressStyle(t *testing.T) {
	style := "color: #888888; background-color: #faffff"
	actual := compressStyle(style)
	expected := "color:#888;background-color:#faffff"
	assert.Equal(t, expected, actual)
}

func BenchmarkHTMLFormatter(b *testing.B) {
	formatter := New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it, err := lexers.Go.Tokenise(nil, "package main\nfunc main()\n{\nprintln(`hello world`)\n}\n")
		assert.NoError(b, err)
		err = formatter.Format(ioutil.Discard, styles.Fallback, it)
		assert.NoError(b, err)
	}
}

func TestSplitTokensIntoLines(t *testing.T) {
	in := []*chroma.Token{
		{Value: "hello", Type: chroma.NameKeyword},
		{Value: " world\nwhat?\n", Type: chroma.NameKeyword},
	}
	expected := [][]*chroma.Token{
		{
			{Type: chroma.NameKeyword, Value: "hello"},
			{Type: chroma.NameKeyword, Value: " world\n"},
		},
		{
			{Type: chroma.NameKeyword, Value: "what?\n"},
		},
		{
			{Type: chroma.NameKeyword},
		},
	}
	actual := splitTokensIntoLines(in)
	assert.Equal(t, expected, actual)
}

func TestIteratorPanicRecovery(t *testing.T) {
	it := func() *chroma.Token {
		panic(errors.New("bad"))
	}
	err := New().Format(ioutil.Discard, styles.Fallback, it)
	assert.Error(t, err)
}

func TestFormatterStyleToCSS(t *testing.T) {
	builder := styles.Get("github").Builder()
	builder.Add(chroma.LineHighlight, "bg:#ffffcc")
	builder.Add(chroma.LineNumbers, "bold")
	style, err := builder.Build()
	if err != nil {
		t.Error(err)
	}
	formatter := New(WithClasses())
	css := formatter.styleToCSS(style)
	for _, s := range css {
		if strings.HasPrefix(strings.TrimSpace(s), ";") {
			t.Errorf("rule starts with semicolon - expected valid css rule without semicolon: %v", s)
		}
	}
}
