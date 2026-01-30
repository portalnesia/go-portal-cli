APP=portal-cli

gobuild_win:
	go build -ldflags "-s -w" -o build/${APP}.exe main.go

compress_win:
	upx --brute build/${APP}.exe

build_win:
	$(MAKE) gobuild_win
	$(MAKE) compress_win
	@echo "Build completed"

.PHONY: build_win compress_win build_win