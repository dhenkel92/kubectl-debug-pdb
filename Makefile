ARGS=

.PHONY: build
build: kubectl-pdb

kubectl-pdb:
	go build -o bin/kubectl-pdb ./cmd/kubectl-pdb/kubectl-pdb.go

.PHONY: test
test:
	go test -v ./...

.PHONY: run
run:
	go run ./cmd/kubectl-pdb/ cover ${ARGS}
