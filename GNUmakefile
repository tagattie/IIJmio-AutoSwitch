ifeq ($(OS), Windows_NT)
    NAME	:= mioswitch.exe
else
    NAME	:= mioswitch
endif
VERSION		:= 0.3.0
REVISION	:= $(shell git rev-parse --short HEAD)

SRCS		:= $(shell find . -depth -maxdepth 1 -type f -name '*.go')
LDFLAGS		:= -ldflags="-X \"main.version=$(VERSION)\" -X \"main.revision=$(REVISION)\" -extldflags \"-static\""
DOCS		:= LICENSE README.md

all: bin/$(NAME)

bin/$(NAME): $(SRCS)
	go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o bin/$(NAME)

.PHONY: install
install:
	go install $(LDFLAGS)

.PHONY: clean
clean:
	rm -rf bin *.tar.gz *.zip

.PHONY: release
release:
	for os in darwin freebsd linux windows; do \
		GOOS=$$os GOARCH=amd64 CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o $(NAME); \
		if [ $$os = "windows" ]; then \
			zip -r mioswitch-$(VERSION)-$(REVISION).$$os-amd64.zip $(NAME) $(DOCS); \
		else \
			tar -czf mioswitch-$(VERSION)-$(REVISION).$$os-amd64.tar.gz $(NAME) $(DOCS); \
		fi; \
		rm -f $(NAME); \
	done
