# Fnn ![example workflow](https://github.com/triole/fnn/actions/workflows/build.yaml/badge.svg)

<!-- toc -->

- [Help](#help)
- [Disclaimer](#disclaimer)

<!-- /toc -->

File name normalizer that renames files and folder removing special characters. Contains a mighty chunk of extra opinionated renaming rules.

## Help

```go mdox-exec="r -h"

rename files removing special characters, leave only those that do not need to
be escaped in shell

Flags:
  -h, --help                      Show context-sensitive help.
  -f, --folder="/home/ole/rolling/golang/projects/fnn/src"
                                  root folder in which to run the renamer,
                                  default is current folder
  -m, --matcher=".*"              regex matcher, process only paths that fit the
                                  expression
  -r, --recursive                 recurse into sub directories, default is none
  -l, --log-file="/dev/stdout"    log file
  -e, --log-level="info"          log level
  -c, --log-no-colors             disable output colours, print plain text
  -j, --log-json                  enable json log, instead of text one
  -n, --dry-run                   debug mode
  -V, --version-flag              display version
```

## Disclaimer

Warning. Use this software at your own risk. I may not be hold responsible for any data loss, starving your kittens or losing the bling bling powerpoint presentation you made to impress human resources with the efficiency of your employee's performance.
