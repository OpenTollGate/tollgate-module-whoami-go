include $(TOPDIR)/rules.mk

PKG_NAME:=tollgate-module-whoami-go
PKG_VERSION:=$(shell git rev-list --count HEAD 2>/dev/null || echo "0.0.1").$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
PKG_RELEASE:=1

# Place conditional checks EARLY - before variables that depend on them
ifneq ($(TOPDIR),)
	# Feed-specific settings (auto-clone from git)
	PKG_SOURCE_PROTO:=git
	PKG_SOURCE_URL:=https://github.com/OpenTollGate/tollgate-module-whoami-go.git
	PKG_SOURCE_VERSION:=main
	PKG_MIRROR_HASH:=skip
else
	# SDK build context (local files)
	PKG_BUILD_DIR:=$(CURDIR)
endif

PKG_MAINTAINER:=Your Name <your@email.com>
PKG_LICENSE:=CC0-1.0
PKG_LICENSE_FILES:=LICENSE

PKG_BUILD_DEPENDS:=golang/host
PKG_BUILD_PARALLEL:=1
PKG_USE_MIPS16:=0

GO_PKG:=github.com/OpenTollGate/tollgate-module-whoami-go

include $(INCLUDE_DIR)/package.mk
# include $(INCLUDE_DIR)/golang-package.mk
$(eval $(call GoPackage))

define Package/$(PKG_NAME)
	SECTION:=net
	CATEGORY:=Network
	TITLE:=TollGate Whoami Module
	DEPENDS:=$(GO_ARCH_DEPENDS)
endef

define Package/$(PKG_NAME)/description
	TollGate WhoamI Module for OpenWrt
endef

define Build/Prepare
	# First call the default prepare to handle git source fetching
	$(call Build/Prepare/Default)

	# Debug: check what's in the build directory after default prepare
	@echo "Contents of PKG_BUILD_DIR after default prepare:"
	@ls -la $(PKG_BUILD_DIR) || echo "Directory empty or doesn't exist"

	# Clean up Go module cache if needed
	cd $(PKG_BUILD_DIR) && \
	rm -f go.mod go.sum && \
	go mod init tollgate-module-whoami-go && \
	go mod edit -replace github.com/OpenTollgate/whoami=./ && \
	go mod tidy && \
	go get github.com/fiatjaf/khatru && \
	go get github.com/nbd-wtf/go-nostr
endef

define Build/Configure
endef

define Build/Compile
	mkdir -p "$(PKG_BUILD_DIR)/.go_work/build/src/github.com/OpenTollGate"
	ln -sf "$(PKG_BUILD_DIR)" "$(PKG_BUILD_DIR)/.go_work/build/src/github.com/OpenTollGate/tollgate-module-whoami-go"

	cd $(PKG_BUILD_DIR) && \
	GO111MODULE=on \
	CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=arm64 \
	GOFLAGS="-mod=readonly -modcacherw" \
	GOWORK=off \
	go mod tidy && \
	go mod vendor && \
	GODEBUG=netdns=go+1 \
	go build -v -x -mod=vendor -trimpath \
		-tags netgo,osusergo \
		-ldflags "-s -w -extldflags '-static' \
		-X main.Version=$(PKG_VERSION) \
		-X main.CommitHash=$(PKG_RELEASE) \
		-X main.BuildTime=$(shell date -u '+%Y-%m-%d_%I:%M:%S%p')" \
		-o $(PKG_BUILD_DIR)/tollgate-whoami main.go
endef

define Package/$(PKG_NAME)/install
	$(INSTALL_DIR) $(1)/usr/bin
	$(INSTALL_BIN) $(PKG_BUILD_DIR)/tollgate-whoami $(1)/usr/bin/tollgate-whoami
endef

$(eval $(call BuildPackage,$(PKG_NAME)))

# Print IPK path after successful compilation
PKG_FINISH:=$(shell echo "Successfully built: $(IPK_FILE)" >&2)
