package config

import (
	"bufio"
	"github.com/securitym0nkey/icalm/pkg/iplookup"
	"log"
	"net"
	"os"
	"strings"
)

// Loads a comma spearated CIDR file with max 2 cols
// 1st col is the Network in CIDR and 2nd col is the map value
func LoadLookupTableFromFile(path string) (*iplookup.DualLookupTable, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	table := iplookup.NewDualLookupTable()
	sn := bufio.NewScanner(f)
	for sn.Scan() {
		cidr2data := sn.Text()
		sp := strings.SplitN(cidr2data, ",", 2)
		if len(sp) != 2 {
			log.Printf("Invalid format: %v\n", cidr2data)
			continue
		}
		_, cidr, err := net.ParseCIDR(sp[0])
		if err != nil {
			log.Printf("invalid network: %v\n", sp[0])
			continue
		}
		table.AddNetwork(*cidr, sp[1])
	}
	f.Close()
	return table, nil
}
