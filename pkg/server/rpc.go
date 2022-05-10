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
	k.root.Put(req.Path, req.Value)
	return nil
}
