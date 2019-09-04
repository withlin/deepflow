package idmap

import (
	. "encoding/binary"
	"testing"
)

const (
	bytesKeyLen = 40
)

func newNode(key0, key1 uint64) *u320IDMapNode {
	node := &u320IDMapNode{}
	BigEndian.PutUint64(node.key[bytesKeyLen-16-8:bytesKeyLen-16], key0)
	BigEndian.PutUint64(node.key[bytesKeyLen-8:], key1)
	node.hash = uint32(key0>>32) ^ uint32(key0) ^ uint32(key1>>32) ^ uint32(key1)
	return node
}

func TestU320IDMapAddOrGet(t *testing.T) {
	m := NewU320IDMap(1024)

	exp := true
	node := newNode(0, 1)
	if _, ret := m.AddOrGet(node.key[:], node.hash, 1, false); ret != exp {
		t.Errorf("第一次插入，Expected %v found %v", exp, ret)
	}
	exp = false
	if _, ret := m.AddOrGet(node.key[:], node.hash, 2, false); ret != exp {
		t.Errorf("插入同样的值，Expected %v found %v", exp, ret)
	}
	if ret, _ := m.Get(node.key[:], node.hash); ret != 1 {
		t.Errorf("查找失败，Expected %v found %v", 1, ret)
	}
	exp = false
	if _, ret := m.AddOrGet(node.key[:], node.hash, 2, true); ret != exp {
		t.Errorf("插入同样的值，Expected %v found %v", exp, ret)
	}
	if ret, _ := m.Get(node.key[:], node.hash); ret != 2 {
		t.Errorf("查找失败，Expected %v found %v", 2, ret)
	}
	exp = true
	node = newNode(1, 0)
	if _, ret := m.AddOrGet(node.key[:], node.hash, 1, false); ret != exp {
		t.Errorf("插入不同的值，Expected %v found %v", exp, ret)
	}

	if m.Size() != 2 {
		t.Errorf("当前长度，Expected %v found %v", 2, m.Size())
	}
}

func TestU320IDMapSize(t *testing.T) {
	m := NewU320IDMap(1024)
	if m.Size() != 0 {
		t.Errorf("当前长度，Expected %v found %v", 0, m.Size())
	}

	node := newNode(0, 1)
	m.AddOrGet(node.key[:], node.hash, 1, false)
	if m.Size() != 1 {
		t.Errorf("当前长度，Expected %v found %v", 1, m.Size())
	}
	m.AddOrGet(node.key[:], node.hash, 1, false)
	if m.Size() != 1 {
		t.Errorf("当前长度，Expected %v found %v", 1, m.Size())
	}
	node = newNode(0, 2)
	m.AddOrGet(node.key[:], node.hash, 1, false)
	if m.Size() != 2 {
		t.Errorf("当前长度，Expected %v found %v", 2, m.Size())
	}
	node = newNode(1, 0)
	m.AddOrGet(node.key[:], node.hash, 1, false)
	if m.Size() != 3 {
		t.Errorf("当前长度，Expected %v found %v", 3, m.Size())
	}
}

func TestU320IDMapGet(t *testing.T) {
	m := NewU320IDMap(1024)

	node := newNode(0, 1)
	m.AddOrGet(node.key[:], node.hash, 1, false)
	if _, in := m.Get(node.key[:], node.hash); !in {
		t.Errorf("查找失败")
	}
	node = newNode(0, 2)
	if _, in := m.Get(node.key[:], node.hash); in {
		t.Errorf("查找失败")
	}
	node = newNode(1, 0)
	if _, in := m.Get(node.key[:], node.hash); in {
		t.Errorf("查找失败")
	}
	m.AddOrGet(node.key[:], node.hash, 1, false)
	if _, in := m.Get(node.key[:], node.hash); !in {
		t.Errorf("查找失败")
	}
}

func TestU320IDMapClear(t *testing.T) {
	m := NewU320IDMap(4)

	node := newNode(0, 1)
	m.AddOrGet(node.key[:], node.hash, 1, false)
	m.AddOrGet(node.key[:], node.hash, 1, false)
	node = newNode(0, 2)
	m.AddOrGet(node.key[:], node.hash, 1, false)
	node = newNode(1, 0)
	m.AddOrGet(node.key[:], node.hash, 1, false)
	m.Clear()
	if m.Size() != 0 {
		t.Errorf("当前长度，Expected %v found %v", 0, m.Size())
	}
	node = newNode(0, 1)
	m.AddOrGet(node.key[:], node.hash, 1, false)
	if _, in := m.Get(node.key[:], node.hash); !in {
		t.Errorf("查找失败")
	}
	if m.Size() != 1 {
		t.Errorf("当前长度，Expected %v found %v", 1, m.Size())
	}
}

func BenchmarkU320IDMap(b *testing.B) {
	m := NewU320IDMap(1 << 26)
	nodes := make([]*u320IDMapNode, (b.N+3)/4*4)

	for i := uint64(0); i < uint64(b.N); i += 4 {
		// 构造哈希冲突
		nodes[i] = newNode(i, i<<1)
		nodes[i+1] = newNode(i<<1, i)
		nodes[i+2] = newNode(^i, ^(i << 1))
		nodes[i+3] = newNode(^(i << 1), ^i)
	}

	b.ResetTimer()
	for i := uint64(0); i < uint64(b.N); i += 4 {
		m.AddOrGet(nodes[i].key[:], nodes[i].hash, uint32(i<<2), false)
		m.AddOrGet(nodes[i+1].key[:], nodes[i+1].hash, uint32(i<<2), false)
		m.AddOrGet(nodes[i+2].key[:], nodes[i+2].hash, uint32(i<<2), false)
		m.AddOrGet(nodes[i+3].key[:], nodes[i+3].hash, uint32(i<<2), false)
	}
	b.Logf("size=%d, width=%d", m.Size(), m.Width())
}
