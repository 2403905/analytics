run: build download-data start

build: # Build a binary
	@echo "  >  Building binary..."
	go build -o ./analytics ./cli/main.go

download-data: # Download the required dependencies
	@echo "  >  Downloading required dependencies..."
	@wget -q https://github.com/adjust/analytics-software-engineer-assignment/raw/master/data.tar.gz -O data.tar.gz
	@echo "  >  Uncompressing dependency files..."
	@tar -zxf data.tar.gz

start: # Start an application
	@echo "  >  Starting binary..."
	@./analytics

test:
	go test -bench=. ./... -benchmem
