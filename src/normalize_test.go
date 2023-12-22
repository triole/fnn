package main

import (
	"fmt"
	"testing"
)

func TestNormalize(t *testing.T) {
	rs := initReplacerSchemes()
	validateNormalize("HELLO there", "ZiP", "hello_there", "zip", rs, t)
	validateNormalize("hello  there", "ZiP", "hello_there", "zip", rs, t)
	validateNormalize("helläöü", "txt", "hellaeoeue", "txt", rs, t)
	validateNormalize("hello's", "txt", "hello_s", "txt", rs, t)
	validateNormalize("àáâãå", "txt", "aaaaa", "txt", rs, t)
	validateNormalize("èéêë", "txt", "eeee", "txt", rs, t)
	validateNormalize("íìîĩï", "txt", "iiiii", "txt", rs, t)
	validateNormalize("òóôõ", "txt", "oooo", "txt", rs, t)
	validateNormalize("ùúû", "txt", "uuu", "txt", rs, t)
	validateNormalize("ç", "txt", "c", "txt", rs, t)
	validateNormalize("hello there   111", "txt", "hello_there111", "txt", rs, t)

}

func validateNormalize(inpName, inpExt, expName, expExt string, rs tReplacerSchemes, t *testing.T) {
	inp := newPath(inpName, inpExt)
	exp := newPath(expName, expExt)
	res := normalize(inp, rs)
	if res.Extension != exp.Extension || res.Name != exp.Name {
		t.Errorf(
			"\ntest normalize path validation failed"+
				"\ninput       %s\nresult      %s\nexpectation %s",
			printPath(inp), printPath(res), printPath(exp),
		)
	}
}

func printPath(inp tPath) string {
	return fmt.Sprintf("%s/%s.%s", inp.Path, inp.Name, inp.Extension)
}

func newPath(name, ext string) (r tPath) {
	r.Name = name
	r.Extension = ext
	return
}
