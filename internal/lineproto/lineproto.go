package lineproto

import (
	"bufio"
	"github.com/securitym0nkey/icalm/pkg/iplookup"
	"io"
	"net"
)

func ServLineProto(r io.Reader, w io.Writer, table iplookup.LookupTable) {
	scanner := bufio.NewScanner(r)
	wr := bufio.NewWriter(w)

	for scanner.Scan() {
		i := scanner.Text()
		ip := net.ParseIP(i)
		if ip != nil {
			s, o := table.Lookup(ip)
			if o {
				wr.WriteString(s)
			}
		}
		wr.WriteString("\n")
		wr.Flush()
	}

}

func NewLineServer(network string, addr string, table iplookup.LookupTable) (net.Listener, error) {
	server, err := net.Listen(network, addr)
	if err != nil {
		return nil, err
	}
	go func(s net.Listener) {
		defer s.Close()
		for {
			conn, err := s.Accept()
			if err != nil {
				break
			}
			go func(c net.Conn) {
				defer c.Close()
				ServLineProto(c, c, table)
			}(conn)
		}
	}(server)
	return server, nil
}
