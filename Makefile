USERSPACE=COOLizh
SERVICE=sqs-sns-examples
GIT_HOST=github.com

APP=./build/service_${SERVICE}
PROJECT?=${GIT_HOST}/${USERSPACE}/${SERVICE}
DATE := $(shell date +'%Y.%m.%d %H:%M:%S')

GOOS?=linux
GO111MODULE?=on

clean:
	@[ -f ${APP} ] && rm -f ${APP} || true

build: clean
	CGO_ENABLED=0 GOOS=${GOOS} go build -a -installsuffix cgo \
		-ldflags '-s -w -X "${PROJECT}/internal/version.DATE=${DATE}" -X ${PROJECT}/internal/version.COMMIT=${COMMIT} -X ${PROJECT}/internal/version.TAG=${TAG}' \
		-o ${APP} ./cmd/main.go
