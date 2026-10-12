package main

import (
	"github.com/triole/logseal"
)

var (
	lg logseal.Logseal
)

func main() {
	parseArgs()
	lg = logseal.Init(CLI.LogLevel, CLI.LogFile, CLI.LogNoColors, CLI.LogJSON)
	folder := trimSuf(CLI.Folder)

	lg.Info("run "+appName, logseal.F{"folder": folder})

	if CLI.Recursive {
		// scan tree once to get ALL paths
		ops := newPaths(folder, CLI.Matcher)
		if err := ops.find(-1); err != nil {
			lg.Fatal("directory scan failed", logseal.F{"error": err})
		}

		// group paths by depth
		depthGroups := make(map[int]tPathList)
		for _, pth := range ops.List {
			if _, exists := depthGroups[pth.Depth]; !exists {
				depthGroups[pth.Depth] = tPathList{}
			}
			depthGroups[pth.Depth] = append(depthGroups[pth.Depth], pth)
		}

		// process from deepest to shallowest to avoid path conflicts
		for i := ops.MaxDepth; i >= depth(folder); i-- {
			if group, exists := depthGroups[i]; exists {
				tempPs := tPaths{RootDir: folder, List: group, MaxDepth: ops.MaxDepth}
				tempPs.normalizeAll()
			}
		}
	} else {
		pathDepth := depth(folder) + 1
		ps := newPaths(folder, CLI.Matcher)
		if err := ps.find(pathDepth); err != nil {
			lg.Fatal("directory scan failed", logseal.F{"error": err})
		}
		ps.normalizeAll()
	}

	if CLI.DryRun {
		lg.Info("dry run, no files were renamed")
	}
}
