ifeq ($(OS), Windows_NT)
    NAME	:= mioswitch.exe
else
    NAME	:= mioswitch
endif
SRCS	:= $(shell find . -depth -maxdepth 1 -type f -name '*.go')
LDFLAGS	:= -ldflags="-extldflags \"-static\""

all: bin/$(NAME)

bin/$(NAME): $(SRCS)
	go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o bin/$(NAME)

.PHONY: install
install:
	go install $(LDFLAGS)

.PHONY: clean
clean:
	rm -rf bin
