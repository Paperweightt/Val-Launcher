$env:GOOS = "windows"
$env:GOARCH = "amd64"

go build -o ../bin/launch.exe ../cmd/launch
go build -o ../bin/clean.exe ../cmd/reset_defaults

# TODO: replace above lines with these
# go build -ldflags "-H windowsgui" -o ../bin/launch.exe ../cmd/launch
# go build -ldflags "-H windowsgui" -o ../bin/clean.exe ../cmd/reset_defaults


# ISCC .\installer.iss
# go build -o bin/configur.exe ./cmd/configure
