package server

import (
	"fmt"

	"github.com/pawalt/kvstore/pkg/persist"
)

type GetRequest struct {
	Path []string
}

type GetResponse struct {
	Value []byte
}

func (k *KVServer) Get(req *GetRequest, resp *GetResponse) error {
	val := k.root.FindValue(req.Path)
	resp.Value = val
	return nil
}

type PutRequest struct {
	Path  []string
	Value []byte
}

type PutResponse struct {
}

func (k *KVServer) Put(req *PutRequest, resp *PutResponse) error {
	err := persist.WriteOp(k.file, k.writer, req.Path, req.Value)
	if err != nil {
		return fmt.Errorf("error while writing: %v", err)
	}

	k.root.Put(req.Path, req.Value)
	return nil
}
