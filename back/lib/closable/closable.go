// This package prints error when defer *.Close() returns it.
package closable

import "fmt"

// Closeable is a wrapper which helps to `defer SafeClose()`.
type Closeable interface {
	Close() error
}

// SafeClose is a function, performing a safe close. Intended to be used inside `defer`;
func SafeClose(c Closeable) {
	if err := c.Close(); err != nil && err.Error() != "invalid argument" {
		fmt.Printf("Error %T: %v\n", err, err)
	}
}
