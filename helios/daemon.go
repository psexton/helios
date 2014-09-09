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
	"time"
	log "github.com/cihub/seelog"
)

// Daemon loops infinitely
// * If getListOfJsonFiles returns empty, run Helios.Sunrise
// * Otherwise run Helios.Sunset
// * Wait 30 seconds
func Daemon(conf Config) (err error) {
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
			log.Debug("Sunrise time!")
			err = Sunrise(conf)
			if err != nil {
				return
			}
		} else {
			log.Debug("Sunset time!")
			err = Sunset(conf)
			if err != nil {
				return
			}
		}

		// Go have some tea
		log.Debug("Nipping out for a bit of tea...")
		time.Sleep(30 * time.Second)
	}
}

