package main

import (
	"encoding/json"
	"github.com/coreos/go-etcd/etcd"
	"github.com/golang/glog"
	"github.com/skynetservices/skydns/msg"
	"strings"
)

//global TTL
//global skydns url

//adds a service entry and a PTR entry
func addEntry(hostname string, ip string) {

	glog.Infoln("addEntry called")

	var services = []*msg.Service{
		{Host: ip, Key: hostname + "."},
		{Host: hostname, Key: reverseIP(ip)},
	}

	client := etcd.NewClient([]string{ETCD})

	//delete any existing entries with this name or ip address
	deleteEntry(hostname, ip)

	//add a service

	glog.Infoln("creating A record...")
	serv := services[0]
	b, err := json.Marshal(serv)
	if err != nil {
		glog.Errorln(err.Error())
		return
	}
	path, _ := msg.PathWithWildcard(serv.Key)

	_, err = client.Create(path, string(b), TTL)
	if err != nil {
		// TODO(miek): allow for existing keys...
		glog.Errorln(err.Error())
	}

	//add a PTR
	glog.Infoln("creating PTR record...")
	serv = services[1]
	b, err = json.Marshal(serv)
	if err != nil {
		glog.Errorln(err.Error())
		return
	}
	path, _ = msg.PathWithWildcard(serv.Key)

	_, err = client.Create(path, string(b), TTL)
	if err != nil {
		// TODO(miek): allow for existing keys...
		glog.Errorln(err.Error())
	}

	glog.Infoln("addEntry completed")

}

//delete both the service entry and the PTR entry
func deleteEntry(hostname string, ip string) {
	glog.Infoln("deleteEntry called...")
	var services = []*msg.Service{
		{Host: ip, Key: hostname + "."},
		{Host: hostname, Key: reverseIP(ip)},
	}

	client := etcd.NewClient([]string{ETCD})
	//delete a service

	serv := services[0]
	path, _ := msg.PathWithWildcard(serv.Key)

	_, err := client.Delete(path, false)
	if err != nil {
		glog.Errorln(err.Error())
	}

	//delete a PTR

	serv = services[1]
	path, _ = msg.PathWithWildcard(serv.Key)

	_, err = client.Delete(path, false)
	if err != nil {
		glog.Errorln(err.Error())
	}

	glog.Infoln("deleteEntry completed...")

}

//return the reverse ip
func reverseIP(ip string) string {
	//"1.0.0.10.in-addr.arpa."},
	//assume ip has 4 numbers 1.2.3.4
	glog.Flush()
	arr := strings.Split(ip, ".")
	return arr[3] + "." + arr[2] + "." + arr[1] + "." + arr[0] + ".in-addr.arpa"
}
