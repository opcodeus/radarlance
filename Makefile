PROJECT_NAME := radarlance
BUILD_DIR := build
GOFLAGS := -ldflags "-s -w" -trimpath -buildvcs=false
GO_BUILD := go build $(GOFLAGS)
.PHONY: all clean linux windows darwin tidy

all: tidy linux windows darwin

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

tidy:
	go mod tidy

linux: linux-amd64 linux-386

linux-amd64: $(BUILD_DIR)/$(PROJECT_NAME)-linux-amd64

$(BUILD_DIR)/$(PROJECT_NAME)-linux-amd64: tidy | $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GO_BUILD) -o $(BUILD_DIR)/$(PROJECT_NAME)-linux-amd64

linux-386: $(BUILD_DIR)/$(PROJECT_NAME)-linux-386

$(BUILD_DIR)/$(PROJECT_NAME)-linux-386: tidy | $(BUILD_DIR)
	GOOS=linux GOARCH=386 $(GO_BUILD) -o $(BUILD_DIR)/$(PROJECT_NAME)-linux-386

windows: windows-amd64 windows-386

windows-amd64: $(BUILD_DIR)/$(PROJECT_NAME)-windows-amd64.exe

$(BUILD_DIR)/$(PROJECT_NAME)-windows-amd64.exe: tidy | $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GO_BUILD) -o $(BUILD_DIR)/$(PROJECT_NAME)-windows-amd64.exe

windows-386: $(BUILD_DIR)/$(PROJECT_NAME)-windows-386.exe

$(BUILD_DIR)/$(PROJECT_NAME)-windows-386.exe: tidy | $(BUILD_DIR)
	GOOS=windows GOARCH=386 $(GO_BUILD) -o $(BUILD_DIR)/$(PROJECT_NAME)-windows-386.exe

darwin: darwin-amd64 darwin-arm64

darwin-amd64: $(BUILD_DIR)/$(PROJECT_NAME)-darwin-amd64

$(BUILD_DIR)/$(PROJECT_NAME)-darwin-amd64: tidy | $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GO_BUILD) -o $(BUILD_DIR)/$(PROJECT_NAME)-darwin-amd64

darwin-arm64: $(BUILD_DIR)/$(PROJECT_NAME)-darwin-arm64

$(BUILD_DIR)/$(PROJECT_NAME)-darwin-arm64: tidy | $(BUILD_DIR)
	GOOS=darwin GOARCH=arm64 $(GO_BUILD) -o $(BUILD_DIR)/$(PROJECT_NAME)-darwin-arm64

clean:
	rm -rf $(BUILD_DIR)
