.PHONY: build install clean release tag

build:
	go build -o bin/supershell ./cmd/supershell

install: build
	install -m 0755 bin/supershell /usr/local/bin/supershell

clean:
	rm -rf bin dist

release:
	bash scripts/release_build.sh $(VERSION)

tag:
	bash scripts/release_tag.sh $(VERSION)


