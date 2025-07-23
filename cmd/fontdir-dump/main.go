package main

import (
	"fmt"
	"log"
	"os"

	fontdir "github.com/parasense/bdmv_go/pkg/fontdir"
)

// PrintFontDirectory prints the FontDirectory struct to the terminal
func PrintFontDirectory(fontDir *fontdir.FontDirectory) {
	for i, font := range fontDir.Fonts {
		fmt.Printf("Font %d:\n", i+1)
		fmt.Printf("  Name: %s\n", font.Name)
		fmt.Printf("  Font Format: %s\n", font.FontFormat)
		fmt.Printf("  Filename: %s\n", font.Filename)
		fmt.Printf("  Styles: %v\n", font.Styles)
		if font.Size != nil {
			fmt.Printf("  Size: min=%s, max=%s\n", font.Size.Min, font.Size.Max)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run fontdir-dump <path-to-xml-file>")
	}

	filePath := os.Args[1]
	fontDir, err := fontdir.ParseFontDirectory(filePath)
	if err != nil {
		log.Fatalf("Error parsing font directory: %v", err)
	}

	PrintFontDirectory(fontDir)
}
