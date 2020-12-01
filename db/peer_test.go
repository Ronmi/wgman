/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package db

import (
	"net"
	"testing"
)

type PeerTestSuite struct {
	ipnet net.IPNet
}

func (s *PeerTestSuite) genInterfaceSimple(t *testing.T) {
	p := &Peer{
		PrivKey: "priv",
		PubKey:  "pub",
		Port:    0,
		IP:      5,
		ExtIP:   "",
		Inners:  map[string]InnerPeer{},
	}
	const expect = `[Interface]
PrivateKey = priv
Address = 1.2.3.5/24
PostUp = iptables -A FORWARD -i %i -j ACCEPT
PreDown = iptables -D FORWARD -i %i -j ACCEPT
`

	actual := p.GenInterface(s.ipnet)

	t.Log("result:\n" + actual)

	if expect != actual {
		t.Log("unexpected result:")
		t.Fatal(actual)
	}
}

func (s *PeerTestSuite) genPeerSimple(t *testing.T) {
	p := &Peer{
		PrivKey: "priv",
		PubKey:  "pub",
		Port:    0,
		IP:      5,
		ExtIP:   "",
		Inners:  map[string]InnerPeer{},
	}
	const expect = `[Peer]
PublicKey = pub
AllowedIPs = 1.2.3.5/32
`

	actual := p.GenPeer(s.ipnet)

	t.Log("result:\n" + actual)

	if expect != actual {
		t.Log("unexpected result:")
		t.Fatal(actual)
	}
}

func (s *PeerTestSuite) genInterfaceExtern(t *testing.T) {
	p := &Peer{
		PrivKey: "priv",
		PubKey:  "pub",
		Port:    1234,
		IP:      5,
		ExtIP:   "5.6.7.8",
		Inners:  map[string]InnerPeer{},
	}
	const expect = `[Interface]
PrivateKey = priv
Address = 1.2.3.5/24
PostUp = iptables -A FORWARD -i %i -j ACCEPT
PreDown = iptables -D FORWARD -i %i -j ACCEPT
ListenPort = 1234
`

	actual := p.GenInterface(s.ipnet)

	t.Log("result:\n" + actual)

	if expect != actual {
		t.Log("unexpected result:")
		t.Fatal(actual)
	}
}

func (s *PeerTestSuite) genPeerExtern(t *testing.T) {
	p := &Peer{
		PrivKey: "priv",
		PubKey:  "pub",
		Port:    1234,
		IP:      5,
		ExtIP:   "5.6.7.8",
		Inners:  map[string]InnerPeer{},
	}
	const expect = `[Peer]
PublicKey = pub
AllowedIPs = 1.2.3.5/32
Endpoint = 5.6.7.8:1234
`

	actual := p.GenPeer(s.ipnet)

	t.Log("result:\n" + actual)

	if expect != actual {
		t.Log("unexpected result:")
		t.Fatal(actual)
	}
}

func (s *PeerTestSuite) genInterfaceInner(t *testing.T) {
	p := &Peer{
		PrivKey: "priv",
		PubKey:  "pub",
		Port:    1234,
		IP:      5,
		ExtIP:   "5.6.7.8",
		Inners: map[string]InnerPeer{
			"inner": {
				RealIP: "10.1.1.1",
				IP:     6,
			},
		},
	}
	const expect = `[Interface]
PrivateKey = priv
Address = 1.2.3.5/24
PostUp = iptables -A FORWARD -i %i -j ACCEPT
PreDown = iptables -D FORWARD -i %i -j ACCEPT
ListenPort = 1234
#### inner peers
PostUp = echo 1 > /proc/sys/net/ipv4/ip_forward
#### inner
PostUp = iptables -t nat -A PREROUTING -i %i -d 1.2.3.6 -j DNAT --to 10.1.1.1
PostUp = iptables -t nat -A POSTROUTING -s 1.2.3.6/24 -d 10.1.1.1 -j MASQUERADE
PreDown = iptables -t nat -D PREROUTING -i %i -d 1.2.3.6 -j DNAT --to 10.1.1.1
PreDown = iptables -t nat -D POSTROUTING -s 1.2.3.6/24 -d 10.1.1.1 -j MASQUERADE
`

	actual := p.GenInterface(s.ipnet)

	t.Log("result:\n" + actual)

	if expect != actual {
		t.Log("unexpected result:")
		t.Fatal(actual)
	}
}

func (s *PeerTestSuite) genPeerInner(t *testing.T) {
	p := &Peer{
		PrivKey: "priv",
		PubKey:  "pub",
		Port:    1234,
		IP:      5,
		ExtIP:   "5.6.7.8",
		Inners: map[string]InnerPeer{
			"inner": {
				RealIP: "10.1.1.1",
				IP:     6,
			},
		},
	}
	const expect = `[Peer]
PublicKey = pub
AllowedIPs = 1.2.3.5/32
Endpoint = 5.6.7.8:1234
AllowedIPs = 1.2.3.6/32 #inner
`

	actual := p.GenPeer(s.ipnet)

	t.Log("result:\n" + actual)

	if expect != actual {
		t.Log("unexpected result:")
		t.Fatal(actual)
	}
}

func TestPeerGen(t *testing.T) {
	suite := &PeerTestSuite{
		ipnet: net.IPNet{
			IP:   []byte{1, 2, 3, 0},
			Mask: net.CIDRMask(24, 32),
		},
	}
	t.Run("interface-simple", suite.genInterfaceSimple)
	t.Run("peer-simple", suite.genPeerSimple)
	t.Run("interface-external", suite.genInterfaceExtern)
	t.Run("peer-external", suite.genPeerExtern)
	t.Run("interface-inner", suite.genInterfaceInner)
	t.Run("peer-inner", suite.genPeerInner)
}
