NAME      := tcping
BUILD_DIR := build                     # root directory for build outputs

TARGET    = $(GOOS)-$(GOARCH)$(GOARM)  # target platform identifier
BIN_DIR   = $(BUILD_DIR)/$(TARGET)     # Platform-specific binary directory

VERSION   ?= dev

ifeq ($(GOOS),windows)
  EXT=.exe
  PACK_CMD=zip -9 -r $(NAME)-$(TARGET)-$(VERSION).zip $(TARGET)
else
  EXT=
  PACK_CMD=tar czpvf $(NAME)-$(TARGET)-$(VERSION).tar.gz $(TARGET)
endif

define check_env
	@ if [ "$(GOOS)" = "" ]; then echo " <- Env variable GOOS not set"; exit 1; fi
	@ if [ "$(GOARCH)" = "" ]; then echo " <- Env variable GOARCH not set"; exit 1; fi
endef

.PHONY: build build-dev cross-build release test clean

build: clean test
	@CGO_ENABLED=0 go build -ldflags "-s -w -X main.version=$(VERSION)" -o $(BIN_DIR)/$(NAME)$(EXT)
	@echo "binary file: $(BIN_DIR)/$(NAME)$(EXT)"

build-dev: clean test
	@go build -o $(BUILD_DIR)/$(NAME)$(EXT)
	@echo "binary file: $(BIN_DIR)/$(NAME)$(EXT)"

cross-build: clean test
	@$(call check_env)
	@CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "-s -w -X main.version=$(VERSION)" -o $(BIN_DIR)/$(NAME)$(EXT)
	@echo "binary file: $(BIN_DIR)/$(NAME)$(EXT)"

release: clean test
	@$(call check_env)
	@mkdir -p $(BUILD_DIR)
	@cp LICENSE $(BUILD_DIR)/
	@cp README.md $(BUILD_DIR)/
	@cp Examples.md $(BUILD_DIR)/
	@CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "-s -w -X main.version=$(VERSION)" -o $(BIN_DIR)/$(NAME)$(EXT)
	@cd $(BUILD_DIR) ; $(PACK_CMD)

test:
	@go test -race -v -bench=. ./...

clean:
	@go clean
	@rm -rf $(BUILD_DIR)
	@echo "done"
