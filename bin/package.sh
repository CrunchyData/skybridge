#!/bin/bash -x


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

VERSION=1.0.4
ETCD_DIR=$HOME/etcd-v2.0.0-linux-amd64
SKYDNS_DIR=$HOME/skydns
SKYBRIDGE_DIR=$HOME/skybridge
TMPDIR=/tmp/skybridge.package

mkdir -p $TMPDIR/var/cpm/bin
mkdir -p $TMPDIR/var/cpm/config
mkdir -p $TMPDIR/var/cpm/data/etcd

cp $ETCD_DIR/etcd \
	$ETCD_DIR/etcdctl \
	$SKYDNS_DIR/bin/skydns \
	$SKYBRIDGE_DIR/bin/skybridge \
	$SKYBRIDGE_DIR/bin/wait20 \
      	$TMPDIR/var/cpm/bin

cp $SKYBRIDGE_DIR/bin/install.sh \
	$TMPDIR 

cp $SKYBRIDGE_DIR/config/skybridge.service \
	$SKYBRIDGE_DIR/config/etcd.service \
	$SKYBRIDGE_DIR/config/docker \
	$SKYBRIDGE_DIR/config/skydns.service \
      	$TMPDIR/var/cpm/config

cd $TMPDIR

tar cvzf /tmp/skybridge.$VERSION-linux-amd64.tar.gz .
	
rm -rf $TMPDIR

