/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package cmds

import (
	"fmt"
	"strings"
	"wgman/db"
)

type list struct{}

func (l *list) Help() (ret string) {
	return `
Usage: list

List all peers.
`
}

func (l *list) Exec(args []string) {
	cfg, err := db.Load()
	if err != nil {
		fmt.Println("cannot load data: ", err)
		return
	}

	fmt.Printf("%-24s %-16s %s\n", "name", "address", "endpoint")
	fmt.Println(strings.Repeat("-", 80))
	for n, v := range cfg.Peers {
		fmt.Printf(
			"%-24s %-16s %s\n",
			n,
			cfg.IP(v.IP),
			v.Endpoint(),
		)

		for n, i := range v.Inners {
			fmt.Printf(
				"  - %-20s %-16s\n",
				n,
				cfg.IP(i.IP),
			)
		}
	}
}

func init() { register("list", "list all peers", &list{}) }
