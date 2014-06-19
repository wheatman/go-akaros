eval $(go env)

export GOOS=akaros
export GOARCH=amd64
export CGO_ENABLED=1

export PATH=$GOROOT/misc/akaros:$PATH

