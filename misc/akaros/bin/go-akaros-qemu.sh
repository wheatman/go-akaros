#!/usr/bin/env bash
eval $(go env)

: ${QEMU:="qemu-system-x86_64"}
: ${CPU_TYPE:="host"}
: ${NUM_CPUS:="8"}
: ${MEMORY:="4096"}
: ${NETWORK_CARD:="e1000"}
: ${HOST_PORT:="5555"}
: ${AKAROS_PORT:="5555"}
: ${QEMU_KVM:="-enable-kvm"}
: ${QEMU_MONITOR_TTY:=""}

AKAROS_BIN=$AKAROS_ROOT/kern/kfs/bin
GO_SCRIPTS_DIR=$GOROOT/misc/akaros/akaros-bin

#QEMU_NETWORK="-net nic,model=$NETWORK_CARD -net tap,ifname=tap0,script=no"
QEMU_NETWORK="-net nic,model=$NETWORK_CARD -net user,hostfwd=tcp::$HOST_PORT-:$AKAROS_PORT"
if [ "$QEMU_MONITOR_TTY" != "" ]; then
	QEMU_MONITOR="-monitor $QEMU_MONITOR_TTY"
fi

if [ "$AKAROS_ROOT" = "" ]; then
	echo "You must have \$AKAROS_ROOT set in order to run this script!"
	exit 1
fi

if [ "$(which $QEMU)" = "" ]; then
	echo "You must have $QEMU installed in order to run this script!"
	exit 1
fi

if [ "$QEMU_KVM" == "-enable-kvm" ]; then
	groups $USER | grep &>/dev/null '\bkvm\b'
	if [ "$?" != "0" ]; then
		echo "You are not part of the kvm group!"
		echo "    This may cause problems with running qemu with kvm enabled."
		echo "    To disable kvm, rerun this script with QEMU_KVM=\"-no-kvm\""
	fi
fi

# Copy the go scripts into $AKAROS_BIN
echo "Copying scripts from ${GO_SCRIPTS_DIR/$GOROOT/\$GOROOT} into ${AKAROS_BIN/$AKAROS_ROOT/\$AKAROS_ROOT}"
cp $GO_SCRIPTS_DIR/* $AKAROS_BIN

# Rebuilding akaros
echo "Rebuilding akaros"
cd $AKAROS_ROOT
make
make install-libs
make tests
make fill-kfs

# Launching qemu
echo "Launching qemu"
$QEMU -s $QEMU_KVM $QEMU_NETWORK $QEMU_MONITOR -cpu $CPU_TYPE -smp $NUM_CPUS \
      -m $MEMORY -kernel $AKAROS_ROOT/obj/kern/akaros-kernel -nographic

