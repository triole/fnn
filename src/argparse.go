package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/alecthomas/kong"
)

var (
	// BUILDTAGS are injected ld flags during build
	BUILDTAGS      string
	appName        = "file name normalizer"
	appDescription = "rename files removing special characters, leave only those that do not need to be escaped in shell"
	appMainversion = "0.1"
)

var CLI struct {
	RootDir     string `help:"root directory in which to run the renamer, default is current folder" short:"r" default:"${curdir}"`
	RxMatcher   string `help:"regex matcher, process only paths that fit the expression" short:"m" default:".*"`
	LogFile     string `help:"log file" short:"l" default:"/dev/stdout"`
	LogLevel    string `help:"log level" short:"e" default:"info" enum:"trace,debug,info,error"`
	LogNoColors bool   `help:"disable output colours, print plain text" short:"c"`
	LogJSON     bool   `help:"enable json log, instead of text one" short:"j"`
	DryRun      bool   `help:"debug mode" short:"n"`
	VersionFlag bool   `help:"display version" short:"V"`
}

func parseArgs() {
	curdir, _ := os.Getwd()
	ctx := kong.Parse(&CLI,
		kong.Name(appName),
		kong.Description(appDescription),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact:      true,
			Summary:      true,
			NoAppSummary: true,
			FlagsLast:    false,
		}),
		kong.Vars{
			"curdir": curdir,
			"config": path.Join(getBindir(), appName+".toml"),
		},
	)
	_ = ctx.Run()

	if CLI.VersionFlag == true {
		printBuildTags(BUILDTAGS)
		os.Exit(0)
	}
	// ctx.FatalIfErrorf(err)
}

type tPrinter []tPrinterEl
type tPrinterEl struct {
	Key string
	Val string
}

func printBuildTags(buildtags string) {
	regexp, _ := regexp.Compile(`({|}|,)`)
	s := regexp.ReplaceAllString(buildtags, "\n")
	s = strings.Replace(s, "_subversion: ", "version: "+appMainversion+".", -1)
	fmt.Printf("\n%s\n%s\n\n", appName, appDescription)
	arr := strings.Split(s, "\n")
	var pr tPrinter
	var maxlen int
	for _, line := range arr {
		if strings.Contains(line, ":") {
			l := strings.Split(line, ":")
			if len(l[0]) > maxlen {
				maxlen = len(l[0])
			}
			pr = append(pr, tPrinterEl{l[0], strings.Join(l[1:], ":")})
		}
	}
	for _, el := range pr {
		fmt.Printf("%"+strconv.Itoa(maxlen)+"s\t%s\n", el.Key, el.Val)
	}
	fmt.Printf("\n")
}

func getBindir() (s string) {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	s = filepath.Dir(ex)
	return
}
