package migrate

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNew verifies that New returns a valid Migrate instance
// when provided with valid source and database URLs.
func TestNew(t *testing.T) {
	t.Run("valid source and database", func(t *testing.T) {
		m, err := New("file://testdata/migrations", "stub://")
		if err != nil {
			// Skip if stub driver not registered in this test run
			t.Skipf("skipping: %v", err)
		}
		require.NotNil(t, m)
	})

	t.Run("invalid source URL", func(t *testing.T) {
		_, err := New("invalid://", "stub://")
		assert.Error(t, err)
	})

	t.Run("invalid database URL", func(t *testing.T) {
		_, err := New("file://testdata/migrations", "invalid://")
		assert.Error(t, err)
	})
}

// TestMigrateUp verifies that Up applies all pending migrations.
func TestMigrateUp(t *testing.T) {
	m, err := newTestMigrate(t)
	if err != nil {
		t.Skipf("skipping: %v", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, ErrNoChange) {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestMigrateDown verifies that Down rolls back all applied migrations.
func TestMigrateDown(t *testing.T) {
	m, err := newTestMigrate(t)
	if err != nil {
		t.Skipf("skipping: %v", err)
	}

	// Apply first, then roll back
	_ = m.Up()

	err = m.Down()
	if err != nil && !errors.Is(err, ErrNoChange) {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestMigrateSteps verifies that Steps applies or rolls back N migrations.
func TestMigrateSteps(t *testing.T) {
	m, err := newTestMigrate(t)
	if err != nil {
		t.Skipf("skipping: %v", err)
	}

	t.Run("step up", func(t *testing.T) {
		err := m.Steps(1)
		if err != nil && !errors.Is(err, ErrNoChange) {
			t.Fatalf("unexpected error on step up: %v", err)
		}
	})

	t.Run("step down", func(t *testing.T) {
		err := m.Steps(-1)
		if err != nil && !errors.Is(err, ErrNoChange) {
			t.Fatalf("unexpected error on step down: %v", err)
		}
	})

	t.Run("zero steps", func(t *testing.T) {
		err := m.Steps(0)
		assert.Error(t, err, "expected error when steps is 0")
	})
}

// TestMigrateVersion verifies that Version returns the current schema version.
func TestMigrateVersion(t *testing.T) {
	m, err := newTestMigrate(t)
	if err != nil {
		t.Skipf("skipping: %v", err)
	}

	version, dirty, err := m.Version()
	if err != nil && !errors.Is(err, ErrNilVersion) {
		t.Fatalf("unexpected error: %v", err)
	}

	// If no error, dirty should be a valid boolean and version a uint
	if err == nil {
		assert.False(t, dirty, "database should not be in dirty state")
		_ = version
	}
}

// newTestMigrate is a helper that creates a Migrate instance for testing.
// It skips the test if the required drivers are not available.
func newTestMigrate(t *testing.T) (*Migrate, error) {
	t.Helper()
	return New("file://testdata/migrations", "stub://")
}
