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

var version string = "v0.0.0-dev"


func main() {
	version = config.VersionString()

	var httplistenFlag = flag.String("http-listen", "", "Address on which the server will listen for http requests. Example: 127.0.0.1:8226")
	var linelistenFlag = flag.String("line-listen", "", "Address on which the server will listen for line-protocol requests. Example: 127.0.0.1:4226")
	var lineunixFlag = flag.String("line-unix", "", "Path of a unix-socket on which icalm provides a line-protocol lookup service. Example: /var/run/icalm/sock")
	var cidrfileFlag = flag.String("networks", "./networks.csv", "Path to the file containing the CIDR to annotation mapping.")
	var versionFlag = flag.Bool("version", false, "Prints version and exists")

	flag.Parse()

	flag.Usage = func(){
		fmt.Fprintf(os.Stderr, "Usage: %v [flags]\n\nFlags:\n", os.Args[0])
		flag.PrintDefaults()
	}

	if *versionFlag {
		fmt.Println(version)
		os.Exit(0)
	}

	fmt.Printf("IP: CIDR annotation lookup microservice [%v]\n", version)

	if *httplistenFlag == "" && *linelistenFlag == "" && *lineunixFlag == "" {
		fmt.Fprintln(os.Stderr,"Useless. Need to serve on at least one method. Specify -http-listen, -line-listen or -line-unix")
		flag.Usage()
		os.Exit(1)
	}

	servers := make([]*lineproto.LineServer,0)
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
		log.Printf("Listening on %v for lineproto requests \n", lineserver.Listener.Addr())
		toclose = append(toclose, lineserver.Listener)
		servers = append(servers, lineserver)
	}

	if *lineunixFlag != "" {
		lineunixserver, err := lineproto.NewLineServer("unix", *lineunixFlag, lookuptable)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Serving lineproto @ %v \n", lineunixserver.Listener.Addr())
		toclose = append(toclose, lineunixserver.Listener)
		servers = append(servers, lineunixserver)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sigHupChan := make(chan os.Signal, 1)
	signal.Notify(sigHupChan, syscall.SIGHUP)


	run := true
	for run {
		select {
			case <-sigChan:
			run = false
			case <-sigHupChan:
				for _, serv := range servers {
					lookuptable, err := config.LoadLookupTableFromFile(*cidrfileFlag)
					if err != nil {
						log.Fatal(err)
					}
					serv.Reload(lookuptable)
					log.Printf("Reloaded lookuptable (%v entries)\n", lookuptable.Size())
				}
		}
	}

	log.Println("Shutting down")

	for _, cs := range toclose {
		cs.Close()
	}

}
