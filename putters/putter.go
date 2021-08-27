package putters

import (
	"os"
	"path"

	"github.com/pilillo/igovium/utils"
)

func Put(tmpPath string, partName string, tmpFile string, config *utils.RemoteVolumeConfig) {
	if config.S3Config != nil {
		s3put(tmpPath, partName, tmpFile, config.S3Config)
	}

	// delete local file once put is done - iff is necessary
	if config.DeleteLocal {
		// what to do with the partitions? delete everything?
		// there should not be any other temp files around at the time of deletion
		// please check https://github.com/juju/fslock in case of concurrent access
		tmpFilePath := path.Join(tmpPath, partName, tmpFile)
		// remove the path and all its children
		err := os.RemoveAll(tmpFilePath)
		if err != nil {
			panic(err)
		}
	}
}
