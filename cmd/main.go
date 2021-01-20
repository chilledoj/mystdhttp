package main

import (
	"flag"
	"log"
	"net"
)

func main() {
	var (
		host      = flag.String("host", "", "host http address to listen on")
		port      = flag.String("port", "8000", "port number for http listener")
		initTasks = flag.Bool("init-tasks", true, "set to false to not prepopulate the in-memory db")
	)
	flag.Parse()

	addr := net.JoinHostPort(*host, *port)

	if err := runHttp(options{
		listenAddr: addr,
		initTasks:  *initTasks,
	}); err != nil {
		log.Fatal(err)
	}
}
