package pipelinepoc

import (
	"sync"

	"github.com/enriquebris/goconcurrentqueue"
	lru "github.com/hashicorp/golang-lru"
)

const Len = 10240

type Pipeline struct {
	// a comming queue for comming blocks received
	Comming *goconcurrentqueue.FIFO
	// current processing lru
	// it should be a map for key tx mapping
	// it should be concurrent
	// it should be support get oldest as FIFO
	// so it is LRU map
	Current *lru.Cache
	// the next will be used when current is logiclly full
	// when current is finisehd processed, the next and current will be exchanged.(similar with java young gc)
	Next *lru.Cache
	// the output should be able to support block tx mapping
	// should be able to return as FIFO ways
	Output *goconcurrentqueue.FIFO

	PCurrent *lru.Cache
	PNext    *lru.Cache

	CCurrent *lru.Cache
	CNext    *lru.Cache

	lock sync.Mutex
}

func (p *Pipeline) Init() error {
	var err error
	p.Comming = goconcurrentqueue.NewFIFO()
	p.Output = goconcurrentqueue.NewFIFO()
	p.Current, err = lru.New(Len)
	if err != nil {
		return err
	}
	p.Next, err = lru.New(Len)
	if err != nil {
		return err
	}
	p.PCurrent = p.Current
	p.PNext = p.Next
	p.CCurrent = p.Current
	p.CNext = p.Next
	return nil
}

func (p *Pipeline) SwitchP() {
	p.lock.Lock()
	defer p.lock.Unlock()
	GetLogger().Info("Switch P")
	tmp := p.PCurrent
	p.PCurrent = p.PNext
	p.PNext = tmp
}

func (p *Pipeline) SwitchC() {
	p.lock.Lock()
	defer p.lock.Unlock()
	GetLogger().Info("Switch C")
	//if !p.SwitchAble() {
	tmp := p.CCurrent
	p.CCurrent = p.CNext
	p.CNext = tmp
	//}
}

// if Pcurrent same with Ccurrent, able to switch
func (p *Pipeline) SwitchAble() bool {
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.CCurrent == p.PCurrent {
		return true
	}
	return false
}
