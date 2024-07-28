package iplookup

import (
	"github.com/phemmer/go-iptrie"
	"net"
	"net/netip"
)

func NewTrieLookupTable() *TrieLookupTable {
	var t TrieLookupTable
	t.size = 0
	t.trie = iptrie.NewTrie()
	return &t
}

func (t *TrieLookupTable) AddNetwork(network net.IPNet, s string) {
	addr, _ := netip.AddrFromSlice(network.IP)
	bits, _ := network.Mask.Size()
	prefix := netip.PrefixFrom(addr, bits)
	t.trie.Insert(prefix, s)
	t.size += 1
}

func (t *TrieLookupTable) Lookup(ip net.IP) (string, bool) {
	ipa, _ := netip.AddrFromSlice(ip)
	r := t.trie.Find(ipa)

	if r == nil {
		return "", false
	}

	return r.(string), true
}

func (t *TrieLookupTable) Size() int {
	return t.size
}
