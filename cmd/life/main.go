package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/418Coffee/life"
)

var seed int64
var nowrap bool
var ticks uint
var rleFile string
var width, height uint

func printUsageAndExit(err error) {
	if err != nil {
		fmt.Fprint(flag.CommandLine.Output(), err, "\n")
	}
	flag.Usage()
	os.Exit(1)
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s [options] width height\noptions:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Int64Var(&seed, "seed", time.Now().UnixMicro(), "seed for initial state")
	flag.BoolVar(&nowrap, "nowrap", false, "don't wrap field toroidally")
	flag.UintVar(&ticks, "ticks", 100, "amount of generation to run")
	flag.StringVar(&rleFile, "file", "", "load initial state from .rle file (mutually exclusive with width height arguments)")
	flag.Parse()
}

func main() {
	var l *life.Game
	var err error
	if rleFile != "" {
		if l, err = life.LoadGame(rleFile, !nowrap); err != nil {
			printUsageAndExit(err)
		}
	} else {
		args := flag.Args()
		if len(args) != 2 {
			printUsageAndExit(nil)
		}
		w, err := strconv.ParseUint(args[0], 0, 64)
		if err != nil {
			printUsageAndExit(err)
		}
		h, err := strconv.ParseUint(args[1], 0, 64)
		if err != nil {
			printUsageAndExit(err)
		}
		width, height = uint(w), uint(h)
		rand.Seed(seed)
		l = life.NewGame(width, height, !nowrap)
	}

	for i := uint(0); i < ticks; i++ {
		l.Tick()
		fmt.Print("\x1bc", l)
		time.Sleep(time.Second / 30)
	}
}
