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
	$(GOBUILD) ./...
	cd facet.ninja/api/src && $(GOBUILD) -o ../build/$(BINARY_NAME)
	cd ../../../
package:
	zip ./facet.ninja/api/build/main.zip ./facet.ninja/api/build/$(BINARY_NAME)
#test:#
#	$(GOTEST) -v ./sample
clean: 
	rm -f facet.ninja/api/src/$(BINARY_NAME)
	rm -f facet.ninja/api/src/$(BINARY_NAME).zip