package main

import (
	"bytes"
	"strconv"
)

func formatNumber(n int, sep rune) string {
	s := strconv.Itoa(n)

	startOffset := 0
	var buff bytes.Buffer

	if n < 0 {
		startOffset = 1
		buff.WriteByte('-')
	}

	l := len(s)

	commaIndex := 3 - ((l - startOffset) % 3)
	if commaIndex == 3 {
		commaIndex = 0
	}

	for i := startOffset; i < l; i++ {
		if commaIndex == 3 {
			buff.WriteRune(sep)
			commaIndex = 0
		}
		commaIndex++
		buff.WriteByte(s[i])
	}

	return buff.String()
}

func trendSymbol(num float64) string {
	if num < 0 {
		return "↓"
	}
	if num > 0 {
		return "↑"
	}
	return "±"
}
