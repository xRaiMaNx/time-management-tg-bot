package tg

import "unicode"

func CollapseSpaces(input string) string {
	ansRunes := []rune{}
	isPrevSpace := false
	for _, rune := range input {
		isCurSpace := unicode.IsSpace(rune)
		if isCurSpace && isPrevSpace {
			continue
		} else if isCurSpace && !isPrevSpace {
			ansRunes = append(ansRunes, ' ')
			isPrevSpace = true
		} else {
			ansRunes = append(ansRunes, rune)
			isPrevSpace = false
		}
	}
	return string(ansRunes)
}
