DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
GOROOT="$DIR/../../.."
SRCDIR="$GOROOT/src"
if [ -f $SRCDIR/local.bash ]; then
    source $SRCDIR/local.bash
fi
export GOROOT
$DIR/mkall.sh $@
rm -rf _obj
