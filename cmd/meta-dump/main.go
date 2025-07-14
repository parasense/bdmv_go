package main

import (
	"fmt"
	"os"

	meta "github.com/parasense/bdmv_go/pkg/meta"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: mobj-dump <mobj-file>")
		os.Exit(1)
	}

	metaPath := os.Args[1]

	discLib, err := meta.ParseMETA(metaPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing %s: %v\n", metaPath, err)

	}

	// Print parsed data
	fmt.Printf("File: \"%s\"\n", metaPath)
	fmt.Printf("Title: \"%s\"\n", discLib.DiscInfo.Title.Name)

	if discLib.DiscInfo.Title.NumSets != nil {
		fmt.Printf("Number of Sets: %d\n", *discLib.DiscInfo.Title.NumSets)
	}

	if discLib.DiscInfo.Title.SetNumber != nil {
		fmt.Printf("Set Number: %d\n", *discLib.DiscInfo.Title.SetNumber)
	}

	if discLib.DiscInfo.Language != nil {
		fmt.Printf("Language: %s\n", *discLib.DiscInfo.Language)
	}

	fmt.Println("Thumbnails:")
	for _, thumb := range discLib.DiscInfo.Description.Thumbnails {
		fmt.Printf("  Href: %s", thumb.Href)
		if thumb.Size != nil {
			fmt.Printf(", Size: %s", *thumb.Size)
		}
		fmt.Println()
	}
	if discLib.DiscInfo.Description.TableOfContents != nil {
		fmt.Println("Table of Contents:")
		for _, title := range discLib.DiscInfo.Description.TableOfContents.TitleNames {
			fmt.Printf("  Title %s: %s\n", title.TitleNumber, title.Name)
		}
	}
	fmt.Println()
}
