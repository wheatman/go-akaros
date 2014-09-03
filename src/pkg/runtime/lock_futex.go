// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build akaros dragonfly freebsd linux

package runtime

import "unsafe"

// This implementation depends on OS-specific implementations of
//
//	runtime·futexsleep(uint32 *addr, uint32 val, int64 ns)
//		Atomically,
//			if(*addr == val) sleep
//		Might be woken up spuriously; that's allowed.
//		Don't sleep longer than ns; ns < 0 means forever.
//
//	runtime·futexwakeup(uint32 *addr, uint32 cnt)
//		If any procs are sleeping on addr, wake up at most cnt.

const (
	mutex_unlocked = 0
	mutex_locked   = 1
	mutex_sleeping = 2

	active_spin     = 4
	active_spin_cnt = 30
	passive_spin    = 1
)

// Possible lock states are mutex_unlocked, mutex_locked and mutex_sleeping.
// mutex_sleeping means that there is presumably at least one sleeping thread.
// Note that there can be spinning threads during all states - they do not
// affect mutex's state.

func futexsleep(addr *uint32, val uint32, ns int64)
func futexwakeup(addr *uint32, cnt uint32)

// We use the uintptr mutex.key and note.key as a uint32.
func key32(p *uintptr) *uint32 {
	return (*uint32)(unsafe.Pointer(p))
}

func lock(l *mutex) {
	gp := getg()

	if gp.m.locks < 0 {
		gothrow("runtime·lock: lock count")
	}
	gp.m.locks++

	// Speculative grab for lock.
	v := xchg(key32(&l.key), mutex_locked)
	if v == mutex_unlocked {
		return
	}

	// wait is either MUTEX_LOCKED or MUTEX_SLEEPING
	// depending on whether there is a thread sleeping
	// on this mutex.  If we ever change l->key from
	// MUTEX_SLEEPING to some other value, we must be
	// careful to change it back to MUTEX_SLEEPING before
	// returning, to ensure that the sleeping thread gets
	// its wakeup call.
	wait := v

	// On uniprocessors, no point spinning.
	// On multiprocessors, spin for ACTIVE_SPIN attempts.
	spin := 0
	if ncpu > 1 {
		spin = active_spin
	}
	for {
		// Try for lock, spinning.
		for i := 0; i < spin; i++ {
			for l.key == mutex_unlocked {
				if cas(key32(&l.key), mutex_unlocked, wait) {
					return
				}
			}
			procyield(active_spin_cnt)
		}

		// Try for lock, rescheduling.
		for i := 0; i < passive_spin; i++ {
			for l.key == mutex_unlocked {
				if cas(key32(&l.key), mutex_unlocked, wait) {
					return
				}
			}
			osyield()
		}

		// Sleep.
		v = xchg(key32(&l.key), mutex_sleeping)
		if v == mutex_unlocked {
			return
		}
		wait = mutex_sleeping
		futexsleep(key32(&l.key), mutex_sleeping, -1)
	}
}

func unlock(l *mutex) {
	v := xchg(key32(&l.key), mutex_unlocked)
	if v == mutex_unlocked {
		gothrow("unlock of unlocked lock")
	}
	if v == mutex_sleeping {
		futexwakeup(key32(&l.key), 1)
	}

	gp := getg()
	gp.m.locks--
	if gp.m.locks < 0 {
		gothrow("runtime·unlock: lock count")
	}
	if gp.m.locks == 0 && gp.preempt { // restore the preemption request in case we've cleared it in newstack
		gp.stackguard0 = stackPreempt
	}
}

// One-time notifications.
func noteclear(n *note) {
	n.key = 0
}

func notewakeup(n *note) {
	old := xchg(key32(&n.key), 1)
	if old != 0 {
		print("notewakeup - double wakeup (", old, ")\n")
		gothrow("notewakeup - double wakeup")
	}
	futexwakeup(key32(&n.key), 1)
}

func notesleep(n *note) {
	gp := getg()
	if gp != gp.m.g0 {
		gothrow("notesleep not on g0")
	}
	for atomicload(key32(&n.key)) == 0 {
		gp.m.blocked = true
		futexsleep(key32(&n.key), 0, -1)
		gp.m.blocked = false
	}
}

//go:nosplit
func notetsleep_internal(n *note, ns int64) bool {
	gp := getg()

	if ns < 0 {
		for atomicload(key32(&n.key)) == 0 {
			gp.m.blocked = true
			futexsleep(key32(&n.key), 0, -1)
			gp.m.blocked = false
		}
		return true
	}

	if atomicload(key32(&n.key)) != 0 {
		return true
	}

	deadline := nanotime() + ns
	for {
		gp.m.blocked = true
		futexsleep(key32(&n.key), 0, ns)
		gp.m.blocked = false
		if atomicload(key32(&n.key)) != 0 {
			break
		}
		now := nanotime()
		if now >= deadline {
			break
		}
		ns = deadline - now
	}
	return atomicload(key32(&n.key)) != 0
}

func notetsleep(n *note, ns int64) bool {
	gp := getg()
	if gp != gp.m.g0 && gp.m.gcing == 0 {
		gothrow("notetsleep not on g0")
	}

	return notetsleep_internal(n, ns)
}

// same as runtime·notetsleep, but called on user g (not g0)
// calls only nosplit functions between entersyscallblock/exitsyscall
func notetsleepg(n *note, ns int64) bool {
	gp := getg()
	if gp == gp.m.g0 {
		gothrow("notetsleepg on g0")
	}

	entersyscallblock()
	ok := notetsleep_internal(n, ns)
	exitsyscall()
	return ok
}