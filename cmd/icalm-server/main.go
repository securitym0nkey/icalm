package main

import (
	"flag"
	"fmt"
	"github.com/securitym0nkey/icalm/internal/config"
	"github.com/securitym0nkey/icalm/internal/lineproto"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("IP: CIDR annotation lookup microservice")
	var httplistenFlag = flag.String("http-listen", "", "Address on which the server will listen for http requests. Example: 127.0.0.1:8226")
	var linelistenFlag = flag.String("line-listen", "", "Address on which the server will listen for line-protocol requests. Example: 127.0.0.1:4226")
	var lineunixFlag = flag.String("line-unix", "", "Path of a unix-socket on which icalm provides a line-protocol lookup service. Example: /var/run/icalm/sock")
	var cidrfileFlag = flag.String("networks", "./networks.csv", "Path to the file containing the CIDR to annotation mapping.")

	flag.Parse()

	toclose := make([]io.Closer, 0)

	lookuptable, err := config.LoadLookupTableFromFile(*cidrfileFlag)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Loaded lookuptable (%v entries)\n", lookuptable.Size())

	if *httplistenFlag != "" {
		log.Println("HTTP Server NOT implemented yet. :(")
	}

	if *linelistenFlag != "" {
		lineserver, err := lineproto.NewLineServer("tcp", *linelistenFlag, lookuptable)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Listening on %v for lineproto requests \n", lineserver.Addr())
		toclose = append(toclose, lineserver)
	}

	if *lineunixFlag != "" {
		lineunixserver, err := lineproto.NewLineServer("unix", *lineunixFlag, lookuptable)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Serving lineproto @ %v \n", lineunixserver.Addr())
		toclose = append(toclose, lineunixserver)
	}

	if len(toclose) == 0 {
		fmt.Println("Useless. Need to serve on at least one method. Specify -http-listen, -line-listen or -line-unix: ")
		flag.PrintDefaults()
		os.Exit(1)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down")

	for _, cs := range toclose {
		cs.Close()
	}

}
