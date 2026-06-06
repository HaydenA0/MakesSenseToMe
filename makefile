.PHONY: cli srv run install build

cli:
	go run ./src/cmd/cli/main.go

srv:
	go run ./src/cmd/server/main.go

install:
	cd frontend && npm install

build:
	cd frontend && npm run build

run:
	@if [ ! -d frontend/node_modules ]; then \
		echo ">> installing frontend deps"; \
		$(MAKE) install; \
	fi
	@trap 'kill 0' INT TERM EXIT; \
	go run ./src/cmd/server/main.go & \
	cd frontend && npm run dev & \
	wait
