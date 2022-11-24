PORT = 8000 8001 8002 8003

hello: $(addprefix hello-, $(NAMES))

.PHONY:
build-bc-1:
	go run blockchain/main.go --port 8001
.PHONY: build-bc-2
build-bc-2:
	go run blockchain/main.go --port 8002
.PHONY: build-bc-3
build-bc-3:
	go run blockchain/main.go --port 8003
