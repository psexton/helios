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
	"errors"
	"flag"
	"fmt"
	///"launchpad.net/goamz/aws"
	///"launchpad.net/goamz/s3"
	///"log"
)

func main() {
	// read in CLI flags and arguments
	confDir, isSunrise, err := readCliFlags()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("confdir:", confDir) // @DEBUG
	fmt.Println("isSunrise:", isSunrise) // @DEBUG

	// read in data from conf dir

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

func readCliFlags() (confDir string, isSunrise bool, err error) {
	sunrisePtr := flag.Bool("sunrise", false, "Import data from s3 to couchdb")
	sunsetPtr := flag.Bool("sunset", false, "Export data from couchdb to s3")
	confDirPtr := flag.String("confdir", "/path/to/conf/files", "Directory containing conf files")
	flag.Parse()

	// Err out if command isn't "sunrise" or "sunset"
	if !(*sunrisePtr || *sunsetPtr) {
		err = errors.New("Either --sunrise or --sunset must be specified")
		return
	}
	if (*sunrisePtr && *sunsetPtr) {
		err = errors.New("Can't do both sunrise and sunset")
		return
	}
	isSunrise = *sunrisePtr

	// Err out if confdir isn't a directory
	// @TODO
	confDir = *confDirPtr	

	return
}
