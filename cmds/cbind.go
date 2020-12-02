/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package cmds

import (
	"fmt"
	"strings"
	"wgman/db"
)

type cbind struct{}

func (g *cbind) Help() (ret string) {
	return `
Usage: cbind [--help|-h] [domain]

Generate relevent config to use in https://hub.docker.com/r/cytopia/bind

The optional domain puts peers in a domain. For example, "wgman cbind mynet"
generates "host.mynet".
`
}

func (g *cbind) Exec(args []string) {
	var domain = ""
	f := func(n string) string {
		if domain == "" {
			return n
		}

		return n + "." + domain
	}

	if len(args) > 0 {
		if s := strings.TrimSpace(args[0]); s == "--help" || s == "-h" {
			fmt.Println(g.Help())
			return
		}
		domain = strings.TrimLeft(strings.TrimSpace(args[0]), ".")
	}

	cfg, err := db.Load()
	if err != nil {
		fmt.Println("cannot load data: ", err)
		return
	}

	arr := []string{}

	for name, p := range cfg.Peers {
		myip := cfg.IP(p.IP)
		arr = append(arr, fmt.Sprintf("%s=%s", f(name), myip))

		for n, i := range p.Inners {
			myip := cfg.IP(i.IP)
			arr = append(arr, fmt.Sprintf("%s=%s", f(n), myip))
		}
	}

	fmt.Println(strings.Join(arr, ","))
}

func init() { register("cbind", "generate config to use with cytopia/bind", &cbind{}) }
