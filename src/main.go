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

	lg.Info("run " + appName + " in folder " + CLI.RootDir)

	// detect all, unlimited depth
	ps := newPaths(CLI.RootDir, CLI.RxMatcher)
	ps.find(-1)

	// recurse from high depth to lower
	for i := ps.MaxDepth; i >= depth(CLI.RootDir); i-- {
		np := newPaths(CLI.RootDir, CLI.RxMatcher)
		np.find(i)
		np.normalizeAll()
	}

}
