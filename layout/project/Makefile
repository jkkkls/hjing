GITSHA=`git rev-list HEAD -n 1 | cut -c 1-`
BUILD_TIME=`date +%Y-%m-%d_%H:%M:%S`
PC_NAME=`hostname`

LDFLAGS=-ldflags "-X card_server/apps.GitTag=${GITTAG} -X card_server/apps.BuildTime=${BUILD_TIME} -X card_server/apps.PcName=${PC_NAME} -X card_server/apps.GitSHA=${GITSHA} -X card_server/apps.Area=${BUILD_AREA}"

GOCMD=go
GOBUILD=$(BUILD_PERFIX) $(GOCMD) build ${LDFLAGS}
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

all: pb build

.PHONY:pb 


game:
	@$(GOBUILD) -o build/apps/game card_server/apps/game
	@cp libs/node_config/guangbo/game0.toml ./build/apps/
	@echo "编译game完成"
pb:
	@protoc -I. --gogu_out ./ --gogu_opt paths=source_relative ./pb/*.proto
	@$(GOCMD) fmt ./pb/*.go
	@echo "编协议文件完成"

build: game
build_all: build
test:
	@$(GOTEST) -v ./...
clean:
	@$(GOCLEAN)
	@mkdir ./build/conf
	@rm -rf ./build/*
