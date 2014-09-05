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
	"io/ioutil"
	log "github.com/cihub/seelog"
)

// Sunset exports from npm-registry to s3
// 1) talk to couchdb directly to get list of json files
// 2) download all json files
// 3) parse them and download all tgz attachments
// 4) upload everything to s3
func Sunset(conf Config) (err error) {
	///const concurrent = 20 // @MAGIC

	// step 1: talk to couchdb directly to get list of json files
	jsonDocs, err := getListOfJsonFiles(conf)
	if err != nil {
		return
	}
	log.Debug("jsonDocs: ", jsonDocs)

	tempDir, err := ioutil.TempDir("", "helios")
	defer removeTempDir(tempDir) // delete our temp dir on exit
	if err != nil {
		return
	}
	log.Debug("tempDir: ", tempDir)
	
	// step 2: download ALL THE DOCUMENTS	

	return
}

