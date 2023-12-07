ARGS=

.PHONY: build
build: kubectl-debug-pdb

kubectl-debug-pdb:
	go build -o bin/kubectl-debug_pdb ./cmd/kubectl-debug-pdb/kubectl-debug-pdb.go

.PHONY: test
test:
	go test -v ./...

.PHONY: run
run:
	go run ./cmd/kubectl-debug-pdb/ ${ARGS}
