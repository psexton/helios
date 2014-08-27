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
	log "github.com/cihub/seelog"
	"github.com/psexton/helios/helios"
	"os"
)

func main() {
	defer log.Flush()

	// read in CLI flags and arguments
	confDir, isSunrise, err := readCliFlags()
	exitOnError(err)
	log.Info("confDir: ", confDir)
	log.Info("isSunrise: ", isSunrise)

	// read in data from conf dir
	// conf is a map[string]string
	conf, err := readConfDir(confDir)
	exitOnError(err)

	// call sunrise or sunset, passing it the conf data
	if isSunrise {
		err := helios.Sunrise(conf)
		exitOnError(err)
	} else {
		err := helios.Sunset(conf)
		exitOnError(err)
	}
}

// utility func taken from gosync's main.go
func exitOnError(e error) {
	if e != nil {
		log.Errorf("Received error '%s'", e.Error())
		log.Flush()
		os.Exit(1)
	}
}
