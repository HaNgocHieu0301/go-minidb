package storage

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"path/filepath"
)

type Log struct {
	path string
	f    *os.File
	w    *bufio.Writer
}

// OpenLog mở/tạo data.log và chuẩn bị writer ở cuối file để append.
func OpenLog(dir string) (*Log, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	path := filepath.Join(dir, "data.log")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return nil, err
	}
	// Seek đến cuối file để ghi tiếp.
	if _, err := f.Seek(0, io.SeekEnd); err != nil {
		f.Close()
		return nil, err
	}
	return &Log{
		path: path,
		f:    f,
		w:    bufio.NewWriterSize(f, 64<<10),
	}, nil
}

// Append ghi một record key/value (chưa flush/sync).
func (l *Log) Append(key, value []byte) error {
	var hdr [binary.MaxVarintLen64 * 2]byte
	n := binary.PutUvarint(hdr[:], uint64(len(key)))
	m := binary.PutUvarint(hdr[n:], uint64(len(value)))
	if _, err := l.w.Write(hdr[:n+m]); err != nil {
		return err
	}
	if _, err := l.w.Write(key); err != nil {
		return err
	}
	if _, err := l.w.Write(value); err != nil {
		return err
	}
	return nil
}

// Sync đảm bảo dữ liệu xuống đĩa.
func (l *Log) Sync() error {
	if err := l.w.Flush(); err != nil {
		return err
	}
	return l.f.Sync()
}

func (l *Log) Close() error {
	if err := l.w.Flush(); err != nil {
		return err
	}
	return l.f.Close()
}

// Iterate đọc toàn bộ file từ đầu và gọi fn(key,value) cho từng record.
// Nếu đuôi file bị cắt dở -> bỏ qua phần đó (tạm thời, chưa có checksum).
func (l *Log) Iterate(fn func(k, v []byte) error) error {
	rf, err := os.Open(l.path)
	if err != nil {
		return err
	}
	defer rf.Close()

	r := bufio.NewReaderSize(rf, 64<<10)

	for {
		klen, err := binary.ReadUvarint(r)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			if errors.Is(err, io.ErrUnexpectedEOF) {
				return nil
			}
			return err
		}
		vlen, err := binary.ReadUvarint(r)
		if err != nil {
			if errors.Is(err, io.ErrUnexpectedEOF) {
				return nil
			}
			return err
		}
		k := make([]byte, klen)
		if _, err := io.ReadFull(r, k); err != nil {
			if errors.Is(err, io.ErrUnexpectedEOF) {
				return nil
			}
			return err
		}
		v := make([]byte, vlen)
		if _, err := io.ReadFull(r, v); err != nil {
			if errors.Is(err, io.ErrUnexpectedEOF) {
				return nil
			}
			return err
		}
		if err := fn(k, v); err != nil {
			return err
		}
	}
}
