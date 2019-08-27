#Go command line
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Paths
CBIN=LambdaCreate/main/bin/
CSRC=LambdaCreate/main
RBIN=LambdaRead/main/bin/
RSRC=LambdaRead/main

# Outputs
COUTPUT= Create
ROUTPUT= Read

.PHONY: Des LambdaCreate LambdaRead
#LambdaCreate:  des test build
LambdaCreate:  des build test

$(COUTPUT): %: $(CSRC)/%.go
	GOOS=linux $(GOBUILD) -o $(CBIN)$@ $<

build: $(COUTPUT)

#LambdaRead:  des test build
LambdaRead: des build test

$(ROUTPUT): %: $(RSRC)/%.go
	GOOS=linux $(GOBUILD) -o $(RBIN)$@ $<

build: $(ROUTPUT)

.PHONY: test
test:
	$(GOTEST) -v ./...

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -rf $(BIN)

.PHONY: des
des:
	$(GOGET) -u github.com/aws/aws-lambda-go/lambda
	$(GOGET) -u github.com/aws/aws-sdk-go
	$(GOGET) -u github.com/satori/go.uuid
	$(GOGET) -u github.com/gorilla/mux



