/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"os"
	"wgman/cmds"
	"wgman/db"
)

func main() {
	cfg := os.Getenv("WGMAN_CONFIG")
	if cfg != "" {
		db.ConfigFilename = cfg
	}

	cmds.Run(os.Args[1:])
}
