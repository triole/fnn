package main

import "github.com/triole/logseal"

func (ps *tPaths) normalizeAll() (r string) {
	for _, pth := range ps.List {
		npth := normalize(pth)
		lg.Info("normalize path", logseal.F{
			"old":       pth.Path,
			"new":       npth.Path,
			"is_folder": pth.IsFolder,
		})
	}
	return
}

func normalize(pth tPath) (r tPath) {
	r = pth
	return r
}
