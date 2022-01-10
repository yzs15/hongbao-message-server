package registry

import (
	"sync"

	"github.com/pkg/errors"

	"ict.ac.cn/hbmsgserver/pkg/wshub"
)

type Client struct {
	Mac      string
	Location string

	WsClient *wshub.Client
}

type Registry struct {
	mac2id    map[string]uint32
	id2client map[uint32]*Client

	idGenerator chan uint32

	mu sync.Mutex
}

func NewRegistry() *Registry {
	return &Registry{
		mac2id:      make(map[string]uint32),
		id2client:   make(map[uint32]*Client),
		idGenerator: make(chan uint32, 1000),
	}
}

func (r *Registry) Run() {
	var i uint32
	for i = 1; i < ^(uint32(0)); i++ {
		r.idGenerator <- i
	}
}

func (r *Registry) Register(c *Client) uint32 {
	mac := c.Mac

	r.mu.Lock()
	defer r.mu.Unlock()

	var id uint32
	var ok bool
	if id, ok = r.mac2id[mac]; !ok {
		id = <-r.idGenerator
		r.mac2id[mac] = id
	}
	r.id2client[id] = c

	return id
}

func (r *Registry) GetClient(id uint32) (*Client, error) {
	cli, ok := r.id2client[id]
	if !ok {
		return nil, errors.New("no this client")
	}

	if cli.WsClient.Down {
		return nil, errors.New("this client is down")
	}

	return cli, nil
}
