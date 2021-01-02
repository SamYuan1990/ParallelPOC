package pipelinepoc_test

import (
	"time"

	"github.com/SamYuan1990/pipelinepoc"
	lru "github.com/hashicorp/golang-lru"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/goleak"
)

var _ = Describe("Consumer", func() {
	It("should stopped", func() {
		LRU, _ := lru.New(128)
		node := ConstructNodeWithSingleKey("key", 0)
		LRU.Add(node.GetKeys(), node)
		output := make(chan *pipelinepoc.Node, 10)
		defer close(output)
		c := &pipelinepoc.Consumer{
			LRU: LRU,
		}
		go c.Consume(output)
		time.Sleep(time.Millisecond)
		<-output
		c.Stop()
		Expect(len(LRU.Keys())).Should(Equal(0))
		Expect(goleak.Find(goleak.IgnoreTopFunction("github.com/onsi/ginkgo/internal/specrunner.(*SpecRunner).registerForInterrupts"))).NotTo(HaveOccurred())
	})

	It("should stopped", func() {
		LRU, _ := lru.New(128)
		output := make(chan *pipelinepoc.Node, 10)
		defer close(output)

		c := &pipelinepoc.Consumer{
			LRU: LRU,
		}
		go c.Consume(output)
		c1 := &pipelinepoc.Consumer{
			LRU: LRU,
		}
		go c1.Consume(output)
		time.Sleep(time.Millisecond)
		for i := 0; i < 10; i++ {
			current := time.Now()
			node := ConstructNodeWithSingleKey(current.String(), 0)
			LRU.Add(node.GetKeys(), node)
		}
		var n int
		for n < 10 {
			<-output
			n++
		}
		c1.Stop()
		c.Stop()
		Expect(len(LRU.Keys())).Should(Equal(0))
		Expect(goleak.Find(goleak.IgnoreTopFunction("github.com/onsi/ginkgo/internal/specrunner.(*SpecRunner).registerForInterrupts"))).NotTo(HaveOccurred())
	})

})
