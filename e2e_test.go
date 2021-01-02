package pipelinepoc_test

import (
	"github.com/SamYuan1990/pipelinepoc"
	lru "github.com/hashicorp/golang-lru"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/goleak"
)

var _ = Describe("e2e", func() {
	It("should pass", func() {
		LRU, _ := lru.New(128)
		output := make(chan *pipelinepoc.Node, 10)
		defer close(output)
		ProviderImpl := &pipelinepoc.ProviderImpl{
			LRU: LRU,
		}
		Wkeys := make([]pipelinepoc.Key, 0)
		key := ConstructKey("key", 0)
		Wkeys = append(Wkeys, key)
		tximpl := pipelinepoc.TxImpl{
			Wkeys: Wkeys,
		}
		txs := make([]pipelinepoc.TxImpl, 0)
		txs = append(txs, tximpl)

		BlockImpl := &pipelinepoc.BlockImpl{
			Txs: txs,
		}
		ProviderImpl.Convert(BlockImpl)
		Expect(LRU.Len()).Should(Equal(1))
		c := &pipelinepoc.Consumer{
			LRU: LRU,
		}
		go c.Consume(output)
		c1 := &pipelinepoc.Consumer{
			LRU: LRU,
		}
		go c1.Consume(output)
		<-output
		c1.Stop()
		c.Stop()
		Expect(LRU.Len()).Should(Equal(0))
		Expect(goleak.Find(goleak.IgnoreTopFunction("github.com/onsi/ginkgo/internal/specrunner.(*SpecRunner).registerForInterrupts"))).NotTo(HaveOccurred())
	})
})
