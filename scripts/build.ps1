$env:GOOS = "windows"
$env:GOARCH = "amd64"

# dev build
# go build -o ../bin/launch.exe ../cmd/launch
# go build -o ../bin/clean.exe ../cmd/reset_defaults
go build -o ../bin/configure.exe ../cmd/configure/

# release build
go build -ldflags "-H windowsgui" -o ../bin/launch.exe ../cmd/launch
go build -ldflags "-H windowsgui" -o ../bin/clean.exe ../cmd/reset_defaults
# go build -ldflags "-H windowsgui" -o ../bin/configure.exe ../cmd/configure/

# go build -o ../bin/clean.exe ../cmd/reset_defaults
# go build -o ../bin/clean.exe ../cmd/reset_defaults

# ISCC .\installer.iss
# go build -o bin/configur.exe ./cmd/configure
