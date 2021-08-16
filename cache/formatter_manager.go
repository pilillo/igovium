package cache

import (
	"fmt"
)

type FormatManager interface {
	Save(entries *[]DBCacheEntry, path string) error
}

// available formatters - lazy loaded formatters
var availableFormatters = map[string]func() FormatManager{
	"csv":     GetSingletonCSVFormatter,
	"parquet": GetSingletonParquetFormatter,
}

func GetFormatter(formatterType string) (FormatManager, error) {
	// return an instance of FormatManager if it exists
	if singleton, ok := availableFormatters[formatterType]; ok {
		return singleton(), nil
	}
	return nil, fmt.Errorf("Impossible to find specified formatter %s", formatterType)
}
