package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

var version string = "v0.0.0-dev"

func stress(network string, addr string, stop int) (time.Duration, int, int) {
	
	conn, err := net.Dial(network, addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	w := bufio.NewWriter(conn)
	r := bufio.NewReader(conn)

	var counter int = 0
	var hits int = 0
	var misses int = 0

	var ip net.IP = make([]byte, 4)
	start := time.Now()
	for {
		rand.Read(ip)
		w.WriteString(ip.String())
		w.WriteString("\n")
		w.Flush()
		l, _, err := r.ReadLine()
		if err != nil {
			panic(err)
		}
		counter++

		if len(l) == 0 {
			misses++
		} else {
			hits++
		}
		if stop == counter {
			break
		}
	}
	return time.Since(start), hits, misses
}

func main() {
	var pFlag = flag.Int("p", runtime.NumCPU()/2, "Count of parallel worker routines to send queries")
	var cFlag = flag.Int("c", 100000, "Amount of quries to do with each worker routine")
	var versionFlag = flag.Bool("version", false, "Prints version and exists")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] [PROTO]:[ADDR]\n\nFlags:\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if *versionFlag {
		fmt.Println(version)
		os.Exit(0)
	}

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	connect := strings.SplitN(flag.Arg(0), ":", 2)
	if len(connect) != 2 || !(connect[0] == "unix" || connect[0] == "tcp") {
		fmt.Printf("Invalid addr format: %v\n", flag.Arg(0))
		os.Exit(1)
	}

	var wg sync.WaitGroup
	for i := 0; i < *pFlag; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			du, hits, misses := stress(connect[0], connect[1], *cFlag)
			rate := float64(hits+misses) / du.Seconds()
			fmt.Printf("[%v] %v hit/miss(%v/%v) lookups took %v (%.2f lookups/s)\n", n, hits+misses, hits, misses, du, rate)
		}(i)
	}
	wg.Wait()

}
