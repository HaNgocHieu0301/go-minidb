# Mini-DB Storage Engine (Golang)

Dự án xây dựng một hệ quản trị cơ sở dữ liệu (Storage Engine) từ con số không, tập trung vào bản chất của việc lưu trữ dữ liệu, cấu trúc cây B+ (B-Tree) và tính bền vững (Persistence).

Tài liệu tham khảo: 'Build Your Own Database From Scratch' ebook

---

## Tổng quan tiến độ (Project Tracking)

| Bước | Nội dung | Trạng thái |
|------|----------|------------|
| Bước 1 | Thiết kế định dạng Page/Node (Binary Layout) | ✅ Hoàn thành |
| Bước 2 | Logic cây B+ (Insert & Delete) | 🔄 Đang thực hiện |
| Bước 3 | Lưu trữ xuống đĩa (Persistence - mmap/fsync) | ⏳ Chưa bắt đầu |
| Bước 4 | Quản lý bộ nhớ trống (Free List) | ⏳ Chưa bắt đầu |
| **Milestone I** | **Working KV Store (Get/Set/Del persistent)** | ⏳ |
| Bước 5 | Tầng Relational DB (Rows & Columns) | ⏳ Chưa bắt đầu |
| Bước 6 | Range Query & Order-Preserving Encoding | ⏳ Chưa bắt đầu |
| Bước 7 | Secondary Index | ⏳ Chưa bắt đầu |
| Bước 8a | Atomic Transactions | ⏳ Chưa bắt đầu |
| Bước 8b | Concurrent Readers and Writers | ⏳ Chưa bắt đầu |
| Bước 9 | Query Language (Parser & Execution) | ⏳ Chưa bắt đầu |

---

## Cấu trúc dự án (Project Structure)

```
mini-db/
├── btree.go          # BNode: định dạng nhị phân và helper methods (Bước 1)
├── go.mod            # Go module
```

---

## Những gì đã hoàn thành (Bước 1)

Xây dựng xong "trái tim" của hệ thống lưu trữ: **Định dạng nhị phân cho một Node (Page Layout)**. Mọi Node được quản lý dưới dạng một mảng byte thô (`[]byte`) có kích thước cố định là 4096 bytes (4KB).

### Cấu trúc một Node (`BNode` Layout):

| Vùng | Kích thước | Ý nghĩa |
|------|------------|---------|
| Header | 4 bytes | `btype` (2B) + `nkeys` (2B) |
| Pointers | 8B × nkeys | Page number trỏ đến node con |
| Offsets | 2B × nkeys | Vị trí từng cặp KV trong vùng data |
| Key-Value Data | động | `klen (2B) \| vlen (2B) \| key \| value` |

### Các helper methods đã cài đặt:

- `btype()` / `setHeader()` — Đọc/ghi thông tin quản lý node
- `getPtr()` / `setPtr()` — Quản lý các liên kết giữa các node
- `getOffset()` / `setOffset()` — Xác định vị trí dữ liệu động trong node
- `getKey()` / `getVal()` — Trích xuất dữ liệu thực tế từ byte thô
- `kvPos()` / `nbytes()` — Tính vị trí và kích thước node

---

## Đang thực hiện (Bước 2)

Triển khai thuật toán cây B+. Mục tiêu: biến các mảng byte tĩnh thành cấu trúc dữ liệu có khả năng tự sắp xếp và tìm kiếm.

**Kế hoạch cụ thể:**
1. Định nghĩa struct `BTree` với callbacks `get`, `new`, `del`
2. Tìm kiếm (`nodeLookupLE`)
3. Chèn dữ liệu + xử lý Split node
4. Xóa dữ liệu + xử lý Merge node
5. Test với in-memory hashmap + reference map

---

## Thông số kỹ thuật (System Constants)

| Hằng số | Giá trị | Ý nghĩa |
|---------|---------|---------|
| `BTREE_PAGE_SIZE` | 4096 | Kích thước một trang đĩa tiêu chuẩn |
| `BTREE_MAX_KEY_SIZE` | 1000 | Giới hạn độ dài Key |
| `BTREE_MAX_VAL_SIZE` | 3000 | Giới hạn độ dài Value |
| Endianness | Little Endian | Quy chuẩn sắp xếp byte trong file |

> NOTE: Tài liệu này sẽ được cập nhật sau mỗi bước hoàn thành.
