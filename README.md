# MiniDB (Go) — LSM-tree Key–Value Store

> Learning-first storage engine. Small, testable, documented.

## Tính năng

-

## Kiến trúc

- **LSM-tree**: ghi tuần tự (append) → MemTable → flush thành SSTable bất biến → compaction.
- **API**: `Open/Close`, `Put/Get/Delete`, `Iterator(start,end)`.

## Quick Start

### Yêu cầu

- Go **1.2x**.
- GNU Make.

### Cài & chạy

```bash
# clone
git clone https://github.com/HaNgocHieu0301/go-minidb
cd go-minidb

# build & test
make fmt && make test && make build

# chạy bản server khởi động
go run ./cmd/minidbd -v
```
### Kiểm thử

```bash
make test      # unit tests
make race      # chạy test với -race
```

## Cấu trúc dự án

```
/cmd/minidbd/           # binary server
/cmd/minidb-cli/        # client CLI
/pkg/minidb/            # public API: DB, Options, Iterator
/internal/storage/      # WAL, SSTable, compaction, file utils
/internal/memtable/     # skiplist/map + iterator
/internal/server/       # HTTP/TCP handlers
/internal/manifest/     # theo dõi level, file metadata
/internal/testutil/     # sinh dữ liệu, chai benchmark
/docs/                  # tài liệu thiết kế & BRD
```

## Roadmap chi tiết

Xem **_docs/business-requirement.md_** (BRD) để biết phạm vi, FR/NFR, rủi ro và cột mốc từng sprint.

## Changelog
- **2025-08-11 – Sprint 1**: Append-only log + replay map; `Put/Get/Delete` tối giản; thêm test.
- **2025-08-11 – Sprint 0**: Khởi tạo repo, module, API stubs, Makefile/CI; thêm BRD & cập nhật README.
