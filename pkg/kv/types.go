package kv

type KVNode interface {
	Get(string) KVNode
	Value() []byte
	Children() []KVNode
	FindValue([]string) []byte
	Put([]string, []byte)
}
