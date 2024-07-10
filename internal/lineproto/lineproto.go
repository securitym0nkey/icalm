package lineproto

import (
	"bufio"
	"github.com/securitym0nkey/icalm/pkg/iplookup"
	"io"
	"net"
	"log"
	"sync"
)


type LineServerConnection struct {
	conn net.Conn
	update_chan chan *iplookup.LookupTable
}

type LineServer struct {
	Listener net.Listener
	clients map[net.Conn]*LineServerConnection
	table iplookup.LookupTable
	clientsMutex sync.Mutex
}



func ServLineProto(r io.Reader, w io.Writer, inittable *iplookup.LookupTable, newtable_chan chan *iplookup.LookupTable) {

	quit_chan := make(chan int)
	defer close(quit_chan)

	request_chan := make(chan string)
	defer close(request_chan)


	wr := bufio.NewWriter(w)
	table := *(inittable)


	// wait for new request (lines) in separate goroutine
	go func(){
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			r := scanner.Text()
			request_chan <- r
		}
		quit_chan <- 0
	}()


	// main loop per client
	for {
		select {
			case request := <- request_chan:
				ip := net.ParseIP(request)
				s, o := table.Lookup(ip)
				if o {
					wr.WriteString(s)
				}
				wr.WriteString("\n")
				wr.Flush()
			case newtable := <- newtable_chan:
				table = *(newtable)
			case <-quit_chan:
				return
		}
	}

}

func (s *LineServer) AddClientConnection(c *LineServerConnection) {
	s.clientsMutex.Lock()
	defer s.clientsMutex.Unlock()
	log.Printf("Client %s connected", c.conn.RemoteAddr())
	s.clients[c.conn] = c
}


func (s *LineServer) RemoveClientConnection(c *LineServerConnection) {
	s.clientsMutex.Lock()
	defer s.clientsMutex.Unlock()
	log.Printf("Client %s disconnected", c.conn.RemoteAddr())
	delete(s.clients, c.conn)
}

func (s *LineServer) Reload(newtable iplookup.LookupTable){
	s.clientsMutex.Lock()
	defer s.clientsMutex.Unlock()

	s.table = newtable
	for _, client := range s.clients {
		client.update_chan <- &(newtable)
	}

}

func NewLineServer(network string, addr string, table iplookup.LookupTable) (*LineServer, error) {
	listener, err := net.Listen(network, addr)

	if err != nil {
		return nil, err
	}

	server := &LineServer {
		Listener: listener,
		clients: make(map[net.Conn]*LineServerConnection),
		table: table,
	}

	go func(s *LineServer) {
		defer s.Listener.Close()
		for {
			conn, err := s.Listener.Accept()
			if err != nil {
				break
			}

			client := &LineServerConnection{
				conn: conn,
				update_chan: make(chan *iplookup.LookupTable),
			}
			server.AddClientConnection(client)

			go func(c *LineServerConnection) {
				defer c.conn.Close()
				defer server.RemoveClientConnection(c)
				ServLineProto(c.conn, c.conn, &server.table, c.update_chan)

			}(client)
		}
	}(server)

	return server, nil
}
