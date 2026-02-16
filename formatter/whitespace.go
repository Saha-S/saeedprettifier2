package formatter

import (
	"regexp"
	"strings"
)

// WhitespaceFormatter handles whitespace-related formatting
type WhitespaceFormatter interface {
	ConvertControlChars(text string) string
	CollapseBlankLines(text string) string
	TrimExcessiveWhitespace(text string) string
}

type WhitespaceProcessor struct{}

func NewWhitespaceFormatter() *WhitespaceProcessor {
	return &WhitespaceProcessor{}
}

func (f *WhitespaceProcessor) ConvertControlChars(text string) string {
	var result strings.Builder
	for _, r := range text {
		switch r {
		case '\v', '\f', '\r':
			result.WriteRune('\n')
		default:
			result.WriteRune(r)
		}
	}
	return result.String()
}

func (f *WhitespaceProcessor) CollapseBlankLines(text string) string {
	lines := strings.Split(text, "\n")
	result := make([]string, 0, len(lines))

	blankCount := 0
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			blankCount++
			if blankCount <= 1 {
				result = append(result, "")
			}
		} else {
			blankCount = 0
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}

func (f *WhitespaceProcessor) TrimExcessiveWhitespace(text string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		re := regexp.MustCompile(`[ \t]+`)
		lines[i] = re.ReplaceAllString(line, " ")
		lines[i] = strings.TrimSpace(lines[i])
	}
	return strings.Join(lines, "\n")
}
