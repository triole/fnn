package main

import (
	"regexp"
	"strings"
)

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func rxMatch(compiledRx *regexp.Regexp, str string) bool {
	return compiledRx.MatchString(str)
}

func rxReplaceAll(basestring string, regex *regexp.Regexp, newstring string) (r string) {
	r = regex.ReplaceAllString(basestring, newstring)
	return
}

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
