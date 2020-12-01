/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package cmds

import (
	"fmt"
	"log"
	"strings"
)

type Cmd interface {
	Exec(args []string)
	Help() (ret string)
}

type data struct {
	Cmd
	Desc string
}

var Cmds = map[string]data{}

func register(n, desc string, c Cmd) {
	if c, ok := Cmds[n]; ok {
		log.Fatalf("command %s already exists with %T", n, c)
	}

	Cmds[n] = data{
		Cmd:  c,
		Desc: desc,
	}
}

func Run(args []string) {
	if len(args) < 1 {
		ListCommands()
		return
	}

	c, ok := Cmds[strings.TrimSpace(args[0])]
	if !ok {
		ListCommands()
		return
	}

	c.Exec(args[1:])
}

func ListCommands() {

	fmt.Print(`
wgman is a simple tool to manage wireguard configurations with wireguard binary
"wg". It saves peers in config file ("data.json" by default) and generate wg-quick
compitable config files.

Use environmental variable "WGMAN_CONFIG" to specify a different config file.

Available commands:
`)

	for n, v := range Cmds {
		fmt.Printf("   %-10s %s\n", n, v.Desc)
	}
}
