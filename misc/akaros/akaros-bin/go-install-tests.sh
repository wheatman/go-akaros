#!/bin/ash

AKAROS_MNT=/mnt
AKAROS_TMP=/tmp
CPIO_FILE=$AKAROS_MNT/gotests.cpio

# Make sure we are at /
cd /

# Make the /tmp directory in case it doesnt exist yet
mkdir -p ${AKAROS_TMP#/}

# Extract the $CPIO_FILE into /
cpio -d -i < $CPIO_FILE


