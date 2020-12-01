/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package cmds

import (
	"fmt"
	"net"
	"wgman/db"
)

type create struct{}

func (n *create) Help() (ret string) {
	return `
Usage: new network_addr

Creates new config data, using network_addr as network. (always with 24bit mask)

Example: new 10.1.1.1
This will create a new network from 10.1.1.1 to 10.1.1.254.
`
}

func (n *create) Exec(args []string) {
	if len(args) < 1 {
		fmt.Print(n.Help())
		return
	}

	addr := args[0]
	ipn := net.ParseIP(addr)
	cfg := &db.Network{
		Net: net.IPNet{
			IP:   ipn,
			Mask: net.IPv4Mask(255, 255, 255, 0),
		},
		Peers: map[string]*db.Peer{},
	}

	err := db.Save(cfg)
	if err != nil {
		fmt.Println("Failed to save config: ", err)
		return
	}

	fmt.Println("Config file created.")
}

func init() { register("new", "initialize config file", &create{}) }
