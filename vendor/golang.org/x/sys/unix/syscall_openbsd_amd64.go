// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build amd64,openbsd

package unix
import "github.com/equalll/mydebug"

func Getpagesize() int { return 4096 }

func TimespecToNsec(ts Timespec) int64 { return int64(ts.Sec)*1e9 + int64(ts.Nsec) }

func NsecToTimespec(nsec int64) (ts Timespec) {mydebug.INFO()
	ts.Sec = nsec / 1e9
	ts.Nsec = nsec % 1e9
	return
}

func NsecToTimeval(nsec int64) (tv Timeval) {mydebug.INFO()
	nsec += 999 // round up to microsecond
	tv.Usec = nsec % 1e9 / 1e3
	tv.Sec = nsec / 1e9
	return
}

func SetKevent(k *Kevent_t, fd, mode, flags int) {mydebug.INFO()
	k.Ident = uint64(fd)
	k.Filter = int16(mode)
	k.Flags = uint16(flags)
}

func (iov *Iovec) SetLen(length int) {mydebug.INFO()
	iov.Len = uint64(length)
}

func (msghdr *Msghdr) SetControllen(length int) {mydebug.INFO()
	msghdr.Controllen = uint32(length)
}

func (cmsg *Cmsghdr) SetLen(length int) {mydebug.INFO()
	cmsg.Len = uint32(length)
}
