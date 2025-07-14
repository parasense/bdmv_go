package meta

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

type LanguageCode [3]byte

type MetaData struct {
	Language LanguageCode
	//DiscLib
}

func ParseMETA(filePath string) (discLib *DiscLib, err error) {
	discLib = &DiscLib{}

	// Open the XML file
	file, err := os.Open(filePath)
	if err != nil {
		return discLib, fmt.Errorf("error opening file %s: %w", filePath, err)
	}
	defer file.Close()

	// Read the entire file
	data, err := io.ReadAll(file)
	if err != nil {
		return discLib, fmt.Errorf("error reading file %s: %w", filePath, err)
	}

	// Unmarshal the XML data into the DiscLib struct
	err = xml.Unmarshal(data, &discLib)
	if err != nil {
		return discLib, fmt.Errorf("error unmarshaling XML from %s: %w", filePath, err)
	}

	return discLib, nil
}
