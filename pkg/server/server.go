package server

import (
	"bufio"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"

	"github.com/pawalt/kvstore/pkg/kv"
	"github.com/pawalt/kvstore/pkg/persist"
)

const (
	PORT     = 1337
	PROTOCOL = "tcp"
)

type KVServer struct {
	root   kv.KVNode
	writer *bufio.Writer
	file   *os.File
}

func New(filePath string) (*KVServer, error) {
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}

	r := bufio.NewReader(f)
	rootNode, err := persist.Restore(r)
	if err != nil {
		return nil, err
	}

	w := bufio.NewWriter(f)

	return &KVServer{
		root:   rootNode,
		writer: w,
	}, nil
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
