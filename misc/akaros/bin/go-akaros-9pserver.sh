#!/usr/bin/env bash
eval $(go env)

: ${UFS_PORT:="1025"}

HOST_MNT=$GOROOT/misc/akaros/mnt
ARCHIVE_SCRIPT=$GOROOT/misc/akaros/bin/go-akaros-archive.sh

if [ "$GOPATH" = "" ]; then
	echo "You must have \$GOPATH set in order to run this script!"
	exit 1
fi

# Get the latest 9p server which supports akaros
echo "Downloading and installing the latest supported 9p server"
GOOS=$GOHOSTOS
GOARCH=$GOHOSTARCH
CGO_ENABLED=0
go get -a github.com/rminnich/go9p
go get -a github.com/rminnich/go9p/ufs
go install -a github.com/rminnich/go9p/ufs

# Clear out the $HOST_MNT directory
echo "Clearing out ${HOST_MNT/$GOROOT/\$GOROOT}"
rm -rf $HOST_MNT
mkdir -p $HOST_MNT

# Leverage the archive script to put an archive of the go tree at $HOST_MNT
$ARCHIVE_SCRIPT go 2>/dev/null

# Kill any old instances of the ufs server and start a new one
echo "Starting the 9p server port=$UFS_PORT root=${HOST_MNT/$GOROOT/\$GOROOT}"
ps aux | grep "ufs -addr=:$UFS_PORT" | head -1 | awk '{print $2}' | xargs kill >/dev/null 2>&1
nohup ufs -addr=:$UFS_PORT -root=$HOST_MNT >/dev/null 2>&1 &

echo "Done"

