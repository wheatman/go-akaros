#!/usr/bin/env bash
eval $(go env)

HOST_MNT=$AKAROS_ROOT/mnt
HOST_GO_DIR=$GOROOT
HOST_BIN_DIR=$GOROOT/misc/akaros/akaros-bin
GO_CPIO_FILE=$HOST_MNT/go.cpio
BIN_CPIO_FILE=$HOST_MNT/bin.cpio

archive()
{
	local PWD=$1
	local FILES=$2
	local CPIO_FILE=$3

	echo "Packaging up ${PWD/$GOROOT/$\$GOROOT} into a cpio archive in ${HOST_MNT/$GOROOT/\$GOROOT}"
	rm -rf $CPIO_FILE
	echo "$FILES" | cpio --no-absolute-filenames -H newc -o > $CPIO_FILE
}

go_archive()
{
	# Archive all the files in the go tree, excluding the .hg, .git, and
	# $HOST_MNT folders
	cd $HOST_GO_DIR > /dev/null
	local FILES=$(find . \( -path ./.hg -o \
	                  -path ./.git -o \
	                  -path ./misc/akaros/mnt \) \
	                -prune -o -print)
	archive "$HOST_GO_DIR" "$FILES" "$GO_CPIO_FILE"
	cd - > /dev/null
}

bin_archive()
{
	cd $HOST_BIN_DIR > /dev/null
	archive "$HOST_BIN_DIR" "$(find .)" "$BIN_CPIO_FILE"
	cd - > /dev/null
}

# Run the appropriate function
targets="go bin"
for t in $targets; do
    if [ "$t" = "$1" ]; then
		${t}_archive
		exit
    fi
done
echo The first argument to $0 must be one of \{${targets// /, }\}

