/*
 Copyright 2015 Crunchy Data Solutions, Inc.
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package main

// skybridge is meant to run on any Docker host that
// needs to register DNS entries for any started or stopped container
// command line options (required)
// -d domain name to use (default: crunchy.lab)
// -s skydns client url (e.g. http://192.168.0.106:4100)
// -h docker socket (default: unix://var/run/docker.sock)
// -t TTL value (default: 36000000 )

import (
	//"errors"
	"flag"
	"fmt"
	dockerapi "github.com/fsouza/go-dockerclient"
	"github.com/golang/glog"
	"strconv"
	"time"
)

var MAX_TRIES = 3

const delaySeconds = 5
const delay = (delaySeconds * 1000) * time.Millisecond

var DOMAIN string
var ETCD string
var DOCKER_HOST string
var TTL uint64

func init() {
	flag.StringVar(&DOMAIN, "d", "crunchy.lab", "domain name to use when creating DNS entries example crunchy.lab")
	flag.StringVar(&ETCD, "s", "http://127.0.0.1:4001", "URL of etcd client example http://192.168.0.106:4001")
	flag.StringVar(&DOCKER_HOST, "h", "unix:///var/run/docker.sock", "docker socket url")
	flag.Uint64Var(&TTL, "t", 36000000, "dns entries ttl value")
	flag.Parse()
}

func main() {

	var dockerConnected = false
	glog.Infoln("DOCKER_HOST=" + DOCKER_HOST)
	glog.Infoln("ETCD=" + ETCD)
	glog.Infoln("TTL=" + strconv.FormatUint(TTL, 10))
	glog.Infoln("DOMAIN=" + DOMAIN)
	var tries = 0
	var docker *dockerapi.Client
	var err error
	for tries = 0; tries < MAX_TRIES; tries++ {
		docker, err = dockerapi.NewClient(DOCKER_HOST)
		err = docker.Ping()
		if err != nil {
			glog.Errorln("could not ping docker host")
			glog.Errorln("sleeping and will retry in %d sec\n", delaySeconds)
			time.Sleep(delay)
		} else {
			glog.Errorln("no err in connecting to docker")
			dockerConnected = true
			break
		}
	}

	if dockerConnected == false {
		glog.Errorln("failing, could not connect to docker after retries")
		glog.Flush()
		panic("cant connect to docker")
	}

	events := make(chan *dockerapi.APIEvents)
	assert(docker.AddEventListener(events))
	glog.Infoln("skybridge: Listening for Docker events...")
	for msg := range events {
		switch msg.Status {
		//case "start", "create":
		case "start":
			glog.Infoln("event: " + msg.Status + " ID=" + msg.ID + " From:" + msg.From)
			Action(msg.Status, msg.ID, docker)
		case "stop":
			glog.Infoln("event: " + msg.Status + " ID=" + msg.ID + " From:" + msg.From)
			Action(msg.Status, msg.ID, docker)
		case "destroy":
			glog.Infoln("event: " + msg.Status + " ID=" + msg.ID + " From:" + msg.From)
			Action(msg.Status, msg.ID, docker)
		case "die":
			glog.Infoln("event: " + msg.Status + " ID=" + msg.ID + " From:" + msg.From)
		default:
			glog.Infoln("event: " + msg.Status)
		}
	}

}

// Action makes a skydns REST API call based on the docker event
func Action(action string, containerId string, docker *dockerapi.Client) {

	//if we fail on inspection, that is ok because we might
	//be checking for a crufty container that no longer exists
	//due to docker being shutdown uncleanly

	container, dockerErr := docker.InspectContainer(containerId)
	if dockerErr != nil {
		fmt.Printf("skybridge: unable to inspect container:%s %s", containerId, dockerErr)
		return
	}
	var hostname = container.Name[1:] + "." + DOMAIN
	var ipaddress = container.NetworkSettings.IPAddress

	if ipaddress == "" {
		glog.Infoln("no ipaddress returned for container: " + hostname)
		return
	}

	switch action {
	case "start":
		glog.Infoln("new container name=" + container.Name[1:] + " ip:" + ipaddress)
		addEntry(hostname, ipaddress)
	case "stop":
		glog.Infoln("removing container name=" + container.Name[1:] + " ip:" + ipaddress)
		deleteEntry(hostname, ipaddress)
	case "destroy":
		glog.Infoln("removing container name=" + container.Name[1:] + " ip:" + ipaddress)
		deleteEntry(hostname, ipaddress)
	default:
	}

}

func assert(err error) {
	if err != nil {
		fmt.Println("skybridge: ", err)
		glog.Flush()
		panic("can't continue")
	}
}