define make_build
	rm -f builds/tmp/*
	GOOS=$(1) GOARCH=$(2) go build -o builds/tmp/
	cp -f README.md builds/tmp/
	cd builds/tmp && zip --recurse-paths --move ../$(basename $3)-$(1)-$(2).zip . && cd -
endef

# Batch build
build: deps build-linux build-macosx build-windows

# Dependencies
deps:
	go mod tidy && go mod vendor

# Linux
build-linux:
	$(call make_build,linux,amd64,fsync)

# MacOSX
build-macosx:
	$(call make_build,darwin,amd64,fsync)
	$(call make_build,darwin,arm64,fsync)

# Windows
build-windows:
	$(call make_build,windows,amd64,fsync.exe)
