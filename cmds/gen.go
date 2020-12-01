/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package cmds

import (
	"fmt"
	"strings"
	"wgman/db"
)

type gen struct{}

func (g *gen) Help() (ret string) {
	return `
Usage: gen name

Generate relevent, wg-quick compitable config file.
`
}

func (g *gen) Exec(args []string) {
	if len(args) < 1 {
		fmt.Println(g.Help())
		return
	}

	name := strings.TrimSpace(args[0])
	if len(name) == 0 {
		fmt.Println(g.Help())
		return
	}

	cfg, err := db.Load()
	if err != nil {
		fmt.Println("cannot load data: ", err)
		return
	}

	if _, ok := cfg.Peers[name]; !ok {
		(&list{}).Exec(nil)
		return
	}

	txt, err := cfg.GenConf(name)
	if err != nil {
		fmt.Println("cannot generate config: ", err)
		return
	}

	fmt.Println(txt)
}

func init() { register("gen", "generate relevent config file", &gen{}) }
