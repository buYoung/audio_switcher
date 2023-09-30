APP_NAME=audio_switcher
APP_VERSION=1.0.0

MAIN_FILE=main.go

GO=go

BUILD_DIR=build

# Production 빌드를 위한 플래그 (최적화와 디버그 정보 제거)
GO_BUILD_FLAGS=-ldflags="-s -w"

build-prod:
	$(GO) build $(GO_BUILD_FLAGS) -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE)