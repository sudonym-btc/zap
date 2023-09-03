package tip

import (
	"regexp"
)

func ExtractEmails(text string) []string {
	re := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	matches := re.FindAllString(text, -1)
	return matches
}
