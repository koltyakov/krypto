# `make version=0.1.0 tag`
version := 0.0.0-SNAPSHOT
author  := Andrew Koltyakov
app     := Krypto
id      := com.koltyakov.krypto

install:
	go get -u ./... && go mod tidy
	which appify || go get github.com/machinebox/appify
	which create-dmg || npm i -g create-dmg

format:
	gofmt -s -w .

generate:
	cd icon/ && ./gen.sh
	make format

build-darwin:
	GOOS=darwin GOARCH=amd64 go build -v -ldflags "-s -w -X main.version=$(version)" -o bin/darwin/krypto ./

build:
	go build -v -o bin/krypto ./

clean:
	rm -rf bin/ dist/

bundle-darwin: build-darwin
	# Package solution to .app folder
	cd bin/darwin/ && \
		appify \
			-author "$(author)" \
			-id $(id) \
			-version $(version) \
			-name "$(app)" \
			-icon ../../assets/icon.png \
			./krypto
	/usr/libexec/PlistBuddy -c 'Add :LSUIElement bool true' 'bin/darwin/$(app).app/Contents/Info.plist'
	rm 'bin/darwin/$(app).app/Contents/README'
	# Package solution to .dmg image
	cd bin/darwin/ && \
		npx create-dmg --dmg-title='$(app)' '$(app).app' ./ \
			|| true # ignore Error 2
	# Rename .dmg appropriately
	mv 'bin/darwin/$(app) $(version).dmg' bin/darwin/krypto_$(version)_darwin_x86_64.dmg
	# Remove temp files
	rm -rf 'bin/darwin/$(app).app'

tag:
	git tag -fa v$(version) -m "Version $(version)"
	git push origin --tags

release-snapshot:
	curl -sfL https://install.goreleaser.com/github.com/goreleaser/goreleaser.sh | sh
	./bin/goreleaser --rm-dist --skip-publish --snapshot
	cd dist && ls *.dmg | xargs shasum -a256 >> $$(ls *_checksums.txt)

release:
	curl -sfL https://install.goreleaser.com/github.com/goreleaser/goreleaser.sh | sh
	./bin/goreleaser --rm-dist
	cd dist && ls *.dmg | xargs shasum -a256 >> $$(ls *_checksums.txt)

lint-cyclo:
	which gocyclo || go get github.com/fzipp/gocyclo/cmd/gocyclo
	gocyclo ./

start: run # alias for run
run:
	pkill krypto || true
	nohup go run ./ >/dev/null 2>&1 &