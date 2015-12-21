OSFLAVOR=centos7

all:
		rm -rf bin
		mkdir bin
		cp ../../bin/etcd* ./bin
		cp ../../bin/skydns ./bin
		cp ../../bin/skybridge ./bin
		cp ../../bin/start-skybridge.sh ./bin
		sudo docker build -t skybridge -f Dockerfile.$(OSFLAVOR) .
		sudo docker tag -f skybridge:latest crunchydata/skybridge
