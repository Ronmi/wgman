/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package cmds

import (
	"fmt"
	"strconv"
	"strings"
	"wgman/db"
)

type add struct{}

func (a *add) Help() (ret string) {
	return `
Usage: add name [ip:port]

Adds a peer in database. ip:port is optional, which defines ip and port that other
peers to connect to.
`
}

func (a *add) Exec(args []string) {
	if len(args) < 1 {
		fmt.Println(a.Help())
		return
	}

	name := strings.TrimSpace(args[0])
	if len(name) == 0 {
		fmt.Println(a.Help())
		return
	}

	cfg, err := db.Load()
	if err != nil {
		fmt.Println("cannot load data: ", err)
		return
	}

	var ip string
	var port uint
	if len(args) >= 2 {
		arr := strings.Split(args[1], ":")
		if len(arr) != 2 {
			fmt.Println(a.Help())
			return
		}
		p, err := strconv.ParseUint(arr[1], 10, 32)
		if err != nil {
			fmt.Println(a.Help())
			return
		}
		ip = arr[0]
		port = uint(p)
	}

	if err = cfg.AddPeer(name, ip, port); err != nil {
		fmt.Println("Cannot create new peer: ", err)
		return
	}

	if err = db.Save(cfg); err != nil {
		fmt.Println("Cannot save config: ", err)
		return
	}

	fmt.Println("New client " + name + " added.")
}

func init() { register("add", "add a new peer", &add{}) }
