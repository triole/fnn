package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
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
	Rx *regexp.Regexp  // Precompiled regex pattern
	Nu string          // Replacement string
}

// tReplacerSchemeConfig is used for YAML unmarshaling
type tReplacerSchemeConfig struct {
	Rx string `yaml:"rx"`
	Nu string `yaml:"nu"`
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

func initReplacerSchemes() tReplacerSchemes {
	var configSchemes []tReplacerSchemeConfig
	err := yaml.Unmarshal([]byte(replacerSchemesStr), &configSchemes)
	if err != nil {
		fmt.Printf("Fatal error: failed to parse replacer schemes YAML: %v\n", err)
		os.Exit(1)
	}

	// Convert config to precompiled schemes
	schemes := make(tReplacerSchemes, len(configSchemes))
	for i, cfg := range configSchemes {
		rx, err := regexp.Compile(cfg.Rx)
		if err != nil {
			fmt.Printf("Fatal error: invalid regex in scheme %d: %v (pattern: %s)\n", i+1, err, cfg.Rx)
			os.Exit(1)
		}
		schemes[i] = tReplacerScheme{
			Rx: rx,
			Nu: cfg.Nu,
		}
	}
	return schemes
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
