package main

import (
	"errors"
	"os"
	"strings"

	"github.com/triole/logseal"
	yaml "gopkg.in/yaml.v2"

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
	r.Name = toSnakeCase(r.Name)
	r.Name = strings.ToLower(r.Name)
	if r.IsFolder {
		r.Name = strings.Replace(r.Name, ".", "_", -1)
		r.Name = r.Name + r.Extension
		r.Extension = ""
	}
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
	if oldPath != newPath {
		if !exists(newPath) {
			lg.Info("rename file", logseal.F{
				"old":       oldPath,
				"new":       newPath,
				"is_folder": pth.IsFolder,
			})
			if !CLI.DryRun {
				err := os.Rename(oldPath, newPath)
				if err != nil {
					lg.Error("rename file failed", logseal.F{"error": err})
				}
			}
		} else {
			lg.Warn("won't rename, target path exists", logseal.F{
				"old":       oldPath,
				"new":       newPath,
				"is_folder": pth.IsFolder,
			})
		}
	} else {
		lg.Debug("skip, file name wouldn't change", logseal.F{
			"old":       oldPath,
			"new":       newPath,
			"is_folder": pth.IsFolder,
		})
	}
}

func exists(filePath string) bool {
	_, error := os.Stat(filePath)
	return !errors.Is(error, os.ErrNotExist)
}

func trimSuf(s string) string {
	return strings.TrimSuffix(s, string(os.PathSeparator))
}
