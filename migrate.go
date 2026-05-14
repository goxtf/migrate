// Package migrate provides database migration functionality.
// It is a fork of golang-migrate/migrate with additional features and improvements.
package migrate

import (
	"errors"
	"fmt"
	"os"
	"sync"
)

// ErrNoChange is returned when no migration is needed.
var ErrNoChange = errors.New("no change")

// ErrNilVersion is returned when the version is nil.
var ErrNilVersion = errors.New("nil version")

// ErrInvalidVersion is returned when the version is invalid.
var ErrInvalidVersion = errors.New("invalid version")

// ErrLocked is returned when the database is locked.
var ErrLocked = errors.New("database locked")

// ErrLockTimeout is returned when the lock times out.
var ErrLockTimeout = errors.New("lock timeout")

// DefaultPrefetchMigrations is the default number of migrations to prefetch.
const DefaultPrefetchMigrations = 10

// DefaultLockTimeout is the default timeout for acquiring a lock in seconds.
const DefaultLockTimeout = 15

// Migrate is the main struct for managing database migrations.
type Migrate struct {
	// sourceName is the registered source driver name.
	sourceName string
	// sourceDrv is the source driver instance.
	sourceDrv Source

	// databaseName is the registered database driver name.
	databaseName string
	// databaseDrv is the database driver instance.
	databaseDrv Database

	// Log is an optional logger.
	Log Logger

	// GracefulStop is a channel to signal a graceful stop.
	GracefulStop chan bool
	isGracefulStop bool

	isLockedMu *sync.Mutex
	isLocked   bool

	// PrefetchMigrations is the number of migrations to prefetch.
	PrefetchMigrations uint

	// LockTimeout is the timeout in seconds for acquiring a lock.
	LockTimeout uint
}

// Logger is the interface for logging migration activity.
type Logger interface {
	Printf(format string, v ...interface{})
	Verbose() bool
}

// New returns a new Migrate instance from the provided source and database URLs.
func New(sourceURL, databaseURL string) (*Migrate, error) {
	m := &Migrate{
		GracefulStop:       make(chan bool, 1),
		PrefetchMigrations: DefaultPrefetchMigrations,
		LockTimeout:        DefaultLockTimeout,
		isLockedMu:         &sync.Mutex{},
	}

	sourceDrv, err := newSource(sourceURL, m)
	if err != nil {
		return nil, fmt.Errorf("source: %w", err)
	}
	m.sourceDrv = sourceDrv

	databaseDrv, err := newDatabase(databaseURL, m)
	if err != nil {
		return nil, fmt.Errorf("database: %w", err)
	}
	m.databaseDrv = databaseDrv

	return m, nil
}

// Close closes the source and database connections.
func (m *Migrate) Close() (source error, database error) {
	databaseSrvClose := make(chan error)
	sourceSrvClose := make(chan error)

	go func() {
		databaseSrvClose <- m.databaseDrv.Close()
	}()

	go func() {
		sourceSrvClose <- m.sourceDrv.Close()
	}()

	return <-sourceSrvClose, <-databaseSrvClose
}

// logPrintf logs a formatted message if a logger is configured.
func (m *Migrate) logPrintf(format string, v ...interface{}) {
	if m.Log != nil {
		m.Log.Printf(format, v...)
	}
}

// logVerbosePrintf logs a formatted message if verbose logging is enabled.
func (m *Migrate) logVerbosePrintf(format string, v ...interface{}) {
	if m.Log != nil && m.Log.Verbose() {
		m.Log.Printf(format, v...)
	}
}

// newSource creates a new source driver from the given URL.
func newSource(url string, m *Migrate) (Source, error) {
	_ = url
	_ = m
	// Source driver registration and instantiation will be implemented
	// as source drivers are added to the project.
	return nil, fmt.Errorf("no source driver registered for url: %s", url)
}

// newDatabase creates a new database driver from the given URL.
func newDatabase(url string, m *Migrate) (Database, error) {
	_ = url
	_ = m
	// Database driver registration and instantiation will be implemented
	// as database drivers are added to the project.
	return nil, fmt.Errorf("no database driver registered for url: %s", url)
}

// isGracefulStopSet returns true if the graceful stop signal has been received.
func (m *Migrate) isGracefulStopSet() bool {
	select {
	case <-m.GracefulStop:
		return true
	default:
		return false
	}
}

// stderr writes a message to standard error.
func stderr(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", v...)
}
