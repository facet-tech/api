# Go parameters
GOCMD=go
GOBUILD=GOARCH=amd64 GOOS=linux $(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
CORE_BINARY_NAME=aws-lambda-go-api-proxy-core
GIN_BINARY_NAME=aws-lambda-go-api-proxy-gin
BINARY_NAME=main
    
all: clean build package
build: 
	#cd src/facet.ninja/api
	cd ../../.. && rm -rf ./build/
	$(GOBUILD) ./...
	$(GOBUILD) -o ../../../build/$(BINARY_NAME)
	#cd ../../..
package:
	cd ../../.. && cp -r ./resources ./build/
	cd ../../../build && zip -r ./main.zip ./
#test:#
#	$(GOTEST) -v ./sample
clean: 
	rm -rf ../../../build/