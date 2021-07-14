package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/gosuri/uitable"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/flavio/rego-builtin-finder/rego"
)

func main() {
	debug := flag.Bool("debug", false, "sets log level to debug")
	flag.Parse()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if len(flag.Args()) != 1 {
		fmt.Println("Must provide the path to the directory/file to scan")
		os.Exit(1)
	}

	fi, err := os.Lstat(flag.Arg(0))
	if err != nil {
		log.Error().
			Err(err).
			Str("target", flag.Arg(0))
		os.Exit(1)
	}

	var regoFiles []string
	switch mode := fi.Mode(); {
	case mode.IsDir():
		regoFiles, err = findRegoFiles(flag.Arg(0))
		if err != nil {
			log.Error().
				Err(err).
				Str("path", flag.Arg(0)).
				Msg("error walking the path")
			os.Exit(1)
		}
	case mode.IsRegular():
		regoFiles = []string{flag.Arg(0)}
	default:
		log.Error().
			Msg("Provide either a dir or a regular file")
		os.Exit(1)
	}

	inspector := rego.NewInspector()

	// a map with builtin name as key, and the number of rego policies using
	// it as value
	builtinsUsage := make(map[string]int)

	for _, filename := range regoFiles {
		required, err := inspector.InspectPolicy(filename)
		if err != nil {
			log.Error().
				Err(err).
				Msg(fmt.Sprintf("Something went wrong while inspecting %s", filename))
		}
		for builtin := range required.Iterator().C {
			val, found := builtinsUsage[builtin.(string)]
			if !found {
				val = 1
			} else {
				val += 1
			}
			builtinsUsage[builtin.(string)] = val
		}
	}

	builtinsSortedByUsage := SortMapByValue(builtinsUsage)

	table := uitable.New()
	table.MaxColWidth = 80
	table.Wrap = true
	table.AddRow("NAME", "OCCURRENCES")

	for _, builtin := range builtinsSortedByUsage {
		table.AddRow(builtin.Key, builtin.Value)
	}

	fmt.Printf("Rego files analyzed: %d\n", len(regoFiles))
	fmt.Println("List of builtins that have to be provided by the SDK")
	fmt.Println()
	fmt.Println(table)
}

func findRegoFiles(dirToScan string) ([]string, error) {
	regoFiles := []string{}
	err := filepath.Walk(dirToScan, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(info.Name()) == ".rego" {
			filename := strings.TrimSuffix(filepath.Base(info.Name()), ".rego")
			if !strings.HasPrefix(filename, "test_") && !strings.HasSuffix(filename, "_test") {
				regoFiles = append(regoFiles, path)
			}
		}
		return nil
	})

	if err != nil {
		return []string{}, err
	}

	return regoFiles, nil
}
