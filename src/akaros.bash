#!/usr/bin/env bash
# Copyright 2013 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# Error out if any subcommand or pipeline returns a non-zero status
set -e

#Grab a reference to the command so we don't always just refer to it as $0
COMMAND=$0

# We are only using this script to build for Akaros
export GOOS=akaros

# Always build with cgo support
export CGO_ENABLED=1

# Print the usage information for this script
usage()
{
cat << EOF
This script will build a cross compiler for Akaros on your local machine
usage: $0 options
OPTIONS:
   -a      This is the default option
           Build everything (same as '$COMMAND -b all')
   -b      Pass one of the following arguments to build the following:
           host   - the dist tool, go_bootstrap, runtime and other packages
                    for your host machine
           akaros - the runtime and other packages for akaros
           all    - build both the host and the akaros packages
   -v      Verbose mode
   -e      Show the list of environment variables that can be set to control
           this script
   -h      Show this message
EOF
}

# Print the environment variables that can be set to control this script
environment()
{
cat << EOF
Environment variables that control akaros.bash:

GOROOT_FINAL: The expected final Go root, baked into binaries.
The default is the location of the Go tree during the build.

GOHOSTARCH: The architecture for host tools (compilers and
binaries).  Binaries of this type must be executable on the current
system, so the only common reason to set this is to set
GOHOSTARCH=386 on an amd64 machine.

GOARCH: The target architecture for installed packages and tools.

GO_GCFLAGS: Additional 5g/6g/8g arguments to use when
building the packages and commands.

GO_LDFLAGS: Additional 5l/6l/8l arguments to use when
building the commands.

GO_CCFLAGS: Additional 5c/6c/8c arguments to use when
building.

GO_EXTLINK_ENABLED: Set to 1 to invoke the host linker when building
packages that use cgo.  Set to 0 to do all linking internally.  This
controls the default behavior of the linker's -linkmode option.  The
default value depends on the system.

CC: Command line to run to get at host C compiler.
Default is "gcc". Also supported: "clang".
EOF
}

# Check to make sure this script is being invoked from the proper location
if [ ! -f make.bash ]; then
	echo "$COMMAND must be run from \$GOROOT/src" 1>&2
    usage
	exit 1
fi

# Source a local.bash file in case the user wants to put his environment
# variable setup in there
if [ -f local.bash ]; then
    source local.bash
fi

# Get the options from the command line
BUILD="all"
while getopts “ab:veh” OPTION
do
     case $OPTION in
         a)
             ;;
         b)
             BUILD=$OPTARG
             ;;
         v)
             VERBOSE=1
             ;;
         e)
             environment
             exit 1
             ;;
         h)
             usage
             exit 1
             ;;
     esac
done

# Make sure that only valid build options have been passed
if [[ "$BUILD" != "all" ]] &&
   [[ "$BUILD" != "host" ]] &&
   [[ "$BUILD" != "akaros" ]]
then
	echo "Illegal build option for $COMMAND: $BUILD" 1>&2
    echo
    usage
	exit 1
fi

# Build the run helper differently depending on whether verbose was passed in
if [[ "$VERBOSE" = "1" ]]; then
  run_helper() { echo "$@"; "$@"; }
else
  run_helper() { "$@"; }
fi
if [[ "$VERBOSE" = "1" ]]; then
  run_helper_eval() { echo "$@"; eval "$@"; }
else
  run_helper_eval() { eval "$@"; }
fi


# Pepare for building host tools by saving off some environment variables
prepare_host_env()
{
  OLD_GOOS=$GOOS
  OLD_GOARCH=$GOARCH
  OLD_CC=$CC
  OLD_CXX=$CXX
  OLD_GO_LDFLAGS="$GO_LDFLAGS"
  export GOOS=$GOHOSTOS
  export GOARCH=$GOHOSTARCH
  export CC=$CC
  export CXX=$CXX
  export GO_LDFLAGS="$GO_LDFLAGS"
}

# Pepare for building target tools by saving off some environment variables
prepare_target_env()
{
  OLD_GOOS=$GOOS
  OLD_GOARCH=$GOARCH
  OLD_CC=$CC
  OLD_CXX=$CXX
  OLD_GO_LDFLAGS="$GO_LDFLAGS"
  export GOOS=$GOOS
  export GOARCH=$GOARCH
  export CC=$TARGETCC
  export CXX=$TARGETCXX
  export GO_LDFLAGS="-extld=$TARGETCC $GO_LDFLAGS"
}

# Restore the environment variables saved by prepare_{host,target}_env
restore_env()
{
  export GOOS=$OLD_GOOS
  export GOARCH=$OLD_GOARCH
  export CC=$OLD_CC
  export CXX=$OLD_CXX
  export GO_LDFLAGS="$OLD_GO_LDFLAGS"
  export CGO_LDFLAGS="$OLD_CGO_LDFLAGS"
}

