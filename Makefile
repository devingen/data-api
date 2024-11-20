# Check if VERSION is provided
ifndef VERSION
$(error VERSION is not defined. Use 'make release VERSION=<version>' to provide it.)
endif

# Build for multiple platforms with version
.PHONY: release
release:
	git tag -a v$(VERSION) -m ""
	git push origin v$(VERSION)
