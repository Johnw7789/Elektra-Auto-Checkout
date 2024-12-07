package shr

import (
	"strings"
)

func Parse(value string, a string, b string) string {
	posFirst := strings.Index(value, a)
	if posFirst == -1 {
		return ""
	}

	posLast := strings.Index(value, b)
	if posLast == -1 {
		return ""
	}

	posFirstAdjusted := posFirst + len(a)
	if posFirstAdjusted >= posLast {
		return ""
	}

	return value[posFirstAdjusted:posLast]
}

func ParseV2(str string, start string, end string) (result string) {
	s := strings.Index(str, start)
	if s == -1 {
		return
	}

	s += len(start)

	e := strings.Index(str[s:], end)
	if e == -1 {
		return
	}

	e += s + e - 1

	return str[s:e]
}

func ParseV3(str string, start string, end string) string {
	var match []byte
	index := strings.Index(str, start)

	if index == -1 {
		return string(match)
	}

	index += len(start)

	for {
		char := str[index]

		if strings.HasPrefix(str[index:index+len(match)], end) {
			break
		}

		match = append(match, char)
		index++
	}

	return string(match)
}

func Reverse(s string) string {
	rs := []rune(s)
	
	for i, j := 0, len(rs)-1; i < j; i, j = i+1, j-1 {
		rs[i], rs[j] = rs[j], rs[i]
	}

	return string(rs)
}
