package main

import (
	"os"
	"strings"

	"github.com/triole/logseal"
	"gopkg.in/yaml.v2"

	_ "embed"
)

//go:embed replacer_schemes.yaml
var replacerSchemesStr string

var replacerSchemes tReplacerSchemes

type tReplacerSchemes []tReplacerScheme
type tReplacerScheme struct {
	Rx string
	Nu string
}

func (ps *tPaths) normalizeAll() (r string) {
	replacerSchemes = initReplacerSchemes()
	for _, pth := range ps.List {
		npth := normalize(pth, replacerSchemes)
		rename(pth, npth)
	}
	return
}

func normalize(pth tPath, replacerSchemes tReplacerSchemes) (r tPath) {
	r = pth
	r.Name = strings.ToLower(r.Name)
	r.Extension = strings.ToLower(r.Extension)
	for _, rs := range replacerSchemes {
		r.Name = rxReplaceAll(r.Name, rs.Rx, rs.Nu)
	}
	return r
}

func initReplacerSchemes() (r tReplacerSchemes) {
	err := yaml.Unmarshal([]byte(replacerSchemesStr), &r)
	if err != nil {
		lg.Fatal("error: %v", err)
	}
	return
}

func rename(pth, npth tPath) {
	oldPath := pathStr(pth)
	newPath := pathStr(npth)
	lg.Info("normalize name", logseal.F{
		"old":       oldPath,
		"new":       newPath,
		"is_folder": pth.IsFolder,
	})
	if !CLI.DryRun {
		os.Rename(oldPath, newPath)
	}
}
