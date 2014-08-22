/*
helios is a simple app for backing up and restoring a npm registry to s3.


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
	"fmt"
)

func main() {
	// read in CLI flags and arguments
	confDir, isSunrise, err := readCliFlags()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("confdir:", confDir)     // @DEBUG
	fmt.Println("isSunrise:", isSunrise) // @DEBUG

	// read in data from conf dir
	// conf is a map[string]string
	conf, err := readConfDir(confDir)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("conf:", conf) // @DEBUG

	// SUNRISE
	// download bucket listing from s3
	// download files from s3
	// use npm to publish all tgz files
	// talk to couchdb directly to overwrite the json files

	// SUNSET
	// talk to couchdb directly to get list of json files
	// download all json files
	// parse them and download all tgz attachments
	// upload everything to s3
}
