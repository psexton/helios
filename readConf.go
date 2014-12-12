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
	"encoding/json"
	"github.com/psexton/helios/helios"
	"io/ioutil"
)

// readConf takes a json file, and decodes it into a helios/Config object
func readConf(confPath string) (conf helios.Config, err error) {
	// Read the file into a byte slice
	fileContents, err := ioutil.ReadFile(confPath)
	if err != nil {
		return
	}

	// Unmarshal
	err = json.Unmarshal(fileContents, &conf)

	// There's no step 3! Also chaos theory.

	return
}
