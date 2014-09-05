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
	"encoding/json"
	"io/ioutil"
	log "github.com/cihub/seelog"
	"net/http"
	"path"
)

// downloadPackage downloads a JSON doc, saves it, then downloads any attachments mentioned in it
func downloadPackage(packageName string, destDir string, conf Config) (err error) {
	log.Info("Downloading json doc for " + packageName)

	docURL := conf.Couch.URL + "registry/" + packageName
	log.Debug("Doc URL: ", docURL)
	
	// Download the JSON into memory
	resp, err := http.Get(docURL)
  	if err != nil {
		return
	}
  	defer resp.Body.Close()
  	body, err := ioutil.ReadAll(resp.Body)
  	if err != nil {
		return
	}

	// Write it to disk
	filePath := path.Join(destDir, packageName + ".json")
	log.Debug("Saving to: ", filePath)
  	err = ioutil.WriteFile(filePath, body, 0666)
  	if err != nil {
		return
	}

	// Parse out the attachments from the json
	var docRoot map[string]interface{} // holder for root JSON object
	err = json.Unmarshal(body, &docRoot)
	if err != nil {
		return
	}
	attachments := docRoot["_attachments"].(map[string]interface{})
	for attachment := range attachments {
		attachmentName := string(attachment)
		log.Debug("Found attachment ", attachmentName)
		attachmentURL := docURL + "/" + attachmentName
		attachmentFilePath := path.Join(destDir, attachmentName)
		err = downloadBinary(attachmentURL, attachmentFilePath, conf)
		if err != nil {
			return
		}
	}

	return
}

func downloadBinary(url string, destPath string, conf Config) (err error) {
	log.Debugf("Downloading %s to %s\n", url, destPath) // @TODO STUB
	return
}

