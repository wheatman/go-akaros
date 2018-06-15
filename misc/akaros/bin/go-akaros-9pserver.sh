#!/usr/bin/env bash
eval $(go env)

: ${CLEAR_MOUNT:=false}
: ${UFS_PORT:="1025"}

HOST_MNT=$AKAROS_ROOT/mnt
ARCHIVE_SCRIPT=$GOROOT/misc/akaros/bin/go-akaros-archive.sh

if [ "$GOPATH" = "" ]; then
	echo "You must have \$GOPATH set in order to run this script!"
	exit 1
fi

# Clear out the $HOST_MNT directory
if [ $CLEAR_MOUNT = true ]; then
	echo "Clearing out ${HOST_MNT/$GOROOT/\$GOROOT}"
	rm -rf $HOST_MNT
fi
mkdir -p $HOST_MNT

# Leverage the archive script to put an archive of the go tree at $HOST_MNT
$ARCHIVE_SCRIPT go 2>/dev/null

# Kill any old instances of the ufs server and start a new one
echo "Starting the 9p server port=$UFS_PORT root=${HOST_MNT/$GOROOT/\$GOROOT}"
ps aux | grep "ufs -akaros=true" | head -1 | awk '{print $2}' | xargs kill >/dev/null 2>&1
ufs -akaros=true -addr=:$UFS_PORT -root=$HOST_MNT &

