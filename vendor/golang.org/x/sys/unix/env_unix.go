// Copyright 2010 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux netbsd openbsd solaris

// Unix environment variables.

package unix
import "github.com/equalll/mydebug"

import "syscall"

func Getenv(key string) (value string, found bool) {mydebug.INFO()
	return syscall.Getenv(key)
}

func Setenv(key, value string) error {mydebug.INFO()
	return syscall.Setenv(key, value)
}

func Clearenv() {mydebug.INFO()
	syscall.Clearenv()
}

func Environ() []string {mydebug.INFO()
	return syscall.Environ()
}
