package iplookup

import (
	"net"
	"sort"
)

func (n *Netmasks[T]) Len() int {
	return len(n.masks)
}

func (n *Netmasks[T]) Less(i, j int) bool {
	is, _ := n.masks[i].Size()
	js, _ := n.masks[j].Size()
	return is < js
}

func (n *Netmasks[T]) Swap(i, j int) {
	n.masks[i], n.masks[j] = n.masks[j], n.masks[i]
}

func (n *Netmasks[T]) Add(ipmask net.IPMask) {
	if _, ok := n.set[T(ipmask)]; !ok {
		n.masks = append(n.masks, ipmask)
		sort.Sort(n)
		n.set[T(ipmask)] = struct{}{}
	}
}
