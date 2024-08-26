build:
	@GOOS=linux GOARCH=amd64 go build -o bootstrap ./cmd/main.go
	@zip -r bootstrap.zip bootstrap

build-windows:
	@GOOS=windows GOARCH=amd64 go build -o bootstrap.exe ./cmd/main.go
	@zip -r bootstrap.zip bootstrap.exe

clean:
	@rm -f bootstrap bootstrap.zip bootstrap.exe
