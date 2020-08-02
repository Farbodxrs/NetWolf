package main

import (
	"./config"
	"./node"
	"flag"
	"fmt"
)


func main() {
	config.Init()
	name := flag.String("name", "", "new node name")

	ip := flag.String("ip", "", "node ip")
	dpn := flag.Int("dpn", 0, "discovery port number")
	dir := flag.String("dir", "", "directory to serve")

	flag.Parse()

	n := node.New(name, ip, dpn, dir, flag.Args())

	go n.DiscoveryClientBegin()
	 n.DiscoveryServerBegin()


	fmt.Printf("%+v", *n)

}

