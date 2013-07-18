package os
import "runtime/parlib"

func force_parlib_import() {
	parlib.Init()
}
