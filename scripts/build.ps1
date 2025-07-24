# go build -o bin/launch.exe ./cmd/launch
# targets widows vs other prev cmd
GOOS=windows GOARCH=amd64 go build -o MyApp.exe ./cmd/launch

# go build -o bin/configur.exe ./cmd/configure
