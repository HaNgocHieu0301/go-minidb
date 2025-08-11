# docs/business-requirement.md

> **MiniDB – Business Requirements (BRD)**\
> Phiên bản: v0.1 (Sprint 0) — Cập nhật: 2025-08-11

## 1) Bối cảnh & mục tiêu

- **Bối cảnh:** Side project học tập của một junior SWE để hiểu sâu cơ chế của storage engine/DB.
- **Mục tiêu:** Xây dựng một **LSM-tree Key–Value Store** nhỏ gọn bằng Go, có độ bền dữ liệu, có test/benchmark cơ bản, có tài liệu rõ ràng.
- **Tinh thần:** Mỗi sprint hoàn thành một lát cắt nhỏ, có giải thích chi tiết, dễ hiểu, có kiểm thử.

## 2) Người dùng mục tiêu & Ứng dụng

- **Người dùng mục tiêu:** người học, SWE muốn hiểu DB internals.
- **Use cases:**
    - Lưu trữ KV cục bộ bền vững cho demo/lab.
    - Chạy một service HTTP nhỏ để thao tác KV từ tool/CLI.
    - So sánh/benchmark và đọc mã nguồn để học.

## 3) Phạm vi (in-scope) & ngoài phạm vi (out-of-scope)

**In-scope (giai đoạn đầu):**

- KV API: `Put/Get/Delete`, iterator scan theo range/prefix.
- WAL, MemTable, flush thành SSTable, compaction, crash-recovery.
- Concurrency cơ bản (RWMutex), nền compaction đơn giản.
- Service HTTP/TCP nhỏ (PUT/GET/DELETE/SCAN).
- Quan sát: log đơn giản, metrics/prom, pprof.

**Out-of-scope (ban đầu):**

- SQL parser, join phức tạp, MVCC đầy đủ, replication/raft, transactions đa bản ghi.

## 4) Yêu cầu chức năng (FR)

- FR-1: `Put(key,value)` ghi WAL trước, cập nhật MemTable.
- FR-2: `Get(key)` hợp nhất kết quả từ MemTable + nhiều SSTable.
- FR-3: `Delete(key)` tạo **tombstone**; thể hiện đúng dù còn bản cũ trong SSTable.
- FR-4: **Flush** MemTable → SSTable có index khối (block index) để binary search.
- FR-5: **Compaction**: gộp file, loại bỏ version cũ/tombstone đã già.
- FR-6: **Iterator** k-way merge, hỗ trợ range scan & prefix.
- FR-7: **Crash-recovery**: đọc WAL, bỏ qua partial record bằng checksum.
- FR-8: **Service** HTTP/TCP; CLI tương tác cơ bản.

## 5) Yêu cầu phi chức năng (NFR)

- NFR-1: **Độ bền**: WAL + ghi an toàn (fsync/atomic rename khi cần).
- NFR-2: **Tính nhất quán**: last-write-wins; đọc sau ghi cùng tiến trình là thấy (read-your-writes) khi không lỗi.
- NFR-3: **Hiệu năng**: tối ưu theo LSM (append, sequential I/O). Benchmark cơ bản sẽ bổ sung ở Sprint 12.
- NFR-4: **Đơn giản – dễ đọc**: mã nguồn, docs, test phải dễ hiểu; ưu tiên giáo dục hơn tối ưu cực hạn.
- NFR-5: **Tính di động**: chạy trên Linux/macOS/Windows; Go ≥ 1.22.

## 6) Tiêu chí thành công

- Chạy được end-to-end, có test xanh (bao gồm `-race`).
- Có tài liệu giải thích thiết kế, trade-off; benchmark tối thiểu.
- Có thể demo bằng CLI/HTTP trong <5 phút.

## 7) Kế hoạch & cột mốc

- Sprint 0: Skeleton dự án, API stubs, CI cơ bản.
- Sprint 1: File append-only tối thiểu → `data.log` + restart get đúng.
- Sprint 2: WAL + MemTable + khôi phục từ WAL.
- Sprint 3: Flush → SSTable (sorted, block index).
- Sprint 4: Tombstone & semantics delete.
- Sprint 5: Compaction.
- Sprint 6: Concurrency (RWMutex) + `-race` xanh.
- Sprint 7: Iterator & Range scan.
- Sprint 8: Cấu hình & Bloom filter (tuỳ chọn).
- Sprint 9: Crash-recovery hoàn chỉnh (checksum/atomic rename).
- Sprint 10: Service hoá (HTTP/TCP) + CLI.
- Sprint 11: Quan sát (metrics/log/pprof).
- Sprint 12: Docs + benchmark + polish.

## 8) Rủi ro & biện pháp

- **Racing/Deadlock** → dùng `-race`, review lock order, test song song.
- **Hỏng dữ liệu do partial write** → checksum record + atomic rename.
- **Compaction bug làm mất key** → test hợp nhất, property/quickcheck/fuzz.
- **Độ phức tạp tăng** → giữ module nhỏ, ghi chú design ngắn mỗi sprint.

## 9) Thuật ngữ (glossary)

- **WAL**: Write-Ahead Log — ghi trước khi áp dụng vào cấu trúc dữ liệu.
- **MemTable**: cấu trúc in-memory giữ các cập nhật mới (map/skiplist).
- **SSTable**: Sorted String Table — file bất biến, dữ liệu đã sort + index.
- **Tombstone**: đánh dấu xoá logic; được dọn trong compaction.
- **Compaction**: gộp/tái sắp xếp SSTable để giảm phân mảnh & xoá bản cũ.

## 10) Trạng thái tiến độ

-

## 11) Cách cập nhật tài liệu này

- Cập nhật **mục 10 (checklist)** và **mục 7 (kế hoạch)** sau mỗi sprint.
- Ghi chú các quyết định kỹ thuật quan trọng ngay dưới mục liên quan.

## 12) Ghi chú thiết kế — Sprint 1

- **Record format:** `uvarint(keyLen) | uvarint(valueLen) | key | value`.
- **Durability:** mỗi `Put/Delete` gọi `Flush + fsync` (vì `SyncEveryWrite=true`).
- **Replay:** đọc toàn file, áp dụng "last-write-wins"; `valueLen==0` xem như delete tối giản.
- **Hạn chế hiện tại:** chưa có checksum, chưa tách WAL riêng, chưa có MemTable/SSTable/Compaction → file log sẽ phình to. Các sprint 2–5 sẽ xử lý.
- **Lý do chọn:** tối giản để đảm bảo persistence + dễ kiểm thử khởi động lại.