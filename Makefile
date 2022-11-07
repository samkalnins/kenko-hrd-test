run:
	go run *.go

build:
	GOOS=darwin GOARCH=arm64 go build -o bin/kenko_bt_grabber.macm1