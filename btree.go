package gosqldb

import (
	"bytes"
	"encoding/binary"
)

// B+Tree ds
// goals
// - design a node format
// - manipulate nodes in a copy-on-write fashion
// - split and merge utility functions
// - tree insertion and deletion
//
//

// Node format
// type nkeys pointers offsets key-values unused
// | type | nkeys | pointers | offsets | key-values | unused |
// | 2B | 2B | nkeys * 8B | nkeys * 2B | ... | |
// this is the format of each KV pair. lengths followed by data
type BNode []byte

// we use same structure for both leaf and internal nodes, this wastes some space but it's easier to implement
//
// Header Size and Structure: The 4-byte header is composed of:
//
// - Node Type (2 bytes): Specifies if the node is a leaf or an internal node. Using uint16 (2 bytes) allows for representation of BNODE_NODE or BNODE_LEAF.
//
// - Number of Keys (2 bytes): Indicates how many keys are stored in the node. Using uint16 (2 bytes) allows the representation of the number of keys in the node.
const HEADER = 4
const BTREE_PAGE_SIZE = 4096 //bytes
const BTREE_MAX_KEY_SIZE = 1000 //bytes
const BTREE_MAX_VALUE_SIZE = 3000 //bytes

func init() {
	node1max := HEADER + 8 + 2 + 4 + BTREE_MAX_KEY_SIZE + BTREE_MAX_VALUE_SIZE
	// if node1max > BTREE_PAGE_SIZE {
	// 	panic(")
	// }
	assert(node1max <= BTREE_PAGE_SIZE, "node size exceeds max btree page size") //maximum KV
}

type BTree struct {
	root uint64

	get func(uint64) []byte //dereference a pointer
	new func([]byte) uint64 //allocate a new page
	del func(uint64)        //deallocate a page
}

//For an on-disk B+tree, the database file is an array of pages (nodes) referenced by page
// numbers (pointers). We’ll implement these callbacks as follows:
// get reads a page from disk.
// new allocates and writes a new page (copy-on-write).
// del deallocates a page.

// HEADER
const (
	BNODE_NODE = 1 // internal nodes without values
	BNODE_LEAF = 2 // leaf nodes with values
)

func (node BNode) btype() uint16 {
	return binary.LittleEndian.Uint16(node[0:2])
}

func (node BNode) nkeys() uint16 {
	return binary.LittleEndian.Uint16(node[2:4])
}

func (node BNode) setHeader(btype uint16, nkeys uint16) {
	binary.LittleEndian.PutUint16(node[0:2], btype)
	binary.LittleEndian.PutUint16(node[2:4], nkeys)
}

// Child pointers
func (node BNode) getPtr(idx uint16) uint64 {
	assert(idx < node.nkeys(), "idx is beyond no of keys")
	pos := HEADER + 8*idx
	return binary.LittleEndian.Uint64(node[pos:])
}

func offsetPos(node BNode, idx uint16) uint16 {
	return HEADER + 8*node.nkeys() + 2*(idx-1)
}
func (node BNode) getOffset(idx uint16) uint16 {
	if idx == 0 {
		return 0
	}
	return binary.LittleEndian.Uint16(node[offsetPos(node, idx):])
}

func (node BNode) setOffset(idx uint16, offset uint16) {
	//?
}

// key-values
func (node BNode) kvPos(idx uint16) uint16 {
	return HEADER + 8*node.nkeys() + 2*node.nkeys() + node.getOffset(idx)
}
func (node BNode) getKey(idx uint16) []byte {
	pos := node.kvPos(idx)
	klen := binary.LittleEndian.Uint16(node[pos:])
	return node[pos+4:][:klen]
}

func (node BNode) nbytes() uint16 {
	return node.kvPos(node.nkeys())
}

//KV lookups within a node

// The function is called nodeLookupLE because it uses the Less-than-or-Equal operator. For point
// queries, we should use the equal operator instead, which is a step we can add later.
func nodeLookupLE(node BNode, key []byte) uint16 {
	nkeys := node.nkeys()
	found := uint16(0)

	for i := uint16(1); i < nkeys; i++ {
		cmp := bytes.Compare(node.getKey(i), key)
		if cmp <= 0 {
			found = i
		}
		if cmp >= 0 {
			break
		}
	}
	return found
}


// Insert into leaf nodes
func leafInsert(
	new BNode, old BNode, idx uint16,
	key []byte, val[]byte
) {
	new.setHeader(BNODE_LEAF,  old.nkeys()+1) //setup the header
	nodeAppendRange(new, old, 0,0, idx)
	nodeAppendKV(new, idx, 0, key, val)
	nodeAppendRange(new, old, idx+1, idx, old.nkeys()-idx)
}
