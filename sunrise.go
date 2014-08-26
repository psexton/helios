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

package main

import (
	"github.com/brettweavnet/gosync/gosync"
	log "github.com/cihub/seelog"
	"github.com/mitchellh/goamz/aws"
	"io/ioutil"
	"os/exec"
	"path"
	"strings"
)

// sunrise imports from s3 to npm-registry
// 1) download bucket listing from s3
// 2) download files from s3
// 3) use npm to publish all tgz files
// 4) talk to couchdb directly to overwrite the json files
func sunrise(conf map[string]string) (err error) {
	const concurrent = 20 // @MAGIC

	// 1 & 2: use s3sync to download the bucket

	auth, authErr := aws.GetAuth(conf["aws_access_key_id"], conf["aws_secret_access_key"])
	if authErr != nil {
		err = authErr
		return
	}
	source := "s3://" + conf["s3_bucket"]
	dest, tempDirErr := ioutil.TempDir("", "helios")
	// @TODO Add deferred call to delete dest
	if tempDirErr != nil {
		err = tempDirErr
		return
	}
	log.Debug("auth:", auth)
	log.Debug("source:", source)
	log.Debug("dest:", dest)

	syncPair := gosync.NewSyncPair(auth, source, dest)
	syncPair.Concurrent = concurrent
	syncErr := syncPair.Sync()
	if syncErr != nil {
		err = syncErr
		return
	}

	// 2.5: Get directory listing and separate tgz from json
	fileInfos, readDirErr := ioutil.ReadDir(dest)
	if readDirErr != nil {
		err = readDirErr
		return
	}
	tgzFiles := []string{}
	jsonFiles := []string{}
	for _, fileInfo := range fileInfos {
		filePath := path.Join(dest, fileInfo.Name())
		if strings.HasSuffix(filePath, ".tgz") {
			tgzFiles = append(tgzFiles, filePath)
		}
		if strings.HasSuffix(filePath, ".json") {
			jsonFiles = append(jsonFiles, filePath)
		}
	}

	// 3: use npm to publish all tgz files
	for _, tgzFile := range tgzFiles {
		sunriseStep3(tgzFile, conf)
	}

	// 4: talk to couchdb directly to overwrite the json files
	for _, jsonFile := range jsonFiles {
		sunriseStep4(jsonFile, conf)
	}

	return
}

func sunriseStep3(filepath string, conf map[string]string) (err error) {
	log.Info("Publishing ", filepath)
	err = exec.Command("npm", "publish", filepath).Run()
	return
}

func sunriseStep4(filepath string, conf map[string]string) (err error) {
	log.Info("[PLACEHOLDER] Overwriting ", filepath)
	return
}
