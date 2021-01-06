/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package db

import (
	"bytes"
	"fmt"
	"net"
	"os/exec"
	"strings"
)

func genPriv() (ret string, err error) {
	c := exec.Command("wg", "genkey")
	buf, err := c.Output()
	if err == nil {
		ret = strings.TrimSpace(string(buf))
	}
	return
}

func genPub(priv string) (ret string, err error) {
	c := exec.Command("wg", "pubkey")
	c.Stdin = strings.NewReader(priv)

	buf, err := c.Output()
	if err == nil {
		ret = strings.TrimSpace(string(buf))
	}
	return
}

type InnerPeer struct {
	RealIP string
	IP     byte
}

func (p *InnerPeer) safe() {}

type Peer struct {
	PrivKey string
	PubKey  string
	Port    uint
	IP      byte
	ExtIP   string
	Inners  map[string]*InnerPeer
}

func NewPeer(ip byte, port uint, ext string) (ret *Peer, err error) {
	priv, err := genPriv()
	if err != nil {
		return
	}
	pub, err := genPub(priv)
	if err != nil {
		return
	}

	ret = &Peer{
		PrivKey: priv,
		PubKey:  pub,
		Port:    port,
		IP:      ip,
		ExtIP:   ext,
		Inners:  map[string]*InnerPeer{},
	}
	return
}

func (p *Peer) safe() {
	if p.Inners == nil {
		p.Inners = map[string]*InnerPeer{}
	}

	for _, i := range p.Inners {
		i.safe()
	}
}

func (p *Peer) AddInner(name, innerIP string, ip byte) {
	p.Inners[name] = &InnerPeer{
		RealIP: innerIP,
		IP:     ip,
	}
}

func (p *Peer) Endpoint() (ret string) {
	if p.ExtIP == "" || p.Port == 0 {
		return
	}

	return fmt.Sprintf("%s:%d", p.ExtIP, p.Port)
}

func (p *Peer) GenInterface(ip net.IPNet) (ret string) {
	buf := &bytes.Buffer{}
	myip := net.IP(append([]byte{}, ip.IP...)).To4()
	myip[3] = p.IP

	fmt.Fprintf(buf, `[Interface]
PrivateKey = %s
Address = %s/24
PostUp = iptables -A FORWARD -i %%i -j ACCEPT
PreDown = iptables -D FORWARD -i %%i -j ACCEPT
`, p.PrivKey, myip.String(),
	)
	if p.Port > 0 {
		fmt.Fprintf(buf, "ListenPort = %d\n", p.Port)
	}

	if len(p.Inners) > 0 {
		buf.WriteString("#### inner peers\nPostUp = echo 1 > /proc/sys/net/ipv4/ip_forward\n")
	}
	for name, i := range p.Inners {
		myip := net.IP(append([]byte{}, ip.IP...)).To4()
		myip[3] = i.IP
		fmt.Fprintf(buf, "#### %s\n", name)
		fmt.Fprintf(buf,
			"PostUp = iptables -t nat -A PREROUTING -i %%i -d %s -j DNAT --to %s\n",
			myip, i.RealIP,
		)
		fmt.Fprintf(buf,
			"PostUp = iptables -t nat -A POSTROUTING -s %s/24 -d %s -j MASQUERADE\n",
			myip, i.RealIP,
		)
		fmt.Fprintf(buf,
			"PreDown = iptables -t nat -D PREROUTING -i %%i -d %s -j DNAT --to %s\n",
			myip, i.RealIP,
		)
		fmt.Fprintf(buf,
			"PreDown = iptables -t nat -D POSTROUTING -s %s/24 -d %s -j MASQUERADE\n",
			myip, i.RealIP,
		)
		fmt.Fprintf(buf,
			`#### To connect %s into this virtual network, add a routing rule for %s:
####     sudo route add -net %s/24 gw ip_address_of_this_host
#### and uncomment 2 lines below
`,
			name, name, myip,
		)
		fmt.Fprintf(buf,
			"#PostUp = iptables -t nat -A POSTROUTING -s %s -d %s/24 -j MASQUERADE\n",
			i.RealIP, myip,
		)
		fmt.Fprintf(buf,
			"#PreDown = iptables -t nat -D POSTROUTING -s %s -d %s/24 -j MASQUERADE\n",
			i.RealIP, myip,
		)
	}

	return buf.String()
}

func (p *Peer) GenPeer(ip net.IPNet) (ret string) {
	buf := &bytes.Buffer{}
	myip := net.IP(append([]byte{}, ip.IP...)).To4()
	myip[3] = p.IP

	fmt.Fprintf(buf, `[Peer]
PublicKey = %s
AllowedIPs = %s/32
`, p.PubKey, myip.String(),
	)

	if p.ExtIP != "" {
		fmt.Fprintf(buf, "Endpoint = %s:%d\n", p.ExtIP, p.Port)
	}

	for name, i := range p.Inners {
		myip := net.IP(append([]byte{}, ip.IP...)).To4()
		myip[3] = i.IP
		fmt.Fprintf(buf,
			"AllowedIPs = %s/32 #%s\n",
			myip, name,
		)
	}

	return buf.String()
}
