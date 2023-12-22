package main

import (
	"github.com/triole/logseal"
)

var (
	lg = logseal.Init("debug", nil, true, false)
)

func main() {
	parseArgs()
	lg = logseal.Init(CLI.LogLevel, CLI.LogFile, CLI.LogNoColors, CLI.LogJSON)
	folder := trimSuf(CLI.Folder)

	lg.Info("run "+appName, logseal.F{"folder": folder})

	if CLI.Recursive {
		// detect all, unlimited depth
		ps := newPaths(folder, CLI.Matcher)
		ps.find(-1)

		// recurse from high depth to lower
		for i := ps.MaxDepth; i > depth(folder); i-- {
			np := newPaths(folder, CLI.Matcher)
			np.find(i)
			np.normalizeAll()
		}
	} else {
		ps := newPaths(folder, CLI.Matcher)
		ps.find(depth(folder) + 1)
		ps.normalizeAll()
	}

	if CLI.DryRun {
		lg.Info("dry run, no files were renamed")
	}

}
