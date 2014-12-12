/*
This file is part of Helios.

Helios is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Helios is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Helios.  If not, see <http://www.gnu.org/licenses/>.
*/

package helios

import (
	log "github.com/cihub/seelog"
	"github.com/mitchellh/goamz/aws"
	"github.com/psexton/gosync/gosync"
	"io/ioutil"
	"path"
	"strings"
)

// Sunrise imports from s3 to npm-registry
// 1) download bucket listing from s3
// 2) download files from s3
// 3) talk to couchdb directly to overwrite the registry documents
func Sunrise(conf Config) (err error) {
	const concurrent = 20 // @MAGIC

	// 1 & 2: use s3sync to download the bucket

	auth, err := aws.GetAuth(conf.AWS.AccessKeyID, conf.AWS.SecretAccessKey)
	if err != nil {
		return
	}
	// create temp dir and do debug output
	source := "s3://" + conf.AWS.S3BucketName
	dest, err := ioutil.TempDir("", "helios")
	defer removeTempDir(dest) // delete our temp dir on exit
	if err != nil {
		return
	}
	log.Debug("source: ", source)
	log.Debug("temp: ", dest)
	log.Debug("dest: ", conf.Couch.URL)

	syncPair := gosync.NewSyncPair(auth, source, dest)
	syncPair.Concurrent = concurrent
	err = syncPair.Sync()
	if err != nil {
		return
	}

	// 2.5: Get directory listing and separate tgz from json
	fileInfos, err := ioutil.ReadDir(dest)
	if err != nil {
		return
	}
	jsonFiles := []string{}
	for _, fileInfo := range fileInfos {
		filePath := path.Join(dest, fileInfo.Name())
		if strings.HasSuffix(filePath, ".json") {
			jsonFiles = append(jsonFiles, filePath)
		}
	}

	// 3: call restorePackage on each json file
	for _, jsonFile := range jsonFiles {
		err = restorePackage(jsonFile, conf)
		if err != nil {
			return
		}
	}

	return
}

