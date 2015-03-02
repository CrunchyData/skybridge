#!/bin/bash 


# Copyright 2015 Crunchy Data Solutions, Inc.
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
# 
# package up etcd, skydns, and skybridge for a user install archive


CPM_DIR=/opt/cpm

sudo mkdir -p $CPM_DIR
sudo mkdir -p $CPM_DIR/bin
sudo mkdir -p $CPM_DIR/data/etcd
sudo mkdir -p $CPM_DIR/config

sudo cp -r ./opt /

echo -n "enter the ipaddress of this server: "
read SERVERIP

echo "configuring etcd to use " $SERVERIP
sed -i "s/192.168.0.106/$SERVERIP/g" ./opt/cpm/config/etcd.service
echo "configuring skydns to use " $SERVERIP
sed -i "s/192.168.0.106/$SERVERIP/g" ./opt/cpm/config/skydns.service
echo "configuring skybridge to use " $SERVERIP
sed -i "s/192.168.0.106/$SERVERIP/g" ./opt/cpm/config/skybridge.service
echo "configuring docker to use " $SERVERIP
sed -i "s/192.168.0.106/$SERVERIP/g" ./opt/cpm/config/docker

echo -n "enter the domain name to use: "
read DOMAINNAME
echo "configuring skydns to use " $DOMAINNAME
sed -i "s/crunchy.lab/$DOMAINNAME/g" ./opt/cpm/config/skydns.service
echo "configuring skybridge to use " $DOMAINNAME
sed -i "s/crunchy.lab/$DOMAINNAME/g" ./opt/cpm/config/skybridge.service
echo "configuring docker to use " $DOMAINNAME
sed -i "s/crunchy.lab/$DOMAINNAME/g" ./opt/cpm/config/docker

SYSTEMD=/usr/lib/systemd/system
sudo cp ./opt/cpm/config/etcd.service $SYSTEMD
sudo cp ./opt/cpm/config/skydns.service $SYSTEMD
sudo cp ./opt/cpm/config/skybridge.service $SYSTEMD

sudo yum -y install docker-io
sudo usermod -a -G docker $USER
sudo cp ./opt/cpm/config/docker /etc/sysconfig/
sudo systemctl enable docker.service
sudo systemctl start docker.service

sudo systemctl enable etcd.service
sudo systemctl start etcd.service
sudo systemctl enable skydns.service
sudo systemctl start skydns.service
sudo systemctl enable skybridge.service
sudo systemctl start skybridge.service

echo "installation done...etcd, skydns, and skybridge should be running"
