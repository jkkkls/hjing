GITSHA=`git rev-list HEAD -n 1 | cut -c 1-`
BUILD_TIME=`date +%Y-%m-%d_%H:%M:%S`
PC_NAME=`hostname`

LDFLAGS=-ldflags "-X github.com/jkkkls/hjing/rpc.GitTag=${GITTAG} -X github.com/jkkkls/hjing/rpc.BuildTime=${BUILD_TIME} -X github.com/jkkkls/hjing/rpc.PcName=${PC_NAME} -X github.com/jkkkls/hjing/rpc.GitSHA=${GITSHA} -X github.com/jkkkls/hjing/rpc.Area=${BUILD_AREA}"

GOCMD=go
GOBUILD=$(BUILD_PERFIX) $(GOCMD) build ${LDFLAGS}
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

all: pb build

.PHONY: pb


#pb tag
pb:
	@protoc -I. --gogu_out ./ --gogu_opt paths=source_relative ./pb/*.proto
	@$(GOCMD) fmt ./pb/*.go
	@echo "编协议文件完成"

#build tag
build:
build_all: build
test:
	@$(GOTEST) -v ./...
clean:
	@$(GOCLEAN)
	@mkdir ./build/conf
	@rm -rf ./build/*
