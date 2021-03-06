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
	confPath, command, err := readCliFlags()
	exitOnError(err)

	// read in data from json conf file
	// conf is a helios/Config struct
	conf, err := readConf(confPath)
	exitOnError(err)

	// If command is daemon, we want to reroute log output to a file
	if command == daemonCmd {
		sendLogsToFile(conf)
	}

	// call sunrise or sunset, passing it the conf data
	switch command {
	case sunriseCmd:
		log.Debug("command: sunrise")
		err := helios.Sunrise(conf)
		exitOnError(err)
	case sunsetCmd:
		log.Debug("command: sunset")
		err := helios.Sunset(conf)
		exitOnError(err)
	case daemonCmd:
		log.Debug("command: daemon")
		err := helios.Daemon(conf)
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

func sendLogsToFile(conf helios.Config) {
	// put logs in file specified by config
	// filter by level specified in config
	filePath := conf.Log.FilePath
	filterLevel := conf.Log.Level
	maxSize := "1000000"

	loggerConfig := "<seelog minlevel=\"" + filterLevel + "\"><outputs><rollingfile type=\"size\" filename=\"" + filePath + "\" maxsize=\"" + maxSize + "\" maxrolls=\"5\"/></outputs></seelog>"
	logger, _ := log.LoggerFromConfigAsString(loggerConfig)
	log.ReplaceLogger(logger)
}
