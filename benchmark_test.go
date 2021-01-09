package pipelinepoc

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

func ConstructKey(key string, value int) Key {
	return Key{
		Name:    key,
		Version: value,
	}
}

func ConstructBlock(key string, value int) *BlockImpl {
	Wkeys := make([]Key, 0)
	keys := ConstructKey(key, value)
	Wkeys = append(Wkeys, keys)
	tximpl := &TxImpl{
		Wkeys: Wkeys,
	}
	txs := make([]*TxImpl, 0)
	txs = append(txs, tximpl)

	BlockImpl := &BlockImpl{
		Txs: txs,
	}
	return BlockImpl
}

func BenchmarkProvider(b *testing.B) {
	b.ReportAllocs()
	p := &Pipeline{}
	p.Init()
	ProviderImpl := &ProviderImpl{
		Pipeline: p,
	}
	go ProviderImpl.Convert()
	defer ProviderImpl.Stop()
	b.ResetTimer()
	go func() {
		for i := 0; i < b.N; i++ {
			BlockImpl := ConstructBlock(strconv.Itoa(i), rand.Intn(1000))
			p.Comming.Enqueue(BlockImpl)
		}
	}()
	for i := 0; i < b.N; {
		time.Sleep(100 * time.Millisecond)
		v, err := p.Output.Dequeue()
		if err == nil && v != nil {
			i++
		}
	}
	b.StopTimer()
}

func BenchmarkSingle(b *testing.B) {
	b.ReportAllocs()
	p := &Pipeline{}
	p.Init()
	ProviderImpl := &ProviderImpl{
		Pipeline: p,
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	c := &Consumer{
		Pipeline: p,
		Wg:       wg,
	}
	go ProviderImpl.Convert()
	go c.Consume()
	defer ProviderImpl.Stop()
	defer c.Stop()
	b.ResetTimer()
	go func() {
		for i := 0; i < b.N; i++ {
			BlockImpl := ConstructBlock(strconv.Itoa(i), rand.Intn(1000))
			p.Comming.Enqueue(BlockImpl)
		}
	}()
	for i := 0; i < b.N; {
		time.Sleep(100 * time.Millisecond)
		v, err := p.Output.Dequeue()
		if err == nil && v != nil {
			for !v.(*BlockImpl).Txs[0].Processed {

			}
			i++
		}
	}
	b.StopTimer()
}

func BenchmarkParallel2(b *testing.B) {
	b.ReportAllocs()
	p := &Pipeline{}
	p.Init()
	ProviderImpl := &ProviderImpl{
		Pipeline: p,
	}
	wg := &sync.WaitGroup{}
	c := &Consumer{
		Pipeline: p,
		Wg:       wg,
	}
	c1 := &Consumer{
		Pipeline: p,
		Wg:       wg,
	}
	Switcher := &Switcher{
		Pipeline: p,
		Wg:       wg,
		Count:    2,
	}
	Switcher.Init()
	Switcher.AddConsumer(c)
	Switcher.AddConsumer(c1)
	go Switcher.Switch()
	go ProviderImpl.Convert()
	go c.Consume()
	go c1.Consume()
	defer ProviderImpl.Stop()
	defer c.Stop()
	defer c1.Stop()
	defer Switcher.Stop()
	b.ResetTimer()
	go func() {
		for i := 0; i < b.N; i++ {
			BlockImpl := ConstructBlock(strconv.Itoa(i), rand.Intn(1000))
			p.Comming.Enqueue(BlockImpl)
		}
	}()
	for i := 0; i < b.N; {
		time.Sleep(100 * time.Millisecond)
		v, err := p.Output.Dequeue()
		if err == nil && v != nil {
			for !v.(*BlockImpl).Txs[0].Processed {

			}
			i++
		}
	}
	b.StopTimer()
}
