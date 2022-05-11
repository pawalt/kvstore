package server

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
	// have to do producer/consumer pattern so we don't get concurrent
	// file writes

	putOp := PutOp{
		req:      req,
		respChan: make(chan error),
	}

	k.putChan <- &putOp

	err := <-putOp.respChan
	if err != nil {
		return err
	}

	return nil
}
