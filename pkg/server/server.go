package server

import (
	"net"
	"net/http"
	"net/rpc"
	"strconv"

	"github.com/pawalt/kvstore/pkg/kv"
)

const (
	PORT     = 1337
	PROTOCOL = "tcp"
)

type KVServer struct {
	root kv.KVNode
}

func New() *KVServer {
	return &KVServer{
		root: kv.NewMapVKNode(),
	}
}

func (k *KVServer) Serve() error {
	rpc.Register(k)
	rpc.HandleHTTP()

	l, err := net.Listen(PROTOCOL, ":"+strconv.Itoa(PORT))
	if err != nil {
		return err
	}

	http.Serve(l, nil)

	return nil
}
