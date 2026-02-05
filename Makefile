BINARY_NAME=gosh
CMD_DIR=./cmd
DIST_DIR=./dist

BINS := $(notdir $(wildcard $(CMD_DIR)/*))

.PHONY: all build clean help $(BINS)

all: $(BINS)

$(BINS):
	@echo "Building $@..."
	@go build -o $(DIST_DIR)/$@ $(CMD_DIR)/$@/main.go

clean:
	@rm -rf $(DIST_DIR)
	@echo "Cleaned $(DIST_DIR)"

help:
	@echo "Usage: make [target]"
	@sed -n 's/^##//p' $(Makefile) | column -t -s ':' | sed -e 's/^/ /'
