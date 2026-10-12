package main

import (
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/triole/logseal"
)

// find scans directory tree and populates ps.List with matched paths
// Returns error if directory traversal fails
// depthToDetect: target depth to collect, or -1 for all depths
func (ps *tPaths) find(depthToDetect int) error {
	var pathList []tPath
	var err error = filepath.Walk(ps.RootDir, visitPaths(&pathList, ps.RxMatcher))
	if err != nil {
		lg.Error("unable to detect files", logseal.F{"error": err})
		return err
	}

	for _, p := range pathList {
		info := p.Info
		ext := strings.TrimPrefix(path.Ext(p.Path), ".")
		pi := tPath{
			Path:      p.Path,
			Folder:    trimSuf(parentDir(p.Path)),
			Name:      strings.TrimSuffix(info.Name(), "."+ext),
			Extension: ext,
			IsFolder:  info.IsDir(),
			Depth:     depth(p.Path),
		}
		if depthToDetect == pi.Depth || depthToDetect < 0 {
			ps.List = append(ps.List, pi)
		}
		if pi.IsFolder && ps.MaxDepth < pi.Depth {
			ps.MaxDepth = pi.Depth
		}
	}
	
	sort.Sort(tPathList(ps.List))
	lg.Info("got find results", logseal.F{
		"depth":    depthToDetect,
		"folder":   ps.RootDir,
		"matcher":  ps.RxMatcher,
		"no_files": len(ps.List),
	})
	return nil
}

func parentDir(p string) string {
	d := filepath.Dir(p)
	if d == "." {
		return ""
	}
	return d
}

func visitPaths(paths *[]tPath, rx string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if rxMatch(rx, path) {
			*paths = append(*paths, tPath{
				Path: path,
				Info: info,
			})
		}
		return nil
	}
}

func depth(pth string) int {
	return strings.Count(strings.TrimSuffix(pth, "/"), string(os.PathSeparator))
}
