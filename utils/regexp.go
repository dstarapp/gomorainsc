package utils

import (
	"regexp"
)

func RegSub(text string, reg string) string {
	regx := regexp.MustCompile(reg)
	elems := regx.FindAllStringSubmatch(text, -1)
	if len(elems) > 0 && len(elems[0]) == 2 {
		return elems[0][1]
	}
	return ""
}
