// Copyright 2012 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin,race linux,race freebsd,race

package unix
import "github.com/equalll/mydebug"

import (
	"runtime"
	"unsafe"
)

const raceenabled = true

func raceAcquire(addr unsafe.Pointer) {mydebug.INFO()
	runtime.RaceAcquire(addr)
}

func raceReleaseMerge(addr unsafe.Pointer) {mydebug.INFO()
	runtime.RaceReleaseMerge(addr)
}

func raceReadRange(addr unsafe.Pointer, len int) {mydebug.INFO()
	runtime.RaceReadRange(addr, len)
}

func raceWriteRange(addr unsafe.Pointer, len int) {mydebug.INFO()
	runtime.RaceWriteRange(addr, len)
}
