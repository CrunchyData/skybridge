OSFLAVOR=centos7

all:
ifeq ($(OSFLAVOR),centos7)
	cp Dockerfile.centos7 Dockerfile
else
	cp Dockerfile.rhel7 Dockerfile
endif
		rm -rf bin
		mkdir bin
		cp ../../bin/etcd* ./bin
		cp ../../bin/skydns ./bin
		cp ../../bin/skybridge ./bin
		cp ../../bin/start-skybridge.sh ./bin
	        sudo docker build -t skybridge .
		sudo docker tag -f skybridge:latest crunchydata/skybridge
