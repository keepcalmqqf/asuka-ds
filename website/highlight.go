package main

import (
	"strings"
	"unicode"

	"github.com/egoist/asuka/colors"
	"github.com/egoist/asuka/fonts"
	"github.com/egoist/asuka/ui"
)

// tokenKind is one of the few roles a documentation snippet needs. The
// palette is deliberately small: a reader is looking for strings, names
// and comments, not a full language grammar.
type tokenKind uint8

const (
	tokenPlain tokenKind = iota
	tokenComment
	tokenString
	tokenKeyword
	tokenType
	tokenNumber
	tokenFlag
	tokenKey
)

type codeToken struct {
	text string
	kind tokenKind
}

var goKeywords = map[string]bool{
	"break": true, "case": true, "chan": true, "const": true, "continue": true,
	"default": true, "defer": true, "else": true, "fallthrough": true, "for": true,
	"func": true, "go": true, "goto": true, "if": true, "import": true,
	"interface": true, "map": true, "package": true, "range": true, "return": true,
	"select": true, "struct": true, "switch": true, "type": true, "var": true,
}

var goTypes = map[string]bool{
	"bool": true, "byte": true, "complex64": true, "complex128": true,
	"error": true, "float32": true, "float64": true, "int": true, "int8": true,
	"int16": true, "int32": true, "int64": true, "rune": true, "string": true,
	"uint": true, "uint8": true, "uint16": true, "uint32": true, "uint64": true,
	"uintptr": true, "any": true,
}

func highlight(language, src string) ui.Rich {
	tokens := tokenize(language, src)
	spans := make(ui.Rich, 0, len(tokens))
	for _, token := range tokens {
		if token.text == "" {
			continue
		}
		span := ui.Span{Text: token.text}
		switch token.kind {
		case tokenComment:
			span.Color = colors.SecondaryLabel
			span.Italic = true
		case tokenString:
			span.Color = colors.Green
		case tokenKeyword:
			span.Color = colors.Purple
			span.Font.Weight = fonts.Bold
		case tokenType:
			span.Color = colors.Blue
		case tokenNumber:
			span.Color = colors.Orange
		case tokenFlag:
			span.Color = colors.Pink
		case tokenKey:
			span.Color = colors.Blue
		}
		spans = append(spans, span)
	}
	return spans
}

func tokenize(language, src string) []codeToken {
	switch language {
	case "JSON":
		return tokenizeJSON(src)
	case "Shell":
		return tokenizeShell(src)
	default:
		return tokenizeGo(src)
	}
}

func tokenizeGo(src string) []codeToken {
	var out []codeToken
	for pos := 0; pos < len(src); {
		ch := src[pos]
		switch {
		case ch == '/' && pos+1 < len(src) && src[pos+1] == '/':
			end := strings.IndexByte(src[pos:], '\n')
			if end < 0 {
				end = len(src) - pos
			}
			out = append(out, codeToken{src[pos : pos+end], tokenComment})
			pos += end
		case ch == '/' && pos+1 < len(src) && src[pos+1] == '*':
			end := strings.Index(src[pos+2:], "*/")
			if end < 0 {
				end = len(src) - pos
			} else {
				end += 4
			}
			out = append(out, codeToken{src[pos : pos+end], tokenComment})
			pos += end
		case ch == '"' || ch == '\'' || ch == '`':
			end := quotedEnd(src, pos, ch)
			out = append(out, codeToken{src[pos:end], tokenString})
			pos = end
		case isDigit(ch):
			end := scanWhile(src, pos, func(c byte) bool { return isDigit(c) || c == '.' || c == 'x' || c == '_' })
			out = append(out, codeToken{src[pos:end], tokenNumber})
			pos = end
		case isIdentStart(ch):
			end := scanWhile(src, pos, isIdentPart)
			word := src[pos:end]
			kind := tokenPlain
			if goKeywords[word] {
				kind = tokenKeyword
			} else if goTypes[word] {
				kind = tokenType
			} else if word == "true" || word == "false" || word == "nil" {
				kind = tokenNumber
			} else if isGoName(word, src, pos, end) {
				kind = tokenType
			}
			out = append(out, codeToken{word, kind})
			pos = end
		default:
			end := scanUntilTokenStart(src, pos, isGoTokenStart)
			out = append(out, codeToken{src[pos:end], tokenPlain})
			pos = end
		}
	}
	return out
}

