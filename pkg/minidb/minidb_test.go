package minidb

import "testing"

func TestOpenClose(t *testing.T) {
	dir := t.TempDir()
	db, err := Open(Options{DataDir: dir})
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("Close failed: %v", err)
	}
}
