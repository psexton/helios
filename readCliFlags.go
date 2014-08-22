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
	"flag"
	"fmt"
	"os"
)

// readCliFlags reads in command line arguments and flags and validates them
func readCliFlags() (confDir string, isSunrise bool, err error) {
	sunrisePtr := flag.Bool("sunrise", false, "Import data from s3 to couchdb")
	sunsetPtr := flag.Bool("sunset", false, "Export data from couchdb to s3")
	confDirPtr := flag.String("confdir", "/path/to/conf/files", "Directory containing conf files")
	flag.Parse()

	// Err out if command isn't "sunrise" or "sunset"
	if !(*sunrisePtr || *sunsetPtr) {
		err = fmt.Errorf("either --sunrise or --sunset must be specified")
		return
	}
	if *sunrisePtr && *sunsetPtr {
		err = fmt.Errorf("can't do both sunrise and sunset")
		return
	}
	isSunrise = *sunrisePtr

	// Err out if confdir isn't a directory
	//	First check if the path exists by trying to open it as a file
	file, fileErr := os.Open(*confDirPtr)
	if fileErr != nil {
		err = fileErr
		return
	}
	//	Then Stat it to check it's a dir
	fileStat, statErr := file.Stat()
	if statErr != nil {
		err = statErr
		return
	}
	if !fileStat.IsDir() {
		err = fmt.Errorf("\"%s\" is not a directory", *confDirPtr)
		return
	}
	confDir = *confDirPtr
	return
}
