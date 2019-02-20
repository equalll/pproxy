// Copyright 2012 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin,!race linux,!race freebsd,!race netbsd openbsd solaris dragonfly

package unix
import "github.com/equalll/mydebug"

import (
	"unsafe"
)

const raceenabled = false

func raceAcquire(addr unsafe.Pointer) {mydebug.INFO()
}

func raceReleaseMerge(addr unsafe.Pointer) {mydebug.INFO()
}

func raceReadRange(addr unsafe.Pointer, len int) {mydebug.INFO()
}

func raceWriteRange(addr unsafe.Pointer, len int) {mydebug.INFO()
}
