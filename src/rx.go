package main

import (
	"regexp"
	"strings"
)

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func rxMatch(rx string, str string) (b bool) {
	re, _ := regexp.Compile(rx)
	b = re.MatchString(str)
	return
}

// func rxFind(rx string, str string) (r string) {
// 	temp, _ := regexp.Compile(rx)
// 	r = temp.FindString(str)
// 	return
// }

func rxReplaceAll(basestring, regex, newstring string) (r string) {
	rx := regexp.MustCompile(regex)
	r = rx.ReplaceAllString(basestring, newstring)
	return
}

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
