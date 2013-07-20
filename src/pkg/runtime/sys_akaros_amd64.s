// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
// System calls and other sys stuff for Akaros amd64
//

// set tls base to DI
TEXT runtime·settls(SB),7,$32
	SUBQ	$12, SP
	MOVL	$0, 0(SP)  // vcore 0
	MOVQ	DI, 4(SP)  // the new fs base
	CALL	runtime∕parlib·Set_tls_desc(SB)
	RET

