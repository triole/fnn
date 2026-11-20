package main

import (
	"errors"
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

	// Regex patterns for build tag parsing
	buildTagPattern *regexp.Regexp
)

type tCLI struct {
	Folder      string  `help:"root folder in which to run the renamer, default is current folder" short:"f" default:"${curdir}"`
	Matcher     string  `help:"regex matcher, process only paths that fit the expression" short:"m" default:".*"`
	Recursive   bool    `help:"recurse into sub directories, default is none" short:"r"`
	Sanitize    *string `help:"sanitize string and print result" short:"s" placeholder:"STRING"`
	LogFile     string  `help:"log file" default:"/dev/stdout"`
	LogLevel    string  `help:"log level" default:"info" enum:"trace,debug,info,error,fatal"`
	LogNoColors bool    `help:"disable output colours, print plain text"`
	LogJSON     bool    `help:"enable json log, instead of text one"`
	DryRun      bool    `help:"debug mode" short:"n"`
	VersionFlag bool    `help:"display version" short:"V"`
}

var CLI tCLI

func init() {
	// Pre-compile regex patterns at startup
	var err error
	buildTagPattern, err = regexp.Compile(`({|}|,)`)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal: regex compilation failed: %v\n", err)
		os.Exit(1)
	}
}

func parseArgs() {
	parser, err := newArgParser(&CLI)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: cannot initialize argument parser: %v\n", err)
		os.Exit(1)
	}
	ctx, err := parser.Parse(os.Args[1:])
	if err != nil {
		// Diagnostics go to stderr only. Kong's usage-on-error prints the
		// full help to stdout, which a shell capture such as
		// `value=$(fnn -s "$arg")` would silently swallow.
		fmt.Fprintf(os.Stderr, "%s: error: %v\n", appName, err)
		exitCode := 1
		var parseErr *kong.ParseError
		if errors.As(err, &parseErr) {
			exitCode = parseErr.ExitCode()
		}
		os.Exit(exitCode)
	}
	_ = ctx.Run()

	if CLI.VersionFlag {
		printBuildTags(BUILDTAGS)
		os.Exit(0)
	}
}

func newArgParser(cli *tCLI) (*kong.Kong, error) {
	curdir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return kong.New(cli,
		kong.Name(appName),
		kong.Description(appDescription),
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
}

type tPrinter []tPrinterEl
type tPrinterEl struct {
	Key string
	Val string
}

func printBuildTags(buildtags string) {
	s := buildTagPattern.ReplaceAllString(buildtags, "\n")
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
