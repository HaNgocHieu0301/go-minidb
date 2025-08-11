package minidb

import (
	"fmt"
	"os"
	"sync"

	"github.com/HaNgocHieu0301/go-minidb/internal/storage"
)

type Options struct {
	DataDir string
	// SyncEveryWrite: nếu true sẽ gọi fsync mỗi lần Put (an toàn hơn, chậm hơn).
	// Tạm thời bật mặc định true để đơn giản bền vững. Sau có thể tối ưu batch.
	SyncEveryWrite bool
}

type DB struct {
	opts Options

	mu  sync.RWMutex
	idx map[string][]byte // chỉ mục trong RAM (last write wins)
	log *storage.Log
}

func Open(opts Options) (*DB, error) {
	if opts.DataDir == "" {
		return nil, fmt.Errorf("Options.DataDir is empty")
	}
	if err := os.MkdirAll(opts.DataDir, 0o755); err != nil {
		return nil, err
	}

	log, err := storage.OpenLog(opts.DataDir)
	if err != nil {
		return nil, err
	}

	db := &DB{
		opts: opts,
		idx:  make(map[string][]byte),
		log:  log,
	}

	// Replay log -> xây lại idx
	if err := db.replay(); err != nil {
		_ = log.Close()
		return nil, err
	}
	return db, nil
}

func (db *DB) replay() error {
	return db.log.Iterate(func(k, v []byte) error {
		key := string(k) // copy key vào string làm map-key
		if len(v) == 0 {
			// Quy ước tạm: value rỗng là "delete"
			delete(db.idx, key)
			return nil
		}
		// Lưu bản copy của value để tránh bên ngoài sửa chung slice.
		valCopy := make([]byte, len(v))
		copy(valCopy, v)
		db.idx[key] = valCopy
		return nil
	})
}

func (db *DB) Close() error {
	db.mu.Lock()
	defer db.mu.Unlock()
	if db.log != nil {
		return db.log.Close()
	}
	return nil
}

func (db *DB) Put(key, value []byte) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if err := db.log.Append(key, value); err != nil {
		return err
	}
	if db.opts.SyncEveryWrite || db.opts.SyncEveryWrite == false {
		// Tạm thời mặc định sẽ là true.
		if err := db.log.Sync(); err != nil {
			return err
		}
	}
	valCopy := make([]byte, len(value))
	copy(valCopy, value)
	db.idx[string(key)] = valCopy
	return nil
}

func (db *DB) Get(key []byte) (value []byte, ok bool, err error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	v, ok := db.idx[string(key)]
	if !ok {
		return nil, false, nil
	}
	// Trả về bản copy bảo vệ internal state.
	out := make([]byte, len(v))
	copy(out, v)
	return out, true, nil
}

// Delete tạm thời: ghi record value rỗng và xoá khỏi idx.
func (db *DB) Delete(key []byte) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	if err := db.log.Append(key, nil); err != nil {
		return err
	}
	if err := db.log.Sync(); err != nil {
		return err
	}
	delete(db.idx, string(key))
	return nil
}
