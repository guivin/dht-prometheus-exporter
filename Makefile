GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
DEP=dep
BINARY_NAME=dht-prometheus-exporter
BINARY_DEST=/usr/bin

all:
	make dep
	make build
	make install
	make clean

build:
	$(GOBUILD) -o ${BINARY_NAME} -v

clean:
	${GOCLEAN}
	rm -f ${BINARY_NAME}

install:
	sudo cp -f ${BINARY_NAME} ${BINARY_DEST}

uninstall:
	sudo rm -f ${BINARY_DEST}/${BINARY_NAME}

dep:
	${DEP} ensure -v
