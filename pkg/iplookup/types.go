package iplookup

import (
	"github.com/phemmer/go-iptrie"
	"net"
)

type LookupTable interface {
	Lookup(net.IP) (string, bool)
	AddNetwork(net.IPNet, string)
	Size() int
}

type NetAddr [net.IPv4len]byte
type NetAddr6 [net.IPv6len]byte

type Netmasks[T NetAddr | NetAddr6] struct {
	set   map[T]struct{}
	masks []net.IPMask
}

type CIDRLookupTable[T NetAddr | NetAddr6] struct {
	netmasks Netmasks[T]
	networks map[T]string
}

type DualLookupTable struct {
	v4 *CIDRLookupTable[NetAddr]
	v6 *CIDRLookupTable[NetAddr6]
}

type TrieLookupTable struct {
	size int
	trie *iptrie.Trie
}
