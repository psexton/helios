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
)

// Sunset exports from npm-registry to s3
// 1) talk to couchdb directly to get list of json files
// 2) download all json files
// 3) parse them and download all tgz attachments
// 4) upload everything to s3
func Sunset(conf Config) (err error) {
	///const concurrent = 20 // @MAGIC

	// create temp dir and do debug output
	tempDir, err := ioutil.TempDir("", "helios")
	defer removeTempDir(tempDir) // delete our temp dir on exit
	if err != nil {
		return
	}
	dest := "s3://" + conf.AWS.S3BucketName
	log.Debug("source: ", conf.Couch.URL)
	log.Debug("temp: ", tempDir)
	log.Debug("dest: ", dest)

	// step 1: talk to couchdb directly to get list of json files
	jsonDocs, err := getListOfJsonFiles(conf)
	if err != nil {
		return
	}
	log.Debug("jsonDocs: ", jsonDocs)

	// steps 2-3: download ALL THE THINGS
	for _, jsonDoc := range jsonDocs {
		log.Info("Backing up ", jsonDoc)
		err = downloadPackage(jsonDoc, tempDir, conf)
		if err != nil {
			return
		}
	}

	// step 4: sync the dir to s3
	err = syncToS3(tempDir, dest, conf)

	return
}

func syncToS3(sourceDir string, dest string, conf Config) (err error) {
	const concurrent = 20 // @MAGIC
	auth, err := aws.GetAuth(conf.AWS.AccessKeyID, conf.AWS.SecretAccessKey)
	if err != nil {
		return
	}
	

	syncPair := gosync.NewSyncPair(auth, sourceDir, dest)
	syncPair.Concurrent = concurrent
	err = syncPair.Sync()

	return
}
