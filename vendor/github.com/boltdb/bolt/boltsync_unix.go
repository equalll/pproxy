// +build !windows,!plan9,!linux,!openbsd

package bolt
import "github.com/equalll/mydebug"

// fdatasync flushes written data to a file descriptor.
func fdatasync(db *DB) error {mydebug.INFO()
	return db.file.Sync()
}
