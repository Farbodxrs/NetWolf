package node

import (
	"strconv"
	"strings"
)

type Node struct {
	Address      Address
	Cluster      []Address
	Directory    string
	TransferPort int
}
type Address struct {
	Name          string
	Ip            string
	DiscoveryPort int
}

func New(name *string, ip *string, dpn *int, dir *string, cluster []string) *Node {

	//if _, err := os.Stat(*dir); os.IsNotExist(err) {
	//	if e := os.Mkdir(*dir, 0777); e != nil {
	//		logrus.Fatal("cannot make dir")
	//	}
	//} //mkdir

	n := &Node{
		Address: Address{
			Name:          *name,
			Ip:            *ip,
			DiscoveryPort: *dpn,
		},
		Cluster:      nil,
		Directory:    *dir,
		TransferPort: 0,
	}

	for _, v := range cluster {
		tmp := strings.Split(v, ":")
		port, _ := strconv.Atoi(tmp[2])
		n.Cluster = append(n.Cluster, Address{
			Name:          tmp[0],
			Ip:            tmp[1],
			DiscoveryPort: port,
		})
	}
	//validate n

	return n

}
