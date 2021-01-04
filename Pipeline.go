package pipelinepoc

import (
	"github.com/enriquebris/goconcurrentqueue"
	lru "github.com/hashicorp/golang-lru"
)

const Len = 10240

type Pipeline struct {
	FIFO    *goconcurrentqueue.FIFO
	Current *lru.Cache
}

func (p *Pipeline) Init() error {
	var err error
	p.FIFO = goconcurrentqueue.NewFIFO()
	p.Current, err = lru.New(Len)
	err = p.FIFO.Enqueue(p.Current)
	return err
}

func (p *Pipeline) Add(key, value interface{}) {
	p.Current.Add(key, value)
}

func (p *Pipeline) AddWithNewLRU(key, value interface{}) error {
	var err error
	p.Current, err = lru.New(Len)
	err = p.FIFO.Enqueue(p.Current)
	p.Add(key, value)
	return err
}

func (p *Pipeline) RemoveOldest() (key interface{}, value interface{}, ok bool) {
	v, err := p.FIFO.Get(0)
	if err != nil {
		return nil, nil, false
	}
	if v == p.Current {
		if v.(*lru.Cache).Len() == 0 {
			return nil, nil, false
		} else {
			return v.(*lru.Cache).RemoveOldest()
		}
	} else {
		for v.(*lru.Cache).Len() == 0 {
			p.FIFO.Remove(0)
			v, _ = p.FIFO.Get(0)
			if v == p.Current {
				return v.(*lru.Cache).RemoveOldest()
			}
		}
		return v.(*lru.Cache).RemoveOldest()
	}
}
