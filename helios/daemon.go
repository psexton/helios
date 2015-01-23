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
	log "github.com/cihub/seelog"
	"time"
)

// Daemon loops infinitely
// * If getListOfJsonFiles returns empty, run Helios.Sunrise
// * Otherwise run Helios.Sunset
// * Then wait a bit
func Daemon(conf Config) (err error) {
	// First parse the duration from the config json
	pauseDuration, err := time.ParseDuration(conf.DaemonPause);
	if err != nil {
		return
	}

	for {
		log.Info("Daemon running...")

		// First, determine if registry is empty
		// If it is, we'll run sunrise.
		// If it isn't, we'll run sunset.
		jsonDocs, subErr := getListOfJsonFiles(conf)
		if subErr != nil {
			err = subErr
			return
		}
		log.Debug("jsonDocs: ", jsonDocs)

		if len(jsonDocs) == 0 {
			log.Info("Sunrise time!")
			err = Sunrise(conf)
			if err != nil {
				return
			}
		} else {
			log.Info("Sunset time!")
			err = Sunset(conf)
			if err != nil {
				return
			}
		}

		// Go have some tea
		log.Info("Nipping out for a bit of tea...")
		time.Sleep(pauseDuration)
	}
}