func tokenizeShell(src string) []codeToken {
	var out []codeToken
	for _, line := range strings.SplitAfter(src, "\n") {
		content := strings.TrimSuffix(line, "\n")
		newline := strings.TrimPrefix(line, content)
		commentAt := strings.Index(content, " #")
		code := content
		if commentAt >= 0 {
			code = content[:commentAt+1]
		}
		words := strings.Fields(code)
		cursor := 0
		first := true
		for _, word := range words {
			at := strings.Index(code[cursor:], word)
			if at < 0 {
				continue
			}
			start := cursor + at
			out = append(out, codeToken{code[cursor:start], tokenPlain})
			kind := tokenPlain
			if first {
				kind = tokenKeyword
			} else if strings.HasPrefix(word, "-") {
				kind = tokenFlag
			}
			out = append(out, codeToken{word, kind})
			cursor = start + len(word)
			first = false
		}
		out = append(out, codeToken{code[cursor:], tokenPlain})
		if commentAt >= 0 {
			out = append(out, codeToken{content[commentAt+1:], tokenComment})
		}
		out = append(out, codeToken{newline, tokenPlain})
	}
	return out
}

func tokenizeJSON(src string) []codeToken {
	var out []codeToken
	for pos := 0; pos < len(src); {
		ch := src[pos]
		switch {
		case ch == '"':
			end := quotedEnd(src, pos, '"')
			lookahead := strings.TrimLeft(src[end:], " \t")
			kind := tokenString
			if strings.HasPrefix(lookahead, ":") {
				kind = tokenKey
			}
			out = append(out, codeToken{src[pos:end], kind})
			pos = end
		case isDigit(ch) || ch == '-':
			end := scanWhile(src, pos, func(c byte) bool { return isDigit(c) || c == '.' || c == '-' || c == '+' })
			out = append(out, codeToken{src[pos:end], tokenNumber})
			pos = end
		case isIdentStart(ch):
			end := scanWhile(src, pos, isIdentPart)
			out = append(out, codeToken{src[pos:end], tokenKeyword})
			pos = end
		default:
			end := scanUntilTokenStart(src, pos, isJSONTokenStart)
			out = append(out, codeToken{src[pos:end], tokenPlain})
			pos = end
		}
	}
	return out
}

func quotedEnd(src string, pos int, quote byte) int {
	for i := pos + 1; i < len(src); i++ {
		if src[i] == '\\' && quote != '`' {
			i++
			continue
		}
		if src[i] == quote {
			return i + 1
		}
	}
	return len(src)
}

// isGoName colors the identifiers a reader scans for in a short snippet:
// declarations, package qualifiers, types, fields and calls. It is
// intentionally heuristic; comments and strings have already been taken
// out above.
func isGoName(word string, src string, start, end int) bool {
	if word == "" {
		return false
	}
	before := strings.TrimRight(src[:start], " \t")
	if strings.HasSuffix(before, "func ") || strings.HasSuffix(before, "type ") ||
		strings.HasSuffix(before, "var ") || strings.HasSuffix(before, ".") ||
		strings.HasSuffix(before, ":=") || strings.HasSuffix(before, "=") {
		return true
	}
	after := strings.TrimLeft(src[end:], " \t")
	if strings.HasPrefix(after, "(") || strings.HasPrefix(after, "{") ||
		strings.HasPrefix(after, ":") || strings.HasPrefix(after, ".") ||
		strings.HasPrefix(after, ":=") {
		return true
	}
	return word[0] >= 'A' && word[0] <= 'Z'
}

func scanWhile(src string, pos int, accept func(byte) bool) int {
	for pos < len(src) && accept(src[pos]) {
		pos++
	}
	return pos
}

func scanUntilTokenStart(src string, pos int, start func(byte) bool) int {
	beginning := pos
	for pos < len(src) {
		if pos > beginning && start(src[pos]) {
			break
		}
		pos++
	}
	if pos == beginning {
		pos++
	}
	return pos
}

func isGoTokenStart(ch byte) bool {
	return isIdentStart(ch) || isDigit(ch) || ch == '"' || ch == '\'' || ch == '`' || ch == '/' || unicode.IsSpace(rune(ch))
}

func isJSONTokenStart(ch byte) bool {
	return isIdentStart(ch) || isDigit(ch) || ch == '-' || ch == '"'
}

func isDigit(ch byte) bool { return ch >= '0' && ch <= '9' }

func isIdentStart(ch byte) bool {
	return ch == '_' || ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch >= 0x80
}

func isIdentPart(ch byte) bool { return isIdentStart(ch) || isDigit(ch) }
