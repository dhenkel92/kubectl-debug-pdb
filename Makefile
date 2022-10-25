ARGS=

.PHONY: run
run:
	go run ./cmd/kubectl-pdb/ cover ${ARGS}