# Clean old generated file that will cause problems in the build.
rm -f ./pkg/runtime/runtime_defs.go
  
# Build the dist tool, compilers and go bootstrap tool for the host machine
if [[ "$BUILD" = "all" ]] ||
   [[ "$BUILD" = "host" ]] 
then
  # Build the dist tool
  echo '# Building C bootstrap tool.'
  echo cmd/dist
  export GOROOT="$(cd .. && pwd)"
  GOROOT_FINAL="${GOROOT_FINAL:-$GOROOT}"
  DEFGOROOT='-DGOROOT_FINAL="'"$GOROOT_FINAL"'"'

  mflag=""
  case "$GOHOSTARCH" in
    386) mflag=-m32;;
    amd64) mflag=-m64;;
  esac
  if [ "$(uname)" == "Darwin" ]; then
  	# golang.org/issue/5261
  	mflag="$mflag -mmacosx-version-min=10.6"
  fi
  run_helper ${CC:-gcc} $mflag -O2 -Wall -Werror -o cmd/dist/dist -Icmd/dist "$DEFGOROOT" cmd/dist/*.c
  run_helper eval $(./cmd/dist/dist env -p)
  echo

  # Build the compilers and go bootstrap tool
  echo "# Building compilers and Go bootstrap tool for host, $GOHOSTOS/$GOHOSTARCH."
  # Build go bootstrap
  bflags="-a -v"
  run_helper ./cmd/dist/dist bootstrap $bflags 

  # Delay move of dist tool to now, because bootstrap may clear tool directory.
  run_helper mkdir -p "$GOTOOLDIR"
  run_helper mv cmd/dist/dist "$GOTOOLDIR"/dist
  run_helper "$GOTOOLDIR"/go_bootstrap clean -i std
  echo
  
  # Build the packages and commands for the host using the bootstrap tool
  echo "# Building packages and commands for host, $GOHOSTOS/$GOHOSTARCH."
  bflags=""
  if [[ "$VERBOSE" = "1" ]]; then
    bflags="-x -p 1"
  fi
  prepare_host_env
  run_helper "$GOTOOLDIR"/go_bootstrap install \
              $bflags \
              -ccflags "$GO_CCFLAGS" \
              -gcflags "$GO_GCFLAGS" \
              -ldflags "$GO_LDFLAGS" \
              -v std
  restore_env
  echo
fi

# Building packages and commands for akaros 
if [[ "$BUILD" = "all" ]] ||
   [[ "$BUILD" = "akaros" ]] 
then
  echo "# Building packages and commands for $GOOS/$GOARCH."
  bflags=""
  if [[ "$VERBOSE" = "1" ]]; then
    bflags="-x -work -p 1"
  fi
  run_helper eval $($GOBIN/go tool dist env -p)
  prepare_target_env

  # Copy in the Akaros sycall header and generate kenc-style C defs from it
  run_helper cd "$GOROOT"/src/pkg/runtime/parlib
  run_helper cp "$ROSROOT"/kern/include/ros/bits/syscall.h syscall_${GOOS}.h
  run_helper_eval "$GOTOOLDIR/go_bootstrap tool cgo -cdefs types_${GOOS}.go > ztypes_${GOOS}.h"
  run_helper rm -rf _obj
  run_helper_eval "cd - > /dev/null"

  # Run the bootstrapping code to build the runtime
  run_helper $GOTOOLDIR/go_bootstrap install \
             $bflags \
             -ccflags "$GO_CCFLAGS" \
             -gcflags "$GO_GCFLAGS" \
             -ldflags "$GO_LDFLAGS" \
              -v runtime

# Install a wrapper script for building Go applications for Akaros on $GOARCH
cat > $GOBIN/go-$GOOS-$GOARCH << EOF
export GOOS=$GOOS
export GOARCH=$GOARCH
export CGO_ENABLED=1
export CC=$TARGETCC
export CXX=$TARGETCXX
export GO_LDFLAGS=$GO_LDFLAGS

ARGS=\$@
if [[ "\$1" = "build" ]] &&
   [[ "\$CC" != "" ]]
then
  ARGS="\$1 -ldflags $GO_LDFLAGS \${@:2}"
fi
$GOBIN/go \$ARGS
EOF
  chmod a+x $GOBIN/go-$GOOS-$GOARCH

  # Regenerate all of the files needed by the syscall package
  run_helper cd $GOROOT/src/pkg/syscall
  run_helper ./mkall_${GOOS}.sh > /dev/null 2>&1 
  run_helper_eval "cd - > /dev/null"

  # Run the bootstrapping code to build the rest of the Go packages
  run_helper $GOTOOLDIR/go_bootstrap install \
              $bflags \
              -ccflags "$GO_CCFLAGS" \
              -gcflags "$GO_GCFLAGS" \
              -ldflags "$GO_LDFLAGS" \
              -v std
  restore_env
  run_helper "$GOTOOLDIR"/dist banner
  echo
fi

