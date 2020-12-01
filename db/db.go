/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package db

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"
	"os"
)

var ConfigFilename = "data.json"

func Save(n *Network) (err error) {
	buf, err := json.Marshal(n)
	if err != nil {
		return
	}

	dst := &bytes.Buffer{}
	if err = json.Indent(dst, buf, "", "    "); err != nil {
		return
	}

	err = ioutil.WriteFile(ConfigFilename, dst.Bytes(), 0600)
	return
}

func Load() (ret *Network, err error) {
	f, err := os.Open(ConfigFilename)
	if err != nil {
		return
	}
	defer f.Close()

	var x Network
	dec := json.NewDecoder(f)
	if err = dec.Decode(&x); err == nil {
		ret = &x
	}

	x.safe()

	return
}

type Network struct {
	Net   net.IPNet
	Peers map[string]*Peer
}

func (n *Network) safe() {
	if n.Peers == nil {
		n.Peers = map[string]*Peer{}
	}
	for _, p := range n.Peers {
		p.safe()
	}
}

func (n *Network) IP(ip byte) (ret net.IP) {
	ret = net.IP(append([]byte{}, n.Net.IP...)).To4()
	ret[3] = ip
	return
}

func (n *Network) genIP() (ret byte) {
	okips := map[byte]bool{}
	for i := byte(1); i < 255; i++ {
		okips[i] = true
	}
	for _, p := range n.Peers {
		okips[p.IP] = false
	}
	for i, v := range okips {
		if !v {
			continue
		}
		return i
	}

	return
}

func (n *Network) AddInner(parent, name string, ip string) (err error) {
	p, ok := n.Peers[parent]
	if !ok {
		return errors.New("peer " + parent + " not found")
	}

	p.AddInner(name, ip, n.genIP())
	return
}

func (n *Network) AddPeer(name string, ext string, port uint) (err error) {
	ip := n.genIP()

	p, err := NewPeer(ip, port, ext)
	if err == nil {
		n.Peers[name] = p
	}
	return
}

func (n *Network) GenConf(name string) (ret string, err error) {
	me, ok := n.Peers[name]
	if !ok {
		err = errors.New("peer " + name + " not found")
		return
	}

	ret = "# " + name + "\n" + me.GenInterface(n.Net)

	for k, p := range n.Peers {
		if k == name {
			continue
		}

		if me.Port == 0 && p.Port == 0 {
			continue
		}

		ret += "\n# " + k + "\n"
		ret += p.GenPeer(n.Net)
	}

	return
}
