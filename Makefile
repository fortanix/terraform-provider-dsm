TEST?=$$(go list ./... | grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
HOSTNAME=fortanix.com
NAMESPACE=fortanix
NAME=dsm
BINARY=terraform-provider-${NAME}
VERSION=0.5.16
OS_ARCH=linux_amd64

default: install

fmt: 
	gofmt -w $(GOFMT_FILES)

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BINARY}

release:
	mkdir -p ./bin/${VERSION}
	mkdir -p ./bin/${VERSION}/linux_386
	mkdir -p ./bin/${VERSION}/linux_amd64
	mkdir -p ./bin/${VERSION}/linux_arm
	mkdir -p ./bin/${VERSION}/windows_386
	mkdir -p ./bin/${VERSION}/windows_amd64
	mkdir -p ./bin/${VERSION}/darwin_amd64
	mkdir -p ./bin/${VERSION}/darwin_arm64
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o ./bin/${VERSION}/linux_386/${BINARY}_v${VERSION}
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/${VERSION}/linux_amd64/${BINARY}_v${VERSION}
	CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o ./bin/${VERSION}/linux_arm/${BINARY}_v${VERSION}
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o ./bin/${VERSION}/windows_386/${BINARY}_v${VERSION}
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./bin/${VERSION}/windows_amd64/${BINARY}_v${VERSION}
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./bin/${VERSION}/darwin_amd64/${BINARY}_v${VERSION}
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o ./bin/${VERSION}/darwin_arm64/${BINARY}_v${VERSION}
	zip -j ./${BINARY}_${VERSION}_linux_386.zip ./bin/${VERSION}/linux_386/${BINARY}_v${VERSION}
	zip -j ./${BINARY}_${VERSION}_linux_amd64.zip ./bin/${VERSION}/linux_amd64/${BINARY}_v${VERSION}
	zip -j ./${BINARY}_${VERSION}_linux_arm.zip ./bin/${VERSION}/linux_arm/${BINARY}_v${VERSION}
	zip -j ./${BINARY}_${VERSION}_windows_386.zip ./bin/${VERSION}/windows_386/${BINARY}_v${VERSION}
	zip -j ./${BINARY}_${VERSION}_windows_amd64.zip ./bin/${VERSION}/windows_amd64/${BINARY}_v${VERSION}
	zip -j ./${BINARY}_${VERSION}_darwin_amd64.zip ./bin/${VERSION}/darwin_amd64/${BINARY}_v${VERSION}
	zip -j ./${BINARY}_${VERSION}_darwin_arm64.zip ./bin/${VERSION}/darwin_arm64/${BINARY}_v${VERSION}
	shasum -a 256 *.zip > ${BINARY}_${VERSION}_SHA256SUMS


install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

test: 
	go test -i $(TEST) || exit 1                                                   
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4                    

testacc: 
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m   
