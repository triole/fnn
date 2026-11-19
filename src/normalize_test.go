package main

import (
	"fmt"
	"strings"
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
	validateNormalize("Hello---World", "txt", "hello_world", "txt", rs, t)
	validateNormalize("a---_b", "txt", "a_b", "txt", rs, t)
	validateNormalize("a----_----b", "txt", "a_b", "txt", rs, t)
	validateNormalize("a-_-b", "txt", "a_b", "txt", rs, t)
	validateNormalize("hello--world", "txt", "hello-world", "txt", rs, t)
	validateNormalize("hello__world", "txt", "hello_world", "txt", rs, t)
	validateNormalize("----", "txt", "", "txt", rs, t)
	validateNormalize("a---", "txt", "a", "txt", rs, t)
	validateNormalize("a-", "txt", "a", "txt", rs, t)
	validateNormalize("hello-__-", "txt", "hello", "txt", rs, t)
	validateNormalize("hello____", "txt", "hello", "txt", rs, t)
	validateNormalize("hello there   11", "txt", "hello_there11", "txt", rs, t)
	validateNormalize("hello there   222", "txt", "hello_there_222", "txt", rs, t)
	validateNormalize("hello there   3333", "txt", "hello_there_3333", "txt", rs, t)
	validateNormalize("Hello HereAndThere    3333", "txt", "hello_here_and_there_3333", "txt", rs, t)
}

func TestAccentedCharacters(t *testing.T) {
	rs := initReplacerSchemes()

	// Every character in a group must collapse to exactly one base letter.
	groups := []struct {
		in   string
		want string
	}{
		{"àáâãåāăąǎȁȃȧḁạảǡǻấầẩẫậắằẳẵặ", "a"},
		{"ḃḅḇ", "b"},
		{"çćĉċčḉ", "c"},
		{"ďḋḍḏḑḓđð", "d"},
		{"èéêëēĕėęěȅȇȩḙḛẹẻẽḕḗḝếềểễệ", "e"},
		{"ḟ", "f"},
		{"ĝğġģǧǵḡ", "g"},
		{"ĥȟḣḥḧḩḫẖħ", "h"},
		{"ìíîïĩīĭįǐȉȋḭỉịḯı", "i"},
		{"ĵǰ", "j"},
		{"ķǩḱḳḵ", "k"},
		{"ĺļľḷḻḽḹłŀ", "l"},
		{"ḿṁṃ", "m"},
		{"ñńņňǹṅṇṉṋ", "n"},
		{"òóôõōŏőơǒǫȍȏȯọỏǭȭȱṍṏṑṓốồổỗộớờởỡợ", "o"},
		{"ṕṗ", "p"},
		{"ŕŗřȑȓṙṛṟṝ", "r"},
		{"śŝşšșṡṣṥṧṩẛſ", "s"},
		{"ţťțṫṭṯṱẗŧ", "t"},
		{"ùúûũūŭůűųưǔȕȗṳṵṷụủṹṻứừửữự", "u"},
		{"ṽṿ", "v"},
		{"ŵẁẃẅẇẉẘ", "w"},
		{"ẋẍ", "x"},
		{"ýÿŷȳẏẙỳỵỷỹ", "y"},
		{"źżžẑẓẕ", "z"},
	}
	for _, g := range groups {
		validateNormalize(g.in, "txt", strings.Repeat(g.want, len([]rune(g.in))), "txt", rs, t)
	}

	// German umlauts and other Latin letters without canonical decompositions.
	validateNormalize("äǟ", "txt", "aeae", "txt", rs, t)
	validateNormalize("öȫ", "txt", "oeoe", "txt", rs, t)
	validateNormalize("üǖǘǚǜ", "txt", "ueueueueue", "txt", rs, t)
	validateNormalize("ß", "txt", "ss", "txt", rs, t)
	validateNormalize("æǣǽ", "txt", "aeaeae", "txt", rs, t)
	validateNormalize("œ", "txt", "oe", "txt", rs, t)
	validateNormalize("øǿ", "txt", "oo", "txt", rs, t)
	validateNormalize("þ", "txt", "th", "txt", rs, t)
	validateNormalize("ÄÖÜẞ", "txt", "aeoeuess", "txt", rs, t)
	validateNormalize("ÆØÞ", "txt", "aeoth", "txt", rs, t)
	validateNormalize("ĐŁŒ", "txt", "dloe", "txt", rs, t)

	// Real-world mixed-case names, including uppercase accents.
	validateNormalize("Mañana", "txt", "manana", "txt", rs, t)
	validateNormalize("São Tomé", "txt", "sao_tome", "txt", rs, t)
	validateNormalize("Élève à l'hôtel", "txt", "eleve_a_l_hotel", "txt", rs, t)
	validateNormalize("Dvořák", "txt", "dvorak", "txt", rs, t)
	validateNormalize("Łódź", "txt", "lodz", "txt", rs, t)
	validateNormalize("Þingvellir", "txt", "thingvellir", "txt", rs, t)
	validateNormalize("Tiếng Việt", "txt", "tieng_viet", "txt", rs, t)
	validateNormalize("İstanbul", "txt", "istanbul", "txt", rs, t)
	validateNormalize("Œuvre", "txt", "oeuvre", "txt", rs, t)
	validateNormalize("Æble Øl Ål", "txt", "aeble_ol_al", "txt", rs, t)
	validateNormalize("ÑANDÚ ÁÉÍÓÚ", "txt", "nandu_aeiou", "txt", rs, t)
	validateNormalize("Đại Việt", "txt", "dai_viet", "txt", rs, t)

	// Decomposed (NFD) input leaves combining marks behind after lowercasing.
	validateNormalize("man\u0303ana", "txt", "manana", "txt", rs, t)
	validateNormalize("a\u0301", "txt", "a", "txt", rs, t)
	validateNormalize("a\u0308", "txt", "ae", "txt", rs, t)
	validateNormalize("o\u0308", "txt", "oe", "txt", rs, t)
	validateNormalize("u\u0308", "txt", "ue", "txt", rs, t)
	validateNormalize("A\u0308", "txt", "ae", "txt", rs, t)
}

func validateNormalize(inpName, inpExt, expName, expExt string, rs tReplacerSchemes, t *testing.T) {
	inp := newPath(inpName, inpExt)
	exp := newPath(expName, expExt)
	res := normalize(inp, rs)
	if res.Extension != exp.Extension || res.Name != exp.Name {
		t.Errorf(
			"\ntest normalize path validation failed"+
				"\ninp %s\nres %s\nexp %s",
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
