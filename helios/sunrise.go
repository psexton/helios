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
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/brettweavnet/gosync/gosync"
	log "github.com/cihub/seelog"
	"github.com/mitchellh/goamz/aws"
	"io/ioutil"
	"net/http"
	"os/exec"
	"path"
	"strings"
)

// Sunrise imports from s3 to npm-registry
// 1) download bucket listing from s3
// 2) download files from s3
// 3) use npm to publish all tgz files
// 4) talk to couchdb directly to overwrite the json files
func Sunrise(conf Config) (err error) {
	const concurrent = 20 // @MAGIC

	// 1 & 2: use s3sync to download the bucket

	auth, err := aws.GetAuth(conf.AWS.AccessKeyID, conf.AWS.SecretAccessKey)
	if err != nil {
		return
	}
	source := "s3://" + conf.AWS.S3BucketName
	dest, err := ioutil.TempDir("", "helios")
	defer removeTempDir(dest) // delete our temp dir on exit
	if err != nil {
		return
	}
	log.Debug("auth:", auth)
	log.Debug("source:", source)
	log.Debug("dest:", dest)

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
		err = sunriseStep3(tgzFile)
		if err != nil {
			return
		}
	}

	// 4: talk to couchdb directly to overwrite the json files
	for _, jsonFile := range jsonFiles {
		err = sunriseStep4(jsonFile, conf)
		if err != nil {
			return
		}
	}

	return
}

func sunriseStep3(filepath string) (err error) {
	log.Info("Publishing ", filepath)
	err = exec.Command("npm", "publish", filepath).Run()
	return
}

func sunriseStep4(filepath string, conf Config) (err error) {
	log.Info("Restoring ", filepath)

	// Read in the JSON file
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return
	}
	var packageData map[string]interface{} // holder for arbitrary JSON
	err = json.Unmarshal(content, &packageData)
	if err != nil {
		return
	}

	// Extract the name element
	packageName := packageData["name"].(string) // we're pretty sure this is a string
	log.Debug("Package name: ", packageName)
	log.Debug("Old _rev: ", packageData["_rev"])

	// GET that document from Couch
	docURL := conf.Couch.URL + "/registry/" + packageName
	log.Debug("Document URL: ", docURL)
	client := &http.Client{}
	request1, err := http.NewRequest("GET", docURL, nil)
	if err != nil {
		return
	}
	response1, err := client.Do(request1)
	if err != nil {
		return
	}
	log.Debug("GET ", docURL, " returned ", response1.Status)
	if response1.StatusCode != 200 {
		err = fmt.Errorf("GET request to %s returned %d", docURL, response1.Status)
		return
	}
	defer response1.Body.Close()
	responseBody1, err := ioutil.ReadAll(response1.Body)
	if err != nil {
		return
	}
	var serverData map[string]interface{} // holder for arbitrary JSON
	err = json.Unmarshal(responseBody1, &serverData)
	if err != nil {
		return
	}
	log.Debug("New _rev: ", serverData["_rev"])

	// Replace the revision field in our json
	packageData["_rev"] = serverData["_rev"]

	// PUT the json into Couch
	content2, err := json.Marshal(packageData)
	if err != nil {
		return
	}
	request2, err := http.NewRequest("PUT", docURL, bytes.NewReader(content2))
	if err != nil {
		return
	}
	request2.SetBasicAuth(conf.Couch.Username, conf.Couch.Password)
	response2, err := client.Do(request2)
	if err != nil {
		return
	}
	log.Debug("PUT ", docURL, " returned ", response2.Status)
	if response2.StatusCode != 201 {
		err = fmt.Errorf("PUT request to %s returned %d", docURL, response2.Status)
		return
	}
	defer response2.Body.Close()

	return
}

