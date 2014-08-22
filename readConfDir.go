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
	"io/ioutil"
	"path"
)

// readConfDir takes a directory full of text files, and puts it into a map
// Filenames become keys, and file contents become values
// (Filenames are assumed to have no file extensions: e.g. ".txt")
func readConfDir(confDirPath string) (data map[string]string, err error) {
	// Get directory listing
	fileInfos, readDirErr := ioutil.ReadDir(confDirPath)
	if readDirErr != nil { // bail out
		err = readDirErr
		return
	}

	// Read the files
	data = make(map[string]string)
	for _, fileInfo := range fileInfos {
		filePath := path.Join(confDirPath, fileInfo.Name())
		contents, fileErr := ioutil.ReadFile(filePath)
		if fileErr != nil { // bail out
			err = fileErr
			return
		}
		// Assign into map
		key := fileInfo.Name()
		value := string(contents)
		data[key] = value
	}

	return
}
