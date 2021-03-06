Go on Akaros
============

This document serves as an overview of running Go on akaros.

Unlike most of the other OSs supported by Go, akaros does not yet support the
ability to build native Go binaries itself.  Instead, the entire Go setup is
hosted on a non-akaros machine, and all binaries are cross compiled. These
binaries are then shipped over to an akaros machine and executed.

As of Go 1.3, we are able to leverage the use of a `go_${GOOS}_${GOARCH}_exec`
script to allow us to invoke Go commands on our host machine but execute any
resulting Go binaries on a remote akaros machine.  So long as this script is in
our path, calls to things like "go run arguments..." or "go test arguments..."
do the job of compiling and building Go binaries for our desired target (akaros
in our case), but instead of just executing the resulting binary, the path of
the binary along with its arguments are passed to the `go_${GOOS}_${GOARCH}_exec`
script, and that script is executed instead.  We have simply implemented the
akaros version of this script to ship the resulting binary over to a remote
akaros machine and invoke it via an rpc call. For more information on this
script and how it is being used, take a look at `go run --help`.

While testing, our remote akaros machine is typically running inside an
instance of qemu on the same physical machine as the host.  While much of this
document assumes you are running in this environment, nothing prevents you from
running akaros on real hardware -- you may just need to tweak the default
settings a little, and we provide guidance on how to do so.

The rest of this document is dedicated to explaining the setup required and the
helper scripts provided to get you running Go on akaros quickly.

General Setup
-------------

This document assumes you have a few environment variables set up, as well as
some extra paths added to your PATH before invoking any of the scripts
described later on.

The required environment variables are:

	GOROOT      -- The path to your go installation
	GOPATH      -- The path to your external go workspace
	ROSROOT     -- The path to your akaros installation
	ROSXCCROOT  -- The path to your akaros cross-compiler installation

The required path elements are:

	PATH=$GOROOT/bin -- the path to where your go binaries exist
	                    (may be different if $GOBIN is set)
	PATH=$GOPATH/bin -- the bin path for your go workspace
	                    (may be different if $GOBIN is set)
	PATH=$GOROOT/misc/akaros/bin -- the path where the go_${GOOS}_${GOARCH}_exec
	                                exists. (Other scripts exist here too)
	PATH=$ROSXCCROOT/bin -- the path to the akaros x86_64 cross
	                        compiler's bin directory as installed as
	                        part of the akaros setup described below

The settings in my (klueska's) `.bashrc`, for example are:

	# Setup for go
	export GOROOT=$HOME/projects/go-akaros
	export GOPATH=$HOME/projects/go-workspace
	PATH="$GOROOT/bin:$PATH"
	PATH="$GOPATH/bin:$PATH"
	PATH="$GOROOT/misc/akaros/bin:$PATH"
	export PATH

	# Setup for Akaros
	export ROSROOT=$HOME/projects/akaros
	export ROSXCCROOT=$HOME/install/x86_64-ucb-akaros-gcc
	export QEMU_MONITOR_TTY=/dev/pts/29
	PATH="$ROSXCCROOT/bin:$PATH"
	export PATH

The extra `QEMU_MONITOR_TTY` environment variable is used by one of the scripts
described later on to point to the tty device where I run my qemu monitor for
qemu instances of akaros.  I typically have an instance of `screen` up and
running on a particular tty device and just attach and detach from this screen
as necessary to control my running instances of qemu.

Akaros Setup
------------

The full guide to installing and running akaros can be found at:

	https://github.com/brho/akaros/blob/master/GETTING_STARTED

For the purposes of just running Go, however, I've provided a quick start guide
that will walk you through the basic setup.  Following these steps will get you
set up with a Go-capable installation of akaros that you can run inside of qemu
on your host machine.

- Download akaros from github:

		git clone https://github.com/brho/akaros.git $ROSROOT

- Configure akaros with its default settings:

		cd $ROSROOT
		make ARCH=x86 defconfig

- Build and install the `x86_64` cross compiler for akaros:

		X86_64_INSTDIR=$ROSXCCROOT make xcc-upgrade-from-scratch

- Build busybox for akaros:

		cd $ROSROOT/tools/apps/busybox
		make x86_64

You should now have a working akaros installation on your host machine.  You
can test that everything is set up correctly by making sure the kernel compiles
and links properly:

	cd $ROSROOT
	make

You should see it compiling a bunch of source files, and ending with:

	LINK    obj/kern/akaros-kernel

Sometimes it may be required to update your akaros tree from github as changes
and updates are committed.  This may require you to re-install the cross
compiler and/or run various other make targets detailed in the full
installation guide. If you run into problems and don't know what to do next,
feel free to shoot an email to <akaros@lists.eecs.berkeley.edu> and we will be
happy to assist you.

Go Setup
--------

Once the akaros source tree and its `x86_64` cross compiler are installed, the
process for getting go set up is fairly straight forward.

- Download the go-akaros source tree:

		git clone git@github.com:klueska/go-akaros.git $GOROOT

- Build and install the go runtime for akaros:

		cd $GOROOT/src
		./akaros.bash make

As per normal, you will need to set `GOOS` and `GOARCH` appropriately in order to
build binaries for akaros. Additionally, I like to always set `CGO_ENABLED`
since lots of programs I tend to run actually use it, but it's not necessarily
required:

	export GOOS=akaros
	export GOARCH=amd64
	export CGO_ENABLED=1

A convenience script that you can 'source' to set these variables automatically
is located in `$GOROOT/misc/akaros/setup.sh` for your convenience:

	source $GOROOT/misc/akaros/setup.sh

You can test to make sure that everything is set up properly by going to
`$GOROOT/test` and trying to compile one of the binaries in there.  It
shouldn't actually run yet though because you are cross compiling for akaros
and we haven't yet completeled the steps to get you running binaries directly
on akaros.

	cd $GOROOT/test
	go build 64bit.go
	./64bit
		-bash: ./64bit: No such file or directory

And that's it!

Running Go on Akaros
--------------------

With akaros installed, and Go setup for cross compiling for akaros, these last
few steps will take you the rest of the way to actually running Go binaries on
your akaros installation.  They will be run via the standard `go run
arguments...` and `go test arguments...` commands from your host machine.

As mentioned at the beginning of this document, we leverage the use of a
`go_${GOOS}_${GOARCH}_exec` script to allow us to invoke Go commands on our
host machine but execute any resulting Go binaries on a remote akaros machine.

In the case where the remote machine is a qemu instance running on the host
machine, a couple of convenience scripts are provided in
`$GOROOT/misc/akaros/bin` to get you going quickly.

- Run the `go-9pserver.sh` script to start a 9p server that akaros can talk to:

		$GOROOT/misc/akaros/bin/go-akaros-9pserver.sh

- Launch an instance of qemu that runs the akaros kernel:

		$GOROOT/misc/akaros/bin/go-akaros-qemu.sh

- From within akaros itself, drop into the monitor, run busybox, and invoke the
  `go-bootstrap.sh` script:

		<ctrl-G>
		bb
		go-bootstrap.sh

This script may take a while, but make sure and wait for it to finish, and
print the line:

	listen started

Once this is ready, you will be able to invoke go commands on your host machine
and have them proxied through to the akaros instance for execution. The results
will be printed out on your host machine.

To test it, open up another terminal, set your `GOOS` and `GOARCH` appropriately
(or source `$GOROOT/misc/akaros/setup.sh`), and run the following:

	go test bufio

You should see something like the following output on your host machine:

	ok      bufio   3.637s

And that should be it!  Let me know if you have any problems or questions:  
	<akaros@lists.eecs.berkeley.edu>


