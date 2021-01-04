package pipelinepoc

import (
	"math/rand"
	"strconv"
	"testing"
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
	output := make(chan *Node, b.N)
	defer close(output)
	ArrayConsumer := make([]Consumer, concurrent)
	p := &Pipeline{}
	p.Init()
	ProviderImpl := &ProviderImpl{
		Pipeline: p,
	}
	for i := 0; i < concurrent; i++ {
		c := Consumer{
			Pipeline: p,
		}
		go c.Consume(output)
		ArrayConsumer = append(ArrayConsumer, c)
	}
	b.ReportAllocs()
	b.ResetTimer()
	go func() {
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
	}()
	var n int
	for n < b.N {
		<-output
		n++
	}
	b.StopTimer()
	if p.Current.Len() != 0 || len(output) != 0 {
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

func BenchmarkConsumer1(b *testing.B)   { benchmarkNConsumer(1, 10, b) }
func BenchmarkConsumer2(b *testing.B)   { benchmarkNConsumer(2, 10, b) }
func BenchmarkConsumer4(b *testing.B)   { benchmarkNConsumer(4, 10, b) }
func BenchmarkConsumer8(b *testing.B)   { benchmarkNConsumer(8, 10, b) }
func BenchmarkConsumer115(b *testing.B) { benchmarkNConsumer(1, 15, b) }
func BenchmarkConsumer120(b *testing.B) { benchmarkNConsumer(1, 20, b) }
