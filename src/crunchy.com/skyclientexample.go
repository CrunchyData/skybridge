package main

import (
	"encoding/json"
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	"github.com/skynetservices/skydns/msg"
)

func main() {

	fmt.Println("here")
	ttl := uint32(360)

	var services = []*msg.Service{
		{Host: "10.0.0.1", Key: "doo.crunchy.lab."},
		{Host: "doo.crunchy.lab", Key: "1.0.0.10.in-addr.arpa."},
	}

	//client := etcd.NewClient([]string{"http://192.168.0.106:4001"})
	client := etcd.NewClient([]string{"http://127.0.0.1:4001"})

	//add a service

	serv := services[0]
	b, err := json.Marshal(serv)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	path, _ := msg.PathWithWildcard(serv.Key)

	_, err = client.Create(path, string(b), uint64(ttl))
	if err != nil {
		// TODO(miek): allow for existing keys...
		fmt.Println(err.Error())
	}

	//add a PTR
	fmt.Println("creating PTR...")
	serv = services[1]
	b, err = json.Marshal(serv)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	path, _ = msg.PathWithWildcard(serv.Key)

	_, err = client.Create(path, string(b), uint64(ttl))
	if err != nil {
		// TODO(miek): allow for existing keys...
		fmt.Println(err.Error())
	}

	//delete a service

	serv = services[0]
	path, _ = msg.PathWithWildcard(serv.Key)

	_, err = client.Delete(path, false)
	if err != nil {
		fmt.Println(err.Error())
	}

	//delete a PTR

	serv = services[1]
	path, _ = msg.PathWithWildcard(serv.Key)

	_, err = client.Delete(path, false)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("done")
}
