GO ?= go
TARGET := check_network
OS := freebsd netbsd linux
ARCH := 386 amd64
VERSION := "1.0"
HOSTOS := `go env | grep GOHOSTOS | cut -d '"' -f 2 | head -1`

all: build

build:
	@for arch in $(ARCH); do \
		echo "===> building: $(TARGET)-$(HOSTOS)-$$arch-$(VERSION)"; \
		if [ $$arch == "386" ] ; then \
			GOOS=$$os GOARCH=$$arch go build -o $(TARGET)-$(HOSTOS)-"i386"-"$(VERSION)" $^ ;\
		else \
			GOOS=$$os GOARCH=$$arch go build -o $(TARGET)-$(HOSTOS)-$$arch-"$(VERSION)" $^ ;\
		fi \
	done \

release:
	@for os in $(OS); do \
		for arch in $(ARCH); do \
			echo "===> building: $(TARGET)-$$os-$$arch-$(VERSION)"; \
			if [ $$arch == "386" ] ; then \
				GOOS=$$os GOARCH=$$arch go build -o $(TARGET)-$$os-"i386"-"$(VERSION)" $^ ;\
			else \
				GOOS=$$os GOARCH=$$arch go build -o $(TARGET)-$$os-$$arch-"$(VERSION)" $^ ;\
			fi \
		done \
	done \

clean:
	@$(GO) clean
	@for os in $(OS); do \
		for arch in $(ARCH); do \
		echo "===> Removing: $(TARGET)-$$os-$$arch-$(VERSION)"; \
		rm -f $(TARGET)-$$os-$$arch-"$(VERSION)" $^ ;\
		done \
	done \


.PHONY: all release build clean
