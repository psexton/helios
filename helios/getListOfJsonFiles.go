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
	"fmt"
	log "github.com/cihub/seelog"
	"io/ioutil"
	"net/http"
	"strings"
)

func getListOfJsonFiles(conf Config) (packages []string, err error) {
	log.Info("Downloading package list")
	
	docListURL := conf.Couch.URL + "registry/_all_docs"
	log.Debug("Doc List URL: ", docListURL)

	client := &http.Client{}
	request, err := http.NewRequest("GET", docListURL, nil)
	if err != nil {
		return
	}
	response, err := client.Do(request)
	if err != nil {
		return
	}
	log.Debug("GET ", docListURL, " returned ", response.Status)
	if response.StatusCode != 200 {
		err = fmt.Errorf("GET request to %s returned %d", docListURL, response.Status)
		return
	}
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	var serverData map[string]interface{} // holder for root JSON object
	err = json.Unmarshal(responseBody, &serverData)
	if err != nil {
		return
	}
	
	rows := serverData["rows"].([]interface{}) // JSON array

	packages = []string{}
	for i := range rows {
		entry := rows[i].(map[string]interface{}) // JSON object
		packageName := entry["id"].(string)
		if !strings.HasPrefix(packageName, "_") { // ignore design docs
			packages = append(packages, packageName)
		}
	}

	return
}

