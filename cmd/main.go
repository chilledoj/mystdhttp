package main

import (
	"flag"
	"log"
	"net"
	"os"
)

func main() {
	var (
		host      = flag.String("host", "", "host http address to listen on")
		port      = flag.String("port", "8000", "port number for http listener")
		initTasks = flag.Bool("init-tasks", true, "set to false to not prepopulate the in-memory db")
	)
	flag.Parse()

	lg := log.New(os.Stdout, "[stask] ", log.LstdFlags)

	addr := net.JoinHostPort(*host, *port)

	if err := runHttp(options{
		lg:         lg,
		listenAddr: addr,
		initTasks:  *initTasks,
	}); err != nil {
		log.Fatal(err)
	}
}
