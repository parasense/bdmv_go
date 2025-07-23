package fontdir

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

// FontDirectory represents the root fontdirectory element
type FontDirectory struct {
	XMLName xml.Name `xml:"fontdirectory"`
	Fonts   []Font   `xml:"font"`
}

// Font represents a single font element
type Font struct {
	Name       string    `xml:"name"`
	FontFormat string    `xml:"fontformat"`
	Filename   string    `xml:"filename"`
	Styles     []string  `xml:"style"`
	Size       *FontSize `xml:"size"`
}

// FontSize represents the optional size element with min and max attributes
type FontSize struct {
	Min string `xml:"min,attr"`
	Max string `xml:"max,attr"`
}

// ParseFontDirectory parses the XML file at the given path into a FontDirectory struct
func ParseFontDirectory(filePath string) (*FontDirectory, error) {
	// Read the XML file
	xmlData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read XML file %s: %v", filePath, err)
	}

	var fontDir FontDirectory
	decoder := xml.NewDecoder(strings.NewReader(string(xmlData)))
	decoder.Strict = true // Enforce strict XML parsing to align with DTD

	err = decoder.Decode(&fontDir)
	if err != nil {
		return nil, fmt.Errorf("failed to parse XML: %v", err)
	}

	// Validate required fields per DTD
	for i, font := range fontDir.Fonts {
		if font.Name == "" {
			return nil, fmt.Errorf("font %d: missing required name element", i)
		}
		if font.FontFormat == "" {
			return nil, fmt.Errorf("font %d: missing required fontformat element", i)
		}
		if font.Filename == "" {
			return nil, fmt.Errorf("font %d: missing required filename element", i)
		}
		// Style is optional (0 or more), so no validation needed
		// Size is optional (0 or 1), so no validation needed
	}

	return &fontDir, nil
}
