package indx

import (
	"fmt"
	"os"
)

func ParseINDX(filePath string) (
	header *INDXHeader,
	appinfo *AppInfo,
	indexes *Indexes,
	extensiondata *Extensions,
	err error,
) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Header
	if header, err = ReadINDXHeader(file); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to read header: %w", err)
	}

	// AppInfo
	if appinfo, err = ReadAppInfo(file, header.AppInfo); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to read appinfo: %w", err)
	}

	// Indexes
	if indexes, err = ReadIndexes(file, header.Indexes); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to read indexes: %w", err)
	}

	// Extensions
	if header.Extensions.Start != 0 && header.Extensions.Stop != 0 {
		if extensiondata, err = ReadExtensions(file, header.Extensions); err != nil {
			return nil, nil, nil, nil, fmt.Errorf("failed to read Extension Data: %w", err)
		}
	}
	return header, appinfo, indexes, extensiondata, nil
}
