package iplookup

import "net"

func mask[T NetAddr | NetAddr6](a, b []byte) T {
	var r T
	for i := 0; i < len(r) && i < len(a) && i < len(b); i++ {
		r[i] = a[i] & b[i]
	}
	return r
}

func NewCIDRLookupTable[T NetAddr | NetAddr6]() *CIDRLookupTable[T] {
	var r CIDRLookupTable[T]
	r.networks = make(map[T]string)
	r.netmasks.set = make(map[T]struct{})
	return &r
}

func (t *CIDRLookupTable[T]) AddNetwork(network net.IPNet, s string) {
	var netaddr T
	var addrslice []byte
	if len(netaddr) == net.IPv4len {
		addrslice = network.IP.To4()
	}
	if len(netaddr) == net.IPv6len {
		addrslice = network.IP.To16()
	}
	if addrslice == nil {
		return
	}
	netaddr = (T)(addrslice)

	t.networks[netaddr] = s
	t.netmasks.Add(network.Mask)
}

func (t *CIDRLookupTable[T]) Lookup(ip net.IP) (string, bool) {
	for _, m := range t.netmasks.masks {
		var net2lookup T
		var addrslice []byte

		if len(net2lookup) == net.IPv4len {
			addrslice = ip.To4()
		}
		if len(net2lookup) == net.IPv6len {
			addrslice = ip.To16()
		}
		if addrslice == nil {
			return "<NEITHER v4 NOR v6>", false
		}
		net2lookup = mask[T](addrslice, m)
		if r, ok := t.networks[net2lookup]; ok {
			return r, true
		}
	}

	return "<NOT FOUND>", false
}

func (t *CIDRLookupTable[T]) Size() int {
	return len(t.networks)
}

func DualLookup(t4 *CIDRLookupTable[NetAddr], t6 *CIDRLookupTable[NetAddr6], ip net.IP) (string, bool) {
	addrslice := ip.To4()
	if addrslice != nil {
		return t4.Lookup(ip)
	}
	return t6.Lookup(ip)
}

func NewDualLookupTable() *DualLookupTable {
	var r DualLookupTable
	r.v4 = NewCIDRLookupTable[NetAddr]()
	r.v6 = NewCIDRLookupTable[NetAddr6]()
	return &r
}

func (dt *DualLookupTable) Lookup(ip net.IP) (string, bool) {
	return DualLookup(dt.v4, dt.v6, ip)
}

func (dt *DualLookupTable) AddNetwork(network net.IPNet, s string) {
	if network.IP.To4() != nil {
		dt.v4.AddNetwork(network, s)
	} else {
		dt.v4.AddNetwork(network, s)
	}
}

func (dt *DualLookupTable) Size() int {
	return dt.v4.Size() + dt.v6.Size()
}
