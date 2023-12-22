package main

import "regexp"

func rxMatch(rx string, str string) (b bool) {
	re, _ := regexp.Compile(rx)
	b = re.MatchString(str)
	return
}

func rxFind(rx string, str string) (r string) {
	temp, _ := regexp.Compile(rx)
	r = temp.FindString(str)
	return
}

func rxReplaceAll(basestring, regex, newstring string) (r string) {
	rx := regexp.MustCompile(regex)
	r = rx.ReplaceAllString(basestring, newstring)
	return
}
