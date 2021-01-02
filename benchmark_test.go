package pipelinepoc

import (
	"math/rand"
	"strconv"
	"testing"

	lru "github.com/hashicorp/golang-lru"
)

func ConstructKey(key string, value int) Key {
	return Key{
		Name:    key,
		Version: value,
	}
}

func ConstructNodeWithSingleKey(key string, value int) *Node {
	Wkeys0 := make([]Key, 0)
	key0 := ConstructKey(key, value)
	Wkeys0 = append(Wkeys0, key0)
	return &Node{
		Wkeys: Wkeys0,
	}
}

func benchmarkNConsumer(concurrent int, maxMsgCount int, b *testing.B) {
	LRU, _ := lru.New(b.N)
	output := make(chan *Node, b.N)
	defer close(output)
	ArrayConsumer := make([]Consumer, concurrent)

	ProviderImpl := &ProviderImpl{
		LRU: LRU,
	}
	for i := 0; i < concurrent; i++ {
		c := Consumer{
			LRU: LRU,
		}
		go c.Consume(output)
		ArrayConsumer = append(ArrayConsumer, c)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Wkeys := make([]Key, 0)
		for i := 0; i < maxMsgCount; i++ {
			key := Key{
				Name: strconv.Itoa(rand.Intn(100000)),
			}
			Wkeys = append(Wkeys, key)
		}
		tximpl := TxImpl{
			Wkeys: Wkeys,
		}
		txs := make([]TxImpl, 0)
		txs = append(txs, tximpl)
		BlockImpl := &BlockImpl{
			Txs: txs,
		}
		ProviderImpl.Convert(BlockImpl)
	}
	var n int
	for n < b.N {
		<-output
		n++
	}
	b.StopTimer()
	if LRU.Len() != 0 || len(output) != 0 {
		b.Fatal()
	}
	for _, vc := range ArrayConsumer {
		vc.Stop()
	}
	/*err := goleak.Find(goleak.IgnoreTopFunction("github.com/onsi/ginkgo/internal/specrunner.(*SpecRunner).registerForInterrupts"))
	if err != nil {
		b.Fatal(err)
	}*/
}

/*
func BenchmarkSerial(b *testing.B) {
	Pipeline := &Pipeline{}
	Pipeline.Init()
	output := make(chan *Node, 10)
	defer close(output)

	c := &Consumer{
		Pipeline: Pipeline,
	}
	go c.Consume(output)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ProviderImpl := &ProviderImpl{
			Pipeline: Pipeline,
		}
		Wkeys := make([]Key, 0)
		key := Key{
			Name: "abc",
		}
		Wkeys = append(Wkeys, key)

		tximpl := TxImpl{
			Wkeys: Wkeys,
		}
		txs := make([]TxImpl, 0)
		txs = append(txs, tximpl)
		BlockImpl := &BlockImpl{
			Txs: txs,
		}
		ProviderImpl.Convert(BlockImpl)
	}
	var n int
	for n < b.N {
		<-output
		n++
		fmt.Println(n)
	}
	b.StopTimer()
	c.Stop()
}
*/
func BenchmarkConsumer1(b *testing.B)    { benchmarkNConsumer(1, 10, b) }
func BenchmarkConsumer2(b *testing.B)    { benchmarkNConsumer(2, 10, b) }
func BenchmarkConsumer8(b *testing.B)    { benchmarkNConsumer(8, 10, b) }
func BenchmarkConsumer864(b *testing.B)  { benchmarkNConsumer(8, 64, b) }
func BenchmarkConsumer8128(b *testing.B) { benchmarkNConsumer(8, 128, b) }
