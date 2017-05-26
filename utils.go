package main

import (
	"strconv"
	"strings"
)

func conjoin(items []int) string {
	if len(items) == 0 {
		return ""
	}
	if len(items) == 1 {
		return strconv.Itoa(items[0])
	}
	if len(items) == 2 { // "a and b" not "a, and b"
		return strconv.Itoa(items[0]) + " " + "," + " " + strconv.Itoa(items[1])
	}

	sep := ", "
	pieces := []string{strconv.Itoa(items[0])}
	for _, item := range items[1 : len(items)-1] {
		pieces = append(pieces, sep, strconv.Itoa(item))
	}
	pieces = append(pieces, sep, strconv.Itoa(items[len(items)-1]))

	return strings.Join(pieces, "")
}

