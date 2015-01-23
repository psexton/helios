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
	"fmt"
	"encoding/json"
	log "github.com/cihub/seelog"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

func restorePackage(filePath string, conf Config) (err error) {
	// Read in the JSON file
	content, err := ioutil.ReadFile(filePath)
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
	log.Info("Restoring ", packageName)

	// Create an empty document so we can add attachments to it
	docRevision, err := putEmptyDocument(packageName, conf)
	if err != nil {
		return
	}

	log.Debug("Doc revision: ", docRevision)

	// Parse attachments array from JSON
	dirPath, _ := path.Split(filePath)
	attachments := packageData["_attachments"].(map[string]interface{})
	for fileName, _ := range attachments {
		tgzFilePath := path.Join(dirPath, fileName)
		docRevision, err = addAttachment(packageName, tgzFilePath, docRevision, conf)
		if err != nil {
			return
		}
	}

	// Overwrite document with full json
	err = overwriteDocument(packageData, docRevision, conf)

	return
}

func overwriteDocument(packageData map[string]interface{}, revision string, conf Config) (err error) {
	packageName := packageData["name"].(string)
	docURL := conf.Couch.URL + "registry/" + packageName
	log.Debug("Overwriting document at: ", docURL)	

	// We need to retrieve JSON for current document state so we don't lose _attachments when we restore
	// the version from s3.
	response1, err := http.Get(docURL)
	if err != nil {
		return
	}
	defer response1.Body.Close()
	body, err := ioutil.ReadAll(response1.Body)
	if err != nil {
		return
	}
	var serverDoc map[string]interface{} // holder for root JSON object
	err = json.Unmarshal(body, &serverDoc)
	if err != nil {
		return
	}

	// Update our JSON and Marshal it
	delete(packageData, "_id") // Couch sets _id & will balk if we try to set them ourselves
	packageData["_attachments"] = serverDoc["_attachments"] // Copy in _attachments from server
	packageData["_rev"] = revision // need to pass up the current doc revision to overwrite
	content, err := json.Marshal(packageData)
	if err != nil {
		return
	}

	// Make the PUT request.
	client := &http.Client{}
	request, err := http.NewRequest("PUT", docURL, bytes.NewReader(content))
	if err != nil {
		return
	}
	request.SetBasicAuth(conf.Couch.Username, conf.Couch.Password)
	response2, err := client.Do(request)
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

// Adds an attachment to an existing document in the registry database
// Returns the new revision ID for the document
func addAttachment(packageName string, tgzFilePath string, oldRevision string, conf Config) (revision string, err error) {
	// Figure out the URL to upload to
	_, fileName := path.Split(tgzFilePath)	
	tgzURL := conf.Couch.URL + "registry/" + packageName + "/" + fileName + "?rev=" + oldRevision
	log.Debug("Uploading attachment from ", tgzFilePath, " to ", tgzURL)	
	
	// Make an io.reader for the file
	file, err := os.Open(tgzFilePath)
	if err != nil {
		return
	}

	// Make the request
	client := &http.Client{}
	request, err := http.NewRequest("PUT", tgzURL, file)
	if err != nil {
		return
	}
	// NOTE: Apparently we don't need to set Content-Type header to "application/octet-stream"?
	request.SetBasicAuth(conf.Couch.Username, conf.Couch.Password)
	response, err := client.Do(request)
	if err != nil {
		return
	}
	log.Debug("File PUT ", tgzURL, " returned ", response.Status)
	if response.StatusCode != 201 {
		err = fmt.Errorf("PUT request to %s returned %d", tgzURL, response.StatusCode)
		return
	}

	// Get the new doc rev from the response
	defer response.Body.Close()	
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	var serverData map[string]interface{} // holder for arbitrary JSON
	err = json.Unmarshal(responseBody, &serverData)
	if err != nil {
		return
	}
	revision = serverData["rev"].(string) // again, we're pretty sure this is a string

	return
}

// Create an empty couch document in the registry database
// Returns the new revision ID for the document
func putEmptyDocument(packageName string, conf Config) (revision string, err error) {
	docURL := conf.Couch.URL + "registry/" + packageName
	log.Debug("Creating empty document at: ", docURL)	
	
	content := "{\"placeholder_for\":\"" + packageName + "\"}"
	
	// Make the PUT request. Should return 201 Created.
	client := &http.Client{}	
	request, err := http.NewRequest("PUT", docURL, strings.NewReader(content))
	if err != nil {
		return
	}
	request.SetBasicAuth(conf.Couch.Username, conf.Couch.Password)
	response, err := client.Do(request)
	if err != nil {
		return
	}
	log.Debug("PUT ", docURL, " ", content, " returned ", response.Status)
	if response.StatusCode != 201 {
		err = fmt.Errorf("PUT request to %s returned %d", docURL, response.StatusCode)
		return
	}

	// Get the new doc rev from the response
	defer response.Body.Close()	
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	var serverData map[string]interface{} // holder for arbitrary JSON
	err = json.Unmarshal(responseBody, &serverData)
	if err != nil {
		return
	}
	revision = serverData["rev"].(string) // again, we're pretty sure this is a string

	return
}

