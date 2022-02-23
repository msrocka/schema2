package main

import (
	"bytes"
	"strings"
	"unicode"
)

// Converts the given identifier from camelCase to snake_case.
func toSnakeCase(identifier string) string {
	var buff bytes.Buffer
	for i, char := range identifier {
		if i > 0 && unicode.IsUpper(char) {
			buff.WriteRune('_')
		}
		buff.WriteRune(unicode.ToLower(char))
	}
	return buff.String()
}

// Formats the given comment to have a line length of max. 80 characters.
func formatComment(comment string, indent string) string {
	if strings.TrimSpace(comment) == "" {
		return ""
	}

	// split words by whitespaces
	var words []string
	var word bytes.Buffer
	for _, char := range comment {
		if unicode.IsSpace(char) {
			if word.Len() > 0 {
				words = append(words, word.String())
			}
			word.Reset()
			continue
		}
		word.WriteRune(char)
	}
	if word.Len() > 0 {
		words = append(words, word.String())
	}
	if len(words) == 0 {
		return ""
	}

	// format the comment
	text := ""
	line := indent + "//"
	for _, word := range words {
		nextLine := line + " " + word
		if len(nextLine) < 80 {
			line = nextLine
		} else {
			text += line + "\n"
			line = indent + "// " + word
		}
	}
	if line != indent+"// " {
		text += line + "\n"
	}
	return text
}
