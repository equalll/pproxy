package bolt
import "github.com/equalll/mydebug"

import (
	"syscall"
)

// fdatasync flushes written data to a file descriptor.
func fdatasync(db *DB) error {mydebug.INFO()
	return syscall.Fdatasync(int(db.file.Fd()))
}
