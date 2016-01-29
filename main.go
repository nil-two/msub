package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

var (
	version = "0.2.1"

	flagset     = flag.NewFlagSet("msub", flag.ContinueOnError)
	useBoundary = flagset.Bool("boundary", false, "")
	isHelp      = flagset.Bool("help", false, "")
	isVersion   = flagset.Bool("version", false, "")
)

func init() {
	flagset.SetOutput(ioutil.Discard)
	flagset.BoolVar(useBoundary, "b", false, "")
	flagset.BoolVar(isHelp, "h", false, "")
	flagset.BoolVar(isVersion, "v", false, "")
}

func usage() {
	os.Stderr.WriteString(`
Usage: msub [OPTION]... FROM TO [FILE]...
Substitute multiple words at once
by FROM and TO patterns.

Options:
  -b, --boundary    use word boundary in matcher
  -h, --help        show this help message
  -v, --version     output version information and exit

Syntax:
  pattern = group {"/" group}
  group   = branch {"," branch}
  branch  = {letter | "\/" | "\,"}

Examples:
  msub true,false false,true ./file
  msub dog,cat/s cat,dog/s ~/Document/questionnaire
`[1:])
}

func printVersion() {
	fmt.Fprintln(os.Stderr, version)
}

func printError(err interface{}) {
	fmt.Fprintln(os.Stderr, "msub:", err)
}

func do(r *Replacer, src io.Reader) error {
	b := bufio.NewScanner(src)
	for b.Scan() {
		fmt.Println(r.ReplaceAll(b.Text()))
	}
	return b.Err()
}

func _main() int {
	if err := flagset.Parse(os.Args[1:]); err != nil {
		printError(err)
		return 2
	}
	switch {
	case *isHelp:
		usage()
		return 0
	case *isVersion:
		printVersion()
		return 0
	}

	switch flagset.NArg() {
	case 0:
		printError("no specify FROM and TO")
		return 2
	case 1:
		printError("no specify TO")
		return 2
	}
	from, to := flagset.Arg(0), flagset.Arg(1)

	r, err := NewReplacer(from, to, *useBoundary)
	if err != nil {
		printError(err)
		return 2
	}

	if flagset.NArg() < 3 {
		if err = do(r, os.Stdin); err != nil {
			printError(err)
			return 1
		}
		return 0
	}

	var srcls []io.Reader
	for _, name := range flagset.Args()[2:] {
		f, err := os.Open(name)
		if err != nil {
			printError(err)
			return 1
		}
		defer f.Close()
		srcls = append(srcls, f)
	}
	if err = do(r, io.MultiReader(srcls...)); err != nil {
		printError(err)
		return 1
	}
	return 0
}

func main() {
	e := _main()
	os.Exit(e)
}
