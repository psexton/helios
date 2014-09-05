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
	"github.com/psexton/helios/helios"
	"io/ioutil"
	"path"
	"strings"
)

// readConfDir takes a directory full of text files, and puts it into a map
// Filenames become keys, and file contents become values
// (Filenames are assumed to have no file extensions: e.g. ".txt")
func readConfDir(confDirPath string) (conf helios.Config, err error) {
	// Get directory listing
	fileInfos, readDirErr := ioutil.ReadDir(confDirPath)
	if readDirErr != nil { // bail out
		err = readDirErr
		return
	}

	// Read the files
	data := make(map[string]string)
	for _, fileInfo := range fileInfos {
		filePath := path.Join(confDirPath, fileInfo.Name())
		contents, fileErr := ioutil.ReadFile(filePath)
		if fileErr != nil { // bail out
			err = fileErr
			return
		}
		// Assign into map
		key := fileInfo.Name()
		value := strings.Trim(string(contents), "\n") // strip out trailing newline
		data[key] = value
	}

	// Populate the helios.Config struct from the map
	conf.AWS.AccessKeyID = data["aws_access_key_id"]
	conf.AWS.SecretAccessKey = data["aws_secret_access_key"]
	conf.AWS.S3BucketName = data["s3_bucket"]
	conf.Couch.Username = data["couch_username"]
	conf.Couch.Password = data["couch_password"]
	conf.Couch.URL = data["couch_url"]

	return
}
