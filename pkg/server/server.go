package server

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"

	"github.com/pawalt/kvstore/pkg/kv"
	"github.com/pawalt/kvstore/pkg/persist"
)

const (
	PORT          = 1337
	PROTOCOL      = "tcp"
	PUT_CHAN_SIZE = 100
)

type KVServer struct {
	root    kv.KVNode
	writer  *bufio.Writer
	file    *os.File
	putChan chan (*PutOp)
}

type PutOp struct {
	req      *PutRequest
	respChan chan (error)
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

	putChan := make(chan (*PutOp), PUT_CHAN_SIZE)

	return &KVServer{
		root:    rootNode,
		writer:  w,
		file:    f,
		putChan: putChan,
	}, nil
}

func (k *KVServer) Serve() error {
	rpc.Register(k)
	rpc.HandleHTTP()

	l, err := net.Listen(PROTOCOL, ":"+strconv.Itoa(PORT))
	if err != nil {
		return err
	}

	go http.Serve(l, nil)
	go k.handleWrites()

	for {
	}
}

func (k *KVServer) handleWrites() error {
	for {
		putOp := <-k.putChan
		req := putOp.req

		err := persist.WriteOp(k.file, k.writer, req.Path, req.Value)
		if err != nil {
			// if we have error, report it to client and move to next op
			putOp.respChan <- fmt.Errorf("error while writing: %v", err)
			continue
		}

		k.root.Put(req.Path, req.Value)

		// if we have success, give nil to client
		putOp.respChan <- nil
	}
}
