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
	$(GOBUILD) ./...
	$(GOBUILD) -o ../../../build/$(BINARY_NAME)
	#cd ../../..
package:
	cd ../../.. && zip ./build/main.zip ./build/$(BINARY_NAME)
#test:#
#	$(GOTEST) -v ./sample
clean: 
	rm -f ../../../build/$(BINARY_NAME)
	rm -f ..//../../build/$(BINARY_NAME).zip

check-swagger:
	which swagger || (go get -u github.com/go-swagger/go-swagger/cmd/swagger)

swagger: check-swagger
	go mod vendor  && swagger generate spec -o ./swagger.yaml --scan-models

serve-swagger: check-swagger
	swagger serve -F=swagger swagger.yaml --host=localhost --port 8282

swagger-generate: 
	make swagger && make serve-swagger