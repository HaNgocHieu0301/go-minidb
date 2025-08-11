package minidb

// Options định nghĩa cấu hình cơ bản của DB.
type Options struct {
	DataDir string
}

type DB struct {
	// TODO: memtable, sstable manager, wal...
	opts Options
}

// Open khởi tạo DB với cấu hình tối thiểu.
func Open(opts Options) (*DB, error) {
	return &DB{opts: opts}, nil
}

func (db *DB) Close() error { return nil }

func (db *DB) Put(key, value []byte) error { return nil }
func (db *DB) Get(key []byte) (value []byte, ok bool, err error) {
	return nil, false, nil
}
func (db *DB) Delete(key []byte) error { return nil }
