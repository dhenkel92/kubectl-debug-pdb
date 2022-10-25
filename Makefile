ARGS=

.PHONY: run
run:
	go run ./cmd/ cover ${ARGS}
