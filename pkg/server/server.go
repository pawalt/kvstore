package server

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
	"time"

	"github.com/pawalt/kvstore/pkg/kv"
	"github.com/pawalt/kvstore/pkg/persist"
)

const (
	PORT             = 1337
	PROTOCOL         = "tcp"
	PUT_CHAN_SIZE    = 100
	WRITE_QUEUE_SIZE = 50
)

type KVServer struct {
	root     kv.KVNode
	writer   *bufio.Writer
	file     *os.File
	putChan  chan (*PutOp)
	stopChan chan (struct{})
}

type PutOp struct {
	Req      *PutRequest
	RespChan chan (error)
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
	go k.HandleWrites()

	select {}
}

func (k *KVServer) HandleWrites() error {
	toWrite := make([]*PutOp, 0)

	for {
		select {
		case <-k.stopChan:
			return nil
		case putOp := <-k.putChan:
			toWrite = append(toWrite, putOp)
		case <-time.After(time.Second / 4):
			err := k.ExecuteWrites(toWrite)
			for _, op := range toWrite {
				op.RespChan <- err
			}
			toWrite = toWrite[:0]
		}

		if len(toWrite) >= WRITE_QUEUE_SIZE {
			err := k.ExecuteWrites(toWrite)
			for _, op := range toWrite {
				op.RespChan <- err
			}
			toWrite = toWrite[:0]
		}
	}
}

func (k *KVServer) ExecuteWrites(ops []*PutOp) error {
	converted := make([]*persist.Write, 0, len(ops))
	for _, op := range ops {
		converted = append(converted, &persist.Write{
			Path:  op.Req.Path,
			Value: op.Req.Value,
		})
	}

	err := persist.BatchWrite(k.file, k.writer, converted)
	if err != nil {
		return fmt.Errorf("failed trying to write in batch: %v", err)
	}

	for _, op := range ops {
		k.root.Put(op.Req.Path, op.Req.Value)
	}

	return nil
}

func (k *KVServer) Put(path []string, data []byte) error {
	// have to do producer/consumer pattern so we don't get concurrent
	// file writes

	putOp := PutOp{
		Req: &PutRequest{
			Path:  path,
			Value: data,
		},
		RespChan: make(chan error),
	}

	k.putChan <- &putOp

	err := <-putOp.RespChan
	if err != nil {
		return err
	}

	return nil
}
