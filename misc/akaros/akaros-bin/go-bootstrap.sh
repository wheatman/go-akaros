#!/bin/ash

AKAROS_MNT=/mnt
AKAROS_GOROOT=/go
CPIO_FILE=$AKAROS_MNT/go.cpio

# Make sure we are at /
cd /

# Run the root script to get our folder mounted properly
ash root $@

# Create $AKAROS_GOROOT
mkdir -p ${AKAROS_GOROOT#/}

# Copy over the archived go tree into /
cp $CPIO_FILE /

# Extract it into $AKAROS_GOROOT
cd $AKAROS_GOROOT
cpio -d -i < $CPIO_FILE

# Go back to / and start the listen server from there
cd /

# Start the listen daemon to wait for incoming rpc calls
listen1 tcp!*!23 /bin/ash &

