package minidb

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPersistAfterRestart(t *testing.T) {
	dir := t.TempDir()

	db, err := Open(Options{DataDir: dir, SyncEveryWrite: true})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer db.Close()

	kvs := map[string]string{
		"alpha": "1",
		"beta":  "2",
		"gamma": "3",
	}
	for k, v := range kvs {
		if err := db.Put([]byte(k), []byte(v)); err != nil {
			t.Fatalf("put %s: %v", k, err)
		}
	}

	if err := db.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	// Mở lại
	db2, err := Open(Options{DataDir: dir, SyncEveryWrite: true})
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	defer db2.Close()

	for k, v := range kvs {
		got, ok, err := db2.Get([]byte(k))
		if err != nil || !ok || string(got) != v {
			t.Fatalf("get %s: got(%q, ok=%v, err=%v) want(%q,true,nil)", k, got, ok, err, v)
		}
	}
}

func TestLastWriteWins(t *testing.T) {
	dir := t.TempDir()

	db, err := Open(Options{DataDir: dir, SyncEveryWrite: true})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer db.Close()

	key := []byte("dup")
	if err := db.Put(key, []byte("v1")); err != nil {
		t.Fatal(err)
	}
	if err := db.Put(key, []byte("v2")); err != nil {
		t.Fatal(err)
	}
	v, ok, err := db.Get(key)
	if err != nil || !ok || string(v) != "v2" {
		t.Fatalf("get after overwrite: got(%q, ok=%v, err=%v) want(v2,true,nil)", v, ok, err)
	}

	// restart và kiểm tra vẫn "v2"
	db.Close()
	db2, err := Open(Options{DataDir: dir, SyncEveryWrite: true})
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	defer db2.Close()

	v2, ok, err := db2.Get(key)
	if err != nil || !ok || string(v2) != "v2" {
		t.Fatalf("after restart: got(%q, ok=%v, err=%v) want(v2,true,nil)", v2, ok, err)
	}
}

func TestDeleteMinimal(t *testing.T) {
	dir := t.TempDir()

	db, err := Open(Options{DataDir: dir, SyncEveryWrite: true})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer db.Close()

	key := []byte("todel")
	if err := db.Put(key, []byte("x")); err != nil {
		t.Fatal(err)
	}
	if err := db.Delete(key); err != nil {
		t.Fatal(err)
	}

	_, ok, err := db.Get(key)
	if err != nil || ok {
		t.Fatalf("get after delete: ok=%v err=%v, want ok=false", ok, err)
	}

	// restart và kiểm tra vẫn không tồn tại
	db.Close()
	db2, err := Open(Options{DataDir: dir, SyncEveryWrite: true})
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	defer db2.Close()

	_, ok, err = db2.Get(key)
	if err != nil || ok {
		t.Fatalf("after restart: ok=%v err=%v, want ok=false", ok, err)
	}
}

// (Tuỳ chọn) kiểm tra file thực sự được tạo.
func TestDataLogCreated(t *testing.T) {
	dir := t.TempDir()
	db, err := Open(Options{DataDir: dir, SyncEveryWrite: true})
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if _, err := os.Stat(filepath.Join(dir, "data.log")); err != nil {
		t.Fatalf("data.log not found: %v", err)
	}
}
