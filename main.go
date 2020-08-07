package main

import (
	"fmt"
	"log"
	"os"
)

func printUsage() {
	fmt.Printf("\nUsage: go mod graph | modv <golang.org/x/text@v0.3.2> <results-dir>\n\n")
}

func main() {
	info, err := os.Stdin.Stat()

	if err != nil {
		fmt.Println("os.Stdin.Stat:", err)
		printUsage()
		os.Exit(1)
	}

	if info.Mode()&os.ModeNamedPipe == 0 {
		fmt.Println("modv not used in pipe, summarizing results")
		results := make(map[Module][]Module, 0)
		err := Load("all.gob", &results)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(results)
		os.Exit(1)
	}

	args := os.Args[1:]
	if len(args) > 2 {
		fmt.Printf("Expected 2 args got %d", len(args)-1)
		printUsage()
		os.Exit(1)
	}

	mg := NewModuleGraph(os.Stdin)
	if err := mg.Parse(); err != nil {
		fmt.Println("mg.Parse: ", err)
		printUsage()
		os.Exit(1)
	}

	if len(args) > 1 {
		target := args[0]
		if err := mg.FilterShort(newModule(target)); err != nil {
			fmt.Println("mg.Filter: ", err)
			printUsage()
			os.Exit(1)
		}
	}

	var targetFile string
	if len(args) == 2 {
		isDir, err := IsDir(args[1])
		if err != nil {
			log.Println(err)
		}
		if isDir {
			targetFile = fmt.Sprintf("%v/%v.gob", args[1], mg.name)
		} else {
			fmt.Printf("%v is not a directory", args[1])
		}
	} else {
		targetFile = fmt.Sprintf("%v.gob", mg.name)
	}
	if err := Save(mg.branches, targetFile); err != nil {
		fmt.Println("Save: ", err)
		printUsage()
		os.Exit(1)
	}
}
