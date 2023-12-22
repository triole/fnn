package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/triole/logseal"
)

func (ps *tPaths) find(depthToDetect int) {
	lg.Debug("find", logseal.F{
		"depth":   depthToDetect,
		"folder":  ps.RootDir,
		"matcher": ps.RxMatcher,
	})
	var pathList []string
	var tempPaths tPathList
	var err error = filepath.Walk(ps.RootDir, visit(&pathList, ps.RxMatcher))
	lg.IfErrFatal(
		"unable to detect files",
		logseal.F{"error": err},
	)
	for _, p := range pathList {
		pth := strings.TrimSuffix(p, "/")
		inf, err := os.Stat(pth)
		lg.IfErrError("failed to stat file", logseal.F{
			"error": err,
			"path":  pth,
		})
		if err == nil {
			ext := strings.TrimPrefix(path.Ext(pth), ".")
			pi := tPath{
				Path:      pth,
				Folder:    strings.TrimSuffix(strings.Replace(pth, inf.Name(), "", -1), "/"),
				Name:      strings.TrimSuffix(inf.Name(), "."+ext),
				Extension: ext,
				IsFolder:  inf.IsDir(),
				Depth:     depth(pth),
			}
			if depthToDetect == pi.Depth || depthToDetect < 0 {
				tempPaths = append(tempPaths, pi)
			}
			if pi.IsFolder && ps.MaxDepth < pi.Depth {
				ps.MaxDepth = pi.Depth
			}
		}
	}
	if depthToDetect < 0 {
		ps.List = tempPaths
	} else {
		for _, pth := range tempPaths {
			if depthToDetect == pth.Depth {
				ps.List = append(ps.List, pth)
			}
		}
	}
	sort.Sort(tPathList(ps.List))
	// ps.RootDir = ps.RootDir
	lg.Info("got find results", logseal.F{
		"depth":    depthToDetect,
		"folder":   ps.RootDir,
		"matcher":  ps.RxMatcher,
		"no_files": len(ps.List),
	},
	)
	return
}

func visit(files *[]string, rx string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		if rxMatch(rx, path) {
			*files = append(*files, path)
		}
		return nil
	}
}

func getAllMaxDepthFolders(paths tPaths) (r tPathList) {
	fmt.Printf("%+v\n", paths.MaxDepth)
	for _, pth := range paths.List {
		if pth.IsFolder && pth.Depth == paths.MaxDepth {
			r = append(r, pth)
		}
	}
	return
}

func depth(pth string) int {
	return strings.Count(strings.TrimSuffix(pth, "/"), string(os.PathSeparator))
}
