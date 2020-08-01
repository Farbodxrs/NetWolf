package main

import (
	"./node"
	"flag"
	"fmt"
)

func main() {

	name := flag.String("name", "", "new node name")

	ip := flag.String("ip", "", "node ip")
	dpn := flag.Int("dpn", 0, "discovery port number")
	dir := flag.String("dir", "", "directory to serve")

	flag.Parse()

	n := node.New(name, ip, dpn, dir, flag.Args())

	fmt.Printf("%+v", *n)

}
