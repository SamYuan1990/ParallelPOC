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

func benchmarke2e(concurrent int, b *testing.B) {
	b.ReportAllocs()
	p := &Pipeline{}
	p.Init()
	ProviderImpl := &ProviderImpl{
		Pipeline: p,
	}
	wg := &sync.WaitGroup{}
	Switcher := &Switcher{
		Pipeline: p,
		Wg:       wg,
		Count:    concurrent,
	}
	Switcher.Init()
	for i := 0; i < concurrent; i++ {
		c := &Consumer{
			Pipeline: p,
			Wg:       wg,
		}
		Switcher.AddConsumer(c)
		go c.Consume()
		defer c.Stop()
	}
	go Switcher.Switch()
	go ProviderImpl.Convert()
	defer ProviderImpl.Stop()
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

func BenchmarkCE2E1(b *testing.B) { benchmarke2e(1, b) }

func BenchmarkE2E2(b *testing.B) { benchmarke2e(2, b) }

func BenchmarkE2E4(b *testing.B) { benchmarke2e(4, b) }

func BenchmarkE2E8(b *testing.B) { benchmarke2e(8, b) }
