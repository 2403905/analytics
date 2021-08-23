run: build download-data start

build: # Build binary
	@echo "  >  Building binary..."
	go build -o ./analytics ./cli/main.go

download-data: # Download required dependencies
	@echo "  >  Downloading required dependencies..."
	@wget -q https://github.com/adjust/analytics-software-engineer-assignment/raw/master/data.tar.gz -O data.tar.gz
	@echo "  >  Uncompressing dependency files..."
	@tar -zxf data.tar.gz

start: # Start application
	@echo "  >  Start binary..."
	@./analytics
