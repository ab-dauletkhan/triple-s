package core

import (
	"errors"
	"flag"
	"fmt"
)

var (
	Port int
	Dir  string
	Help bool
)

var (
	ErrIncorrectPort = errors.New("incorrect port number, range must be between 1-65535")
	ErrEmptyDir      = errors.New("empty directory path")
)

// Parses the above three flags
// Port number should be in range 1 and 65535 inclusively.
func ParseFlags() error {
	flag.IntVar(&Port, "port", 8080, "server port to listen on")
	flag.StringVar(&Dir, "dir", "./data", "directory to store buckets")
	flag.BoolVar(&Help, "help", false, "print help message")

	flag.Usage = PrintUsage
	flag.Parse()

	if Port < 1 || Port > 65535 {
		return ErrIncorrectPort
	}

	if Dir == "" {
		return ErrEmptyDir
	}

	return nil
}

func PrintUsage() {
	fmt.Println(`Simple Storage Service.
Usage:
	triple-s [-port <N>] [-dir <S>]
	triple-s --help
Options:
	--help     Show this screen.
	--port N   Port number
	--dir S    Path to the directory`)
}
