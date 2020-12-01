/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package cmds

import (
	"fmt"
	"net"
	"strings"
	"wgman/db"
)

type addInner struct{}

func (a *addInner) Help() (ret string) {
	return `
Usage: add-inner parent inner_name inner_ip

Adds a "inner" peer.

An inner peer is a host sits in internal network. The parent peer has to relay the
traffics between you and inner peer.
`
}

func (a *addInner) Exec(args []string) {
	if len(args) < 3 {
		fmt.Println(a.Help())
		return
	}

	parent := strings.TrimSpace(args[0])
	name := strings.TrimSpace(args[1])
	ip := strings.TrimSpace(args[2])
	if parent == "" || name == "" || ip == "" {
		fmt.Println(a.Help())
		return
	}

	if x := net.ParseIP(ip).To4(); x.String() != ip {
		fmt.Println(ip + " is not recongnized.")
		return
	}

	cfg, err := db.Load()
	if err != nil {
		fmt.Println("cannot load data: ", err)
		return
	}

	if _, ok := cfg.Peers[parent]; !ok {
		fmt.Println("peer " + parent + " not found")
		return
	}

	cfg.AddInner(parent, name, ip)

	if err = db.Save(cfg); err != nil {
		fmt.Println("Cannot save config: ", err)
		return
	}

	fmt.Printf("Inner peer %s is added after %s\n", name, parent)
}

func init() { register("add-inner", "adds an inner peer", &addInner{}) }
