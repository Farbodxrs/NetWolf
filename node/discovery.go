package node

import (
	"../config"
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"strconv"
	"strings"
	"time"
)

func (n *Node) DiscoveryClientBegin() {

	logrus.Info("Begin discovery client for node : ", n.Address.Name)
	for true {
		time.Sleep(config.Cfg.DiscoveryDelay)
		logrus.Info("cluster list is : ", n.Cluster)

		message := ClusterToString(n.Cluster)
		message += n.Address.Name + "@" + n.Address.Ip + "@" + strconv.Itoa(n.Address.DiscoveryPort)
		logrus.Info("Going to send every one : ", message)

		for _, v := range n.Cluster {
			add := getAddress(v)
			conn, err := net.Dial("udp", add)
			if err != nil {
				logrus.Warn("Cannot dial on : ", add)
			}
			fmt.Fprintf(conn, message)
			err = conn.Close()
			if err != nil {
				logrus.Warn("cannot close udp conn on err: , ", err)
			}
		}

	}

}

func (n *Node) DiscoveryServerBegin() {
	url := ":" + strconv.Itoa(n.Address.DiscoveryPort) // :port

	logrus.Info("node : ", n.Address.Name, " serving on port : ", n.Address.DiscoveryPort)

	addr, _ := net.ResolveUDPAddr("udp", url)
	msg := make([]byte, 8192)
	conn, _ := net.ListenUDP("udp", addr)
	for {
		size, _, _ := conn.ReadFromUDP(msg)
		tmp := make([]byte, size)
		copy(tmp, msg)
		go n.handleIncomingUDP(tmp)
	}
}

func (n *Node) handleIncomingUDP(a []byte) {

	rows := strings.Split(string(a), "#")

	logrus.Info("Received : ", rows)
	for _, row := range rows {
		logrus.Info("LETS DO ROW : ", row)
		n.Mutex.Lock()
		parts := strings.Split(row, "@")
		if len(parts) != 3 {
			logrus.Warn("Invalid input", row)
			continue
		}
		if ExistsName(parts[0], n.Cluster) {
			logrus.Info("Exist name skipping ", parts[0])
			continue
		}
		port, e := strconv.Atoi(strings.TrimSuffix(parts[2], "\n"))
		if e != nil {
			logrus.Warn("Invalid port : ", e)
			continue
		}
		if ExistsAddress(parts[1], port, n.Cluster) {
			logrus.Info("Exists Address skipping : ", parts[1], port)
			continue
		}
		if parts[0] == n.Address.Name {
			logrus.Info("Do not add yourself (on name) to list skipping : ")
			//continue
		}
		if parts[1] == n.Address.Ip && port == n.Address.DiscoveryPort {
			logrus.Info("Do not add yourself (on address) to list skipping : ")
			//continue

		}

		logrus.Info("Add node to cluster list", parts)
		n.Cluster = append(n.Cluster, Address{
			Name:          parts[0],
			Ip:            parts[1],
			DiscoveryPort: port,
		})

		n.Mutex.Unlock()
	}
}

func ExistsName(n string, a []Address) bool {
	for _, v := range a {
		if v.Name == n {
			return true
		}
	}
	return false
}
func ExistsAddress(ip string, p int, a []Address) bool {
	for _, v := range a {
		if ip == v.Ip && p == v.DiscoveryPort {
			return true
		}
	}
	return false
}
func NodeToString(a Address) string {
	return a.Name + "@" + a.Ip + "@" + strconv.Itoa(a.DiscoveryPort)
}
func ClusterToString(a []Address) string {
	var res string
	for _, v := range a {
		res += NodeToString(v)
		res += "#"
	}

	return res
}

func getAddress(a Address) string {
	return a.Ip + ":" + strconv.Itoa(a.DiscoveryPort)

}
