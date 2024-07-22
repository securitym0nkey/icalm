package config

import (
	"encoding/csv"
	"github.com/securitym0nkey/icalm/pkg/iplookup"
	"io"
	"log"
	"net"
	"os"
)

// LoadLookupTableFromFile loads a comma seperated CIDR file with exactly 2 cols
// 1st col is the Network in CIDR and 2nd col is the map value
func LoadLookupTableFromFile(path string) (*iplookup.DualLookupTable, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	table := iplookup.NewDualLookupTable()

	r := csv.NewReader(f)
	r.FieldsPerRecord = 2
	for {
		record, err := r.Read()

		// End of file
		if err == io.EOF {
			break
		}

		// Something not well formated, skipping
		if err != nil {
			log.Printf("Invalid format: %s\n", err.Error())
			continue
		}

		// Parse CIDR
		_, cidr, err := net.ParseCIDR(record[0])
		if err != nil {
			line, _ := r.FieldPos(0)
			log.Printf("invalid network in line %d: %v\n", line, record[0])
			continue
		}

		// all nice, insert
		table.AddNetwork(*cidr, record[1])
	}
	_ = f.Close()
	return table, nil
}