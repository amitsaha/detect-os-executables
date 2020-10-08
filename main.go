package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gabriel-vasile/mimetype"
)

var executableTypeString = map[string]bool{
	"application/x-mach-binary":                     true,
	"application/vnd.microsoft.portable-executable": true,
	"application/x-elf":                             true,
	"application/x-object":                          true,
	"application/x-executable":                      true,
}

func isBinary(path string) (bool, error) {
	mime, err := mimetype.DetectFile(path)
	if err != nil {
		return false, err
	}
	return executableTypeString[mime.String()], nil
}

type ActionConfig struct {
	gitIgnoreGen bool
	cleanup      bool
}

func main() {
	config := ActionConfig{}
	fs := flag.NewFlagSet("detect-executables", flag.ContinueOnError)
	fs.BoolVar(&config.gitIgnoreGen, "g", false, "Generate a gitignore")
	fs.BoolVar(&config.cleanup, "c", false, "Cleanup")
	fs.Parse(os.Args[1:])
	if fs.NArg() != 1 {
		fs.PrintDefaults()
		os.Exit(1)
	}
	rootDir := fs.Args()[0]

	var executables []string
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if !info.IsDir() {
			b, err1 := isBinary(path)
			if err1 != nil {
				fmt.Printf("Couldn't check file type for: %s, error: %v\n", path, err1)
			}
			if b {
				executables = append(executables, path)
			}
			return nil
		}
		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", os.Args[1], err)
	}

	// the output is also a gitignore format
	for _, e := range executables {
		fmt.Printf("%s\n", e)
	}

	if config.cleanup {
		for _, e := range executables {
			err = os.Remove(e)
			if err != nil {
				fmt.Println(e)
				os.Exit(1)
			}
		}
	}
}
