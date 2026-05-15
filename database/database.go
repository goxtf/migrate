// Package database provides the interface and registry for database drivers
// used by the migrate tool.
package database

import (
	"fmt"
	"io"
	"sync"
)

// Driver is the interface that database drivers must implement.
// It defines the core operations needed to apply and track migrations.
type Driver interface {
	// Open returns a new driver instance configured with the given URL.
	Open(url string) (Driver, error)

	// Close releases any resources held by the driver.
	Close() error

	// Lock acquires an advisory lock on the database to prevent concurrent
	// migrations from running simultaneously.
	Lock() error

	// Unlock releases the advisory lock.
	Unlock() error

	// Run applies a single migration step from the given reader.
	Run(migration io.Reader) error

	// SetVersion stores the current migration version and whether the
	// database is in a dirty state (i.e., a migration failed mid-run).
	SetVersion(version int, dirty bool) error

	// Version returns the currently applied migration version.
	// Returns -1 if no migration has been applied yet.
	// The dirty flag indicates whether the last migration failed.
	Version() (version int, dirty bool, err error)

	// Drop deletes all database objects (tables, views, etc.).
	// This is a destructive operation and should be used with caution.
	Drop() error
}

var (
	driversMu sync.RWMutex
	drivers   = make(map[string]Driver)
)

// Register makes a database driver available under the given name.
// It panics if called twice with the same name or if the driver is nil.
func Register(name string, driver Driver) {
	driversMu.Lock()
	defer driversMu.Unlock()

	if driver == nil {
		panic("database: Register called with nil driver")
	}
	if _, dup := drivers[name]; dup {
		panic("database: Register called twice for driver " + name)
	}
	drivers[name] = driver
}

// Open returns a new database driver instance for the given driver name and URL.
// The driver must have been previously registered via Register.
func Open(name, url string) (Driver, error) {
	driversMu.RLock()
	driver, ok := drivers[name]
	driversMu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("database: unknown driver %q (forgotten import?)", name)
	}

	return driver.Open(url)
}

// List returns the names of all registered database drivers.
func List() []string {
	driversMu.RLock()
	defer driversMu.RUnlock()

	list := make([]string, 0, len(drivers))
	for name := range drivers {
		list = append(list, name)
	}
	return list
}

// ErrLocked is returned when the database is already locked by another process.
var ErrLocked = fmt.Errorf("database: lock already acquired")

// ErrNilVersion is returned when no migration version has been set.
var ErrNilVersion = fmt.Errorf("database: no migration version set")

// ErrDirty is returned when the database is in a dirty state,
// meaning a previous migration failed and must be resolved manually.
type ErrDirty struct {
	Version int
}

func (e ErrDirty) Error() string {
	return fmt.Sprintf("database: dirty migration version %d — resolve manually or use force", e.Version)
}
