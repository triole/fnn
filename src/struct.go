package main

import (
	"fmt"
	"path"
)

type tPath struct {
	Path      string
	Folder    string
	Name      string
	Extension string
	Depth     int
	IsFolder  bool
}

type tPathList []tPath

type tPaths struct {
	RootDir   string
	RxMatcher string
	List      tPathList
	MaxDepth  int
}

func newPaths(rootDir, rxMatcher string) (ps tPaths) {
	ps.RootDir = rootDir
	ps.RxMatcher = rxMatcher
	return
}

func (pl tPathList) Len() int {
	return len(pl)
}

func (pl tPathList) Less(i, j int) bool {
	a := toStr(pl[i].IsFolder) + pl[i].Path
	b := toStr(pl[j].IsFolder) + pl[j].Path
	return a < b
}

func (pl tPathList) Swap(i, j int) {
	pl[i], pl[j] = pl[j], pl[i]
}

func toStr(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

func pathStr(pth tPath) string {
	p := path.Join(pth.Folder, pth.Name)
	if pth.Extension == "" {
		return p
	}
	return fmt.Sprintf("%s.%s", p, pth.Extension)
}
