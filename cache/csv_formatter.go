package cache

import (
	"encoding/csv"
	"fmt"
	"os"
	"sync"

	"github.com/pilillo/igovium/utils"
)

var csvFormatterOnce sync.Once
var csvFormatterInstance *csvFormatter

type csvFormatter struct{}

// NewCSVFormatter ... constructor of FormatManager of csv type
func NewCSVFormatter() FormatManager {
	return &csvFormatter{}
}

// GetSingletonCSVFormatter ... lazy singleton on DAO
func GetSingletonCSVFormatter() FormatManager {
	// once.do is lazy, we use it to return an instance of the DAO
	csvFormatterOnce.Do(func() {
		csvFormatterInstance = &csvFormatter{}
	})
	return csvFormatterInstance
}

func (f *csvFormatter) Save(entries *[]DBCacheEntry, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// shall we write a header?
	header := []string{"key", "value", "created_at", "updated_at"}
	err = writer.Write(header)
	if err != nil {
		return err
	}

	for _, entry := range *entries {
		record := []string{
			entry.Key,
			utils.ToBase64String(entry.Value),
			//entry.Value,
			fmt.Sprint(entry.CreatedAt),
			fmt.Sprint(entry.UpdatedAt),
		}
		err := writer.Write(record)
		if err != nil {
			return err
		}
	}

	return nil
}
