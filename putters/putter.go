package putters

import (
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/pilillo/igovium/utils"
)

func getBasePath(path string) string {
	if filepath.Dir(path) == "." || filepath.Dir(path) == "/" {
		return path
	} else {
		// get parent dir
		return getBasePath(filepath.Dir(path))
	}
}

func Put(tmpPath string, partName string, tmpFile string, config *utils.RemoteVolumeConfig) {
	if config.S3Config != nil {
		s3put(tmpPath, partName, tmpFile, config.S3Config)
	}

	// delete local file once put is done - iff is necessary
	if config.DeleteLocal {
		// what to do with the partitions? delete everything?
		// there should not be any other temp files around at the time of deletion
		// please check https://github.com/juju/fslock in case of concurrent access
		tmpFilePath := path.Join(partName, tmpFile)
		// remove the part path and all its children
		dataRoot := getBasePath(tmpFilePath)
		log.Printf("Removing directory at %s", dataRoot)
		err := os.RemoveAll(dataRoot)
		if err != nil {
			panic(err)
		}
	}
}
