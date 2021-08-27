package cache

import (
	"fmt"
	"os"

	"sync"

	"github.com/pilillo/igovium/utils"
	"github.com/xitongsys/parquet-go/parquet"
	"github.com/xitongsys/parquet-go/writer"
)

var parquetFormatterOnce sync.Once
var parquetFormatterInstance *parquetFormatter

type parquetFormatter struct{}

// NewParquetFormatter ... constructor of FormatManager of parquet type
func NewParquetFormatter() FormatManager {
	return &parquetFormatter{}
}

// GetSingletonParquetFormatter ... lazy singleton on DAO
func GetSingletonParquetFormatter() FormatManager {
	// once.do is lazy, we use it to return an instance of the DAO
	parquetFormatterOnce.Do(func() {
		parquetFormatterInstance = &parquetFormatter{}
	})
	return parquetFormatterInstance
}

type parquetEntry struct {
	Key       string `parquet:"name=key, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	Value     string `parquet:"name=value, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	CreatedAt int64  `parquet:"name=created, type=INT64"`
	UpdatedAt int64  `parquet:"name=updated, type=INT64"`
}

func (f *parquetFormatter) Save(entries *[]DBCacheEntry, path string) error {
	w, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("Can't create local file: %v", err)
	}
	defer w.Close()

	pw, err := writer.NewParquetWriterFromWriter(w, new(parquetEntry), 4)
	if err != nil {
		return fmt.Errorf("Can't create parquet writer: %v", err)
	}

	pw.RowGroupSize = 128 * 1024 * 1024 //128M
	pw.PageSize = 8 * 1024              //8K
	pw.CompressionType = parquet.CompressionCodec_SNAPPY

	for _, cacheEntry := range *entries {
		if err = pw.Write(&parquetEntry{
			Key:       cacheEntry.Key,
			Value:     utils.ToBase64String(cacheEntry.Value),
			CreatedAt: cacheEntry.CreatedAt,
			UpdatedAt: cacheEntry.UpdatedAt,
		}); err != nil {
			return fmt.Errorf("Write error: %v", err)
		}
	}
	pw.Flush(true)
	if err = pw.WriteStop(); err != nil {
		return fmt.Errorf("WriteStop error: %v", err)
	}

	return nil
}
