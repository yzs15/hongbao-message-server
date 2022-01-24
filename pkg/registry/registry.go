package registry

import (
	"sync"
	"time"

	"ict.ac.cn/hbmsgserver/pkg/timeutils"

	"ict.ac.cn/hbmsgserver/pkg/idutils"

	"ict.ac.cn/hbmsgserver/pkg/msgserver"

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

	hub *wshub.Hub
}

func NewRegistry(hub *wshub.Hub) *Registry {
	return &Registry{
		mac2id:      make(map[string]uint32),
		id2client:   make(map[uint32]*Client),
		idGenerator: make(chan uint32, 1000),
		hub:         hub,
	}
}

func (r *Registry) Run() {
	var i uint32
	for i = 4; i < ^(uint32(0)); i++ {
		r.idGenerator <- i
	}
}

func (r *Registry) Register(c *Client, expId uint32) uint32 {
	r.mu.Lock()
	defer r.mu.Unlock()

	var finalId uint32 = 0
	mac := c.Mac
	if id, ok := r.mac2id[mac]; ok {
		finalId = id

	} else if _, ok := r.id2client[expId]; expId > 0 && !ok {
		finalId = expId

	} else {
		for id := range r.idGenerator {
			if _, ok := r.id2client[id]; ok {
				continue
			}
			finalId = id
			break
		}
	}

	r.mac2id[mac] = finalId
	r.id2client[finalId] = c

	return finalId
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

func (r *Registry) Broadcast(msg msgserver.Message) time.Time {
	sendTime := timeutils.GetPtpTime()
	r.hub.Broadcast <- msg
	return sendTime
}

func (r *Registry) Send(msg msgserver.Message) (time.Time, error) {
	cliID := idutils.CliId32(msg.Receiver())
	cli, err := r.GetClient(cliID)
	if err != nil {
		return time.Time{}, err
	}
	sendTime := timeutils.GetPtpTime()
	cli.WsClient.Send <- msg
	return sendTime, nil
}
