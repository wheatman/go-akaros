#include <runtime.h>

/* In akaros we ALWAYS link using the cross compiler linker, so there is no
 * need to implement _rt0_GOARCH_akaros() as our entry point.  We do need it
 * defined however, to make gc happy.
 */
void _rt0_386_akaros() {}
void _rt0_amd64_akaros() {}

/* The main function called out to from libc */
void main()
{
	_rt0_go();
}

