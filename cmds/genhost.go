/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package cmds

import (
	"fmt"
	"strings"
	"wgman/db"
)

type genhost struct{}

func (g *genhost) Help() (ret string) {
	return `
Usage: genhost [--help|-h] [domain]

Generate relevent /etc/hosts entries.

The optional domain puts peers in a domain. For example, "wgman genhost mynet"
generates "host.mynet".
`
}

func (g *genhost) Exec(args []string) {
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

	for name, p := range cfg.Peers {
		myip := cfg.IP(p.IP)
		fmt.Printf("%-16s %s\n", myip, f(name))

		for n, i := range p.Inners {
			myip := cfg.IP(i.IP)
			fmt.Printf("%-16s %s\n", myip, f(n))
		}
	}
}

func init() { register("genhost", "generate relevent /etc/host entries", &genhost{}) }
