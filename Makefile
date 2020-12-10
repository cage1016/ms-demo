add:
	curl -X POST localhost:8180/sum -d '{"a": 1, "b":1}'

gadd:
	grpcurl -plaintext -proto ./pb/add/add.proto -d '{"a": 1, "b":1}' localhost:8181 pb.Add.Sum

tic:
	curl -X POST localhost:9190/tic

tac:
	curl localhost:9190/tac	

a:
	curl -X POST localhost:80/api/v1/add/sum -d '{"a": 1, "b":1}'

b:
	grpcurl -plaintext -proto ./pb/add/add.proto -d '{"a": 1, "b":1}' localhost:443 pb.Add.Sum

c:
	curl -X POST localhost:80/api/v1/tictac/tic

d:
	curl localhost:80/api/v1/tictac/tac

# Regenerates OPA data from rego files
HAVE_GO_BINDATA := $(shell command -v mockgen 2> /dev/null)
generate:
ifndef HAVE_GO_BINDATA
	@echo "requires 'mockgen' (GO111MODULE=on go get github.com/golang/mock/mockgen@v1.4.4)"
	@exit 1 # fail
else
	go generate ./...
endif
