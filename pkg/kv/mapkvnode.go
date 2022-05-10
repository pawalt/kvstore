package kv

type MapKVNode struct {
	children map[string]*MapKVNode
	value    []byte
}

func NewMapVKNode() *MapKVNode {
	return &MapKVNode{
		children: map[string]*MapKVNode{},
		value:    nil,
	}
}

func (k *MapKVNode) Get(name string) KVNode {
	return k.children[name]
}

func (k *MapKVNode) Value() []byte {
	return k.value
}

func (k *MapKVNode) Children() []KVNode {
	ret := make([]KVNode, 0, len(k.children))
	for _, val := range k.children {
		ret = append(ret, val)
	}
	return ret
}

func (k *MapKVNode) FindValue(path []string) []byte {
	if len(path) == 0 {
		return k.value
	}

	if child, ok := k.children[path[0]]; ok {
		return child.FindValue(path[1:])
	} else {
		return nil
	}
}

func (k *MapKVNode) Put(path []string, val []byte) {
	if len(path) == 0 {
		k.value = val
		return
	}

	if child, ok := k.children[path[0]]; ok {
		child.Put(path[1:], val)
	} else {
		child := NewMapVKNode()
		child.Put(path[1:], val)
		k.children[path[0]] = child
	}
}
