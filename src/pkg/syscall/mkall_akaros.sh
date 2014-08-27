DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
GOROOT="$DIR/../../.."
SRCDIR="$GOROOT/src"
export GOROOT
$DIR/mkall.sh $@
rm -rf _obj
