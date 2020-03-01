.PHONY: coepi

GOBIN = $(shell pwd)/bin
GO ?= latest

cen:
		go build -o bin/cen
		@echo "Done building cen.  Run \"$(GOBIN)/cen\" to launch cen."
