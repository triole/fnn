package main

import (
	"regexp"
	"strings"

	"github.com/triole/logseal"
)

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func rxMatch(rx string, str string) bool {
	re, err := regexp.Compile(rx)
	if err != nil {
		lg.Error("invalid regex pattern", logseal.F{"error": err, "pattern": rx})
		return false
	}
	return re.MatchString(str)
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
