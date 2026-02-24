# 项目名称（用于 release 文件名）
APP_NAME=tmdb-list-exporter

# 二进制名称（真正可执行文件名）
BIN_NAME=tmdb

# 输出目录
DIST_DIR=./bin

# 版本号（可通过 make VERSION=1.0.0 传入）
VERSION?=dev

# 构建参数
LDFLAGS=-ldflags="-s -w -X main.version=$(VERSION)"

# 默认目标
all: clean build-all

# 创建目录
prepare:
	mkdir -p $(DIST_DIR)

########################################
# 构建函数模板
########################################

define build_platform
	@echo ">> Building $(1)-$(2)"
	mkdir -p $(DIST_DIR)/$(1)-$(2)
	GOOS=$(1) GOARCH=$(2) go build $(LDFLAGS) -o $(DIST_DIR)/$(1)-$(2)/$(BIN_NAME)$(3)
	cd $(DIST_DIR)/$(1)-$(2) && \
	tar -czf ../$(APP_NAME)-v$(VERSION)-$(1)-$(2).tar.gz $(BIN_NAME)$(3)
	rm -rf $(DIST_DIR)/$(1)-$(2)
endef

########################################
# 各平台构建
########################################

build-darwin-amd64:
	$(call build_platform,darwin,amd64,)

build-darwin-arm64:
	$(call build_platform,darwin,arm64,)

build-linux-amd64:
	$(call build_platform,linux,amd64,)

build-linux-arm64:
	$(call build_platform,linux,arm64,)

build-windows-amd64:
	@echo ">> Building windows-amd64"
	mkdir -p $(DIST_DIR)/windows-amd64
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(DIST_DIR)/windows-amd64/$(BIN_NAME).exe
	cd $(DIST_DIR)/windows-amd64 && \
	zip ../$(APP_NAME)-v$(VERSION)-windows-amd64.zip $(BIN_NAME).exe
	rm -rf $(DIST_DIR)/windows-amd64

########################################
# 构建全部平台
########################################

build-all: prepare \
	build-darwin-amd64 \
	build-darwin-arm64 \
	build-linux-amd64 \
	build-linux-arm64 \
	build-windows-amd64

########################################
# 清理
########################################

clean:
	rm -rf $(DIST_DIR)
