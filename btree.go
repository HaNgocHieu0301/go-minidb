package minidb

import "encoding/binary"

// Cấu trúc một node (page 4096 byte cố định):
// | type | nkeys | pointers        | offsets        | key-values |
// | 2B   | 2B    | nkeys * 8B      | nkeys * 2B     | ...        |
//
// Định dạng một cặp KV:
// | klen | vlen | key | val |
// | 2B   | 2B   | ... | ... |

type BNode struct {
	data []byte // mảng byte thô, ánh xạ trực tiếp lên một disk page
}

const (
	BNODE_NODE = 1 // internal node (chỉ chứa key, không có value)
	BNODE_LEAF = 2 // node lá (chứa cả key và value)
)

const (
	HEADER             = 4    // type (2B) + nkeys (2B)
	BTREE_PAGE_SIZE    = 4096 // kích thước một disk page
	BTREE_MAX_KEY_SIZE = 1000
	BTREE_MAX_VAL_SIZE = 3000
)

// --- Header ---

func (node BNode) btype() uint16 {
	return binary.LittleEndian.Uint16(node.data[0:2])
}

func (node BNode) nkeys() uint16 {
	return binary.LittleEndian.Uint16(node.data[2:4])
}

func (node BNode) setHeader(btype uint16, nkeys uint16) {
	binary.LittleEndian.PutUint16(node.data[0:2], btype)
	binary.LittleEndian.PutUint16(node.data[2:4], nkeys)
}

// --- Pointers (chỉ dùng cho internal node) ---

func (node BNode) getPtr(idx uint16) uint64 {
	pos := HEADER + 8*idx
	return binary.LittleEndian.Uint64(node.data[pos : pos+8])
}

func (node BNode) setPtr(idx uint16, val uint64) {
	pos := HEADER + 8*idx
	binary.LittleEndian.PutUint64(node.data[pos:pos+8], val)
}

// --- Offsets ---

// offsetPos trả về vị trí byte của offset thứ idx trong mảng offset.
// Mảng offset bắt đầu ngay sau mảng pointer.
// idx tính từ 1 vì offset[0] luôn bằng 0 và không được lưu trên đĩa.
func offsetPos(node BNode, idx uint16) uint16 {
	return HEADER + 8*node.nkeys() + 2*(idx-1)
}

// getOffset trả về offset của cặp KV thứ idx so với đầu vùng dữ liệu KV.
// offset[0] luôn bằng 0 (theo quy ước, không lưu trên đĩa).
func (node BNode) getOffset(idx uint16) uint16 {
	if idx == 0 {
		return 0
	}
	pos := offsetPos(node, idx)
	return binary.LittleEndian.Uint16(node.data[pos : pos+2])
}

func (node BNode) setOffset(idx uint16, offset uint16) {
	pos := offsetPos(node, idx)
	binary.LittleEndian.PutUint16(node.data[pos:pos+2], offset)
}

// --- Key-Value ---

// kvPos trả về vị trí byte tuyệt đối của cặp KV thứ idx trong node.data.
func (node BNode) kvPos(idx uint16) uint16 {
	return HEADER + 8*node.nkeys() + 2*node.nkeys() + node.getOffset(idx)
}

func (node BNode) getKey(idx uint16) []byte {
	pos := node.kvPos(idx)
	klen := binary.LittleEndian.Uint16(node.data[pos : pos+2])
	return node.data[pos+4 : pos+4+klen]
}

func (node BNode) getVal(idx uint16) []byte {
	pos := node.kvPos(idx)
	klen := binary.LittleEndian.Uint16(node.data[pos : pos+2])
	vlen := binary.LittleEndian.Uint16(node.data[pos+2 : pos+4])
	return node.data[pos+4+klen : pos+4+klen+vlen]
}

// nbytes trả về tổng số byte đang dùng của node.
// Dùng để kiểm tra node có cần tách (split) không khi vượt quá BTREE_PAGE_SIZE.
func (node BNode) nbytes() uint16 {
	return node.kvPos(node.nkeys())
}
