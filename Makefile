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