TEST?=$$(go list ./... | grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
HOSTNAME=fortanix.com
NAMESPACE=fyoo
NAME=dsm
BINARY=terraform-provider-${NAME}
VERSION=0.4.0
OS_ARCH=linux_amd64

default: install

fmt: 
	gofmt -w $(GOFMT_FILES)

build:
	go build -o ${BINARY}

release:
	mkdir -p ./bin/${VERSION}
	mkdir -p ./bin/${VERSION}/linux_386
	mkdir -p ./bin/${VERSION}/linux_amd64
	mkdir -p ./bin/${VERSION}/linux_arm
	mkdir -p ./bin/${VERSION}/windows_386
	mkdir -p ./bin/${VERSION}/windows_amd64
	mkdir -p ./bin/${VERSION}/darwin_amd64
	mkdir -p ./bin/${VERSION}/darwin_arm64
	GOOS=linux GOARCH=386 go build -o ./bin/${VERSION}/linux_386/${BINARY}
	GOOS=linux GOARCH=amd64 go build -o ./bin/${VERSION}/linux_amd64/${BINARY}
	GOOS=linux GOARCH=arm go build -o ./bin/${VERSION}/linux_arm/${BINARY}
	GOOS=windows GOARCH=386 go build -o ./bin/${VERSION}/windows_386/${BINARY}
	GOOS=windows GOARCH=amd64 go build -o ./bin/${VERSION}/windows_amd64/${BINARY}
	GOOS=darwin GOARCH=amd64 go build -o ./bin/${VERSION}/darwin_amd64/${BINARY}
	GOOS=darwin GOARCH=arm64 go build -o ./bin/${VERSION}/darwin_arm64/${BINARY}

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

test: 
	go test -i $(TEST) || exit 1                                                   
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4                    

testacc: 
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m   