package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"janitor/internal/config"
	"janitor/internal/organizer"
)

// Constants for version and developer info
const (
	version       = "Janitor CLI - v0.0.1"
	developerInfo = "Janitor CLI - Created by Mohd Zaid"
)

func main() {
	// 1. Define our flags
	// The flag package automatically handles -h and -help
	versionFlag := flag.Bool("v", false, "Print Janitor's version")
	infoFlag := flag.Bool("i", false, "Print developer info")

	// 2. Parse the flags
	flag.Parse()

	// 3. Act on the flags if they were provided
	if *versionFlag {
		fmt.Println(version)
		os.Exit(0) // Exit cleanly
	}
	if *infoFlag {
		fmt.Println(developerInfo)
		os.Exit(0) // Exit cleanly
	}

	// 4. Load configuration (same as before)
	// This function will implement the logic of checking ./, ~/, and then defaults
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// 5. Determine the target directory
	var targetPath string

	// flag.Args() holds the remaining arguments after flags are parsed
	if flag.NArg() > 1 {
		log.Println("Error: Too many arguments. Please provide at most one directory path.")
		flag.Usage() // Print the help message
		os.Exit(1)   // Exit with an error
	}

	if flag.NArg() == 1 {
		// A path was provided
		targetPath = flag.Arg(0)
	} else {
		// No path provided, use the current working directory
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Error getting current directory: %v", err)
		}
		targetPath = cwd
	}

	// 6. Create and run the organizer
	log.Printf("Janitor starting... (Target: %s)", targetPath)
	org := organizer.New(cfg, targetPath)

	if err := org.Run(); err != nil {
		log.Fatalf("Error during organization: %v", err)
	}

	log.Println("Janitor run complete. Directory is clean!")
}
