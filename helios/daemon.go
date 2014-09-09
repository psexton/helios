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
		log.Debug("Daemon iterating")
		time.Sleep(30 * time.Second)
	}
}

