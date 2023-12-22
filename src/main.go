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

	lg.Info("run " + appName + " in folder " + CLI.Folder)

	// detect all, unlimited depth
	ps := newPaths(CLI.Folder, CLI.Matcher)
	ps.find(-1)

	// recurse from high depth to lower
	for i := ps.MaxDepth; i >= depth(CLI.Folder); i-- {
		np := newPaths(CLI.Folder, CLI.Matcher)
		np.find(i)
		np.normalizeAll()
	}

	if CLI.DryRun {
		lg.Info("dry run, no files were renamed")
	}

}
