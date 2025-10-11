# This creates Windows executables
$env:GOOS="windows"; $env:GOARCH="amd64"; go build -o ./bin/villain_couch_amd64.exe ./agent/src/
$env:GOOS="windows"; $env:GOARCH="arm64"; go build -o ./bin/villain_couch_arm64.exe ./agent/src/
$env:GOOS="windows"; $env:GOARCH="386"; go build -o ./bin/villain_couch_32bit.exe ./agent/src/
# $env:GOOS="windows"; $env:GOARCH="arm64"; go build -ldflags -H=windowsgui -o ./bin/villain_couch_arm64.exe ./agent/src/
# $env:GOOS="windows"; $env:GOARCH="386";   go build -ldflags -H=windowsgui -o ./bin/villain_couch_32bit.exe ./agent/src/
