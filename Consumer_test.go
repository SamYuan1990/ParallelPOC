package pipelinepoc_test

import (
	"sync"
	"time"

	"github.com/SamYuan1990/pipelinepoc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Consumer", func() {
	It("Basic testing", func() {
		pipelinepoc.GetLogger().Info("Consumer basic testing")
		p := &pipelinepoc.Pipeline{}
		p.Init()
		ProviderImpl := &pipelinepoc.ProviderImpl{
			Pipeline: p,
		}
		go ProviderImpl.Convert()
		Wkeys := make([]pipelinepoc.Key, 0)
		Rkeys := make([]pipelinepoc.Key, 0)
		tximpl := &pipelinepoc.TxImpl{
			Wkeys: Wkeys,
			Rkeys: Rkeys,
		}
		txs := make([]*pipelinepoc.TxImpl, 0)
		txs = append(txs, tximpl)
		BlockImpl := &pipelinepoc.BlockImpl{
			Txs: txs,
		}
		wg := &sync.WaitGroup{}
		c := &pipelinepoc.Consumer{
			Pipeline: p,
			Wg:       wg,
		}
		go c.Consume()
		p.Comming.Enqueue(BlockImpl)
		time.Sleep(500 * time.Millisecond)
		Expect(len(p.PCurrent.Keys())).Should(Equal(0))
		ProviderImpl.Stop()
		c.Stop()
		v, _ := p.Output.Dequeue()
		Expect(v).Should(Equal(BlockImpl))
		Expect(v.(*pipelinepoc.BlockImpl).Txs[0].Processed).Should(BeTrue())
	})
	It("Key Merge case 1", func() {
		pipelinepoc.GetLogger().Info("start Key Merge case 1")
		p := &pipelinepoc.Pipeline{}
		p.Init()
		ProviderImpl := &pipelinepoc.ProviderImpl{
			Pipeline: p,
		}
		go ProviderImpl.Convert()
		wg := &sync.WaitGroup{}
		c := &pipelinepoc.Consumer{
			Pipeline: p,
			Wg:       wg,
		}
		go c.Consume()
		BlockImpl := ConstructBlock("key", 0)
		AnotherBlock := ConstructBlock("key", 1)
		p.Comming.Enqueue(BlockImpl)
		p.Comming.Enqueue(AnotherBlock)
		time.Sleep(500 * time.Millisecond)

		ProviderImpl.Stop()
		c.Stop()
		v, _ := p.Output.Dequeue()
		Expect(v).Should(Equal(BlockImpl))
		Expect(v.(*pipelinepoc.BlockImpl).Txs[0].Processed).Should(BeTrue())
		v, _ = p.Output.Dequeue()
		Expect(v).Should(Equal(AnotherBlock))
		Expect(v.(*pipelinepoc.BlockImpl).Txs[0].Processed).Should(BeTrue())
	})

	It("Key Merge case 2", func() {
		pipelinepoc.GetLogger().Info("start Key Merge case 2")
		p := &pipelinepoc.Pipeline{}
		p.Init()
		ProviderImpl := &pipelinepoc.ProviderImpl{
			Pipeline: p,
		}
		go ProviderImpl.Convert()
		wg := &sync.WaitGroup{}
		c := &pipelinepoc.Consumer{
			Pipeline: p,
			Wg:       wg,
		}
		go c.Consume()
		BlockImpl := ConstructBlocks("key", "abc")
		AnotherBlock := ConstructBlock("key", 1)
		p.Comming.Enqueue(BlockImpl)
		p.Comming.Enqueue(AnotherBlock)
		time.Sleep(500 * time.Millisecond)

		ProviderImpl.Stop()
		c.Stop()

		v, _ := p.Output.Dequeue()
		Expect(v).Should(Equal(BlockImpl))
		Expect(v.(*pipelinepoc.BlockImpl).Txs[0].Processed).Should(BeTrue())
		v, _ = p.Output.Dequeue()
		Expect(v).Should(Equal(AnotherBlock))
		Expect(v.(*pipelinepoc.BlockImpl).Txs[0].Processed).Should(BeTrue())
	})

	It("Key Merge case 3", func() {
		pipelinepoc.GetLogger().Info("start Key Merge case 3")
		p := &pipelinepoc.Pipeline{}
		p.Init()
		ProviderImpl := &pipelinepoc.ProviderImpl{
			Pipeline: p,
		}
		go ProviderImpl.Convert()
		wg := &sync.WaitGroup{}
		c := &pipelinepoc.Consumer{
			Pipeline: p,
			Wg:       wg,
		}
		BlockImpl := ConstructBlock("key", 1)
		AnotherBlock := ConstructBlocks("key", "abc")
		p.Comming.Enqueue(BlockImpl)
		p.Comming.Enqueue(AnotherBlock)
		go c.Consume()
		time.Sleep(500 * time.Millisecond)

		ProviderImpl.Stop()
		c.Stop()

		v, _ := p.Output.Dequeue()
		Expect(v).Should(Equal(BlockImpl))
		Expect(v.(*pipelinepoc.BlockImpl).Txs[0].Processed).Should(BeTrue())
		v, _ = p.Output.Dequeue()
		Expect(v).Should(Equal(AnotherBlock))
		Expect(v.(*pipelinepoc.BlockImpl).Txs[0].Processed).Should(BeTrue())
	})

	It("Key Merge case 4", func() {
		pipelinepoc.GetLogger().Info("start Key Merge case 4")
		p := &pipelinepoc.Pipeline{}
		p.Init()
		ProviderImpl := &pipelinepoc.ProviderImpl{
			Pipeline: p,
		}
		go ProviderImpl.Convert()
		BlockImpl := ConstructBlock("key", 0)
		AnotherBlock := ConstructBlock("abc", 1)
		ThridBlock := ConstructBlocks("key", "abc")
		p.Comming.Enqueue(BlockImpl)
		p.Comming.Enqueue(AnotherBlock)
		p.Comming.Enqueue(ThridBlock)
		time.Sleep(500 * time.Millisecond)

		ProviderImpl.Stop()

		wg := &sync.WaitGroup{}
		Switcher := &pipelinepoc.Switcher{
			Pipeline: p,
			Wg:       wg,
			Count:    1,
		}
		Switcher.Init()
		go Switcher.Switch()
		c := &pipelinepoc.Consumer{
			Pipeline: p,
			Wg:       wg,
		}
		Switcher.AddConsumer(c)
		go c.Consume()
		time.Sleep(1000 * time.Millisecond)
		c.Stop()
		Switcher.Stop()

		v, _ := p.Output.Dequeue()
		Expect(v).Should(Equal(BlockImpl))
		Expect(v.(*pipelinepoc.BlockImpl).Txs[0].Processed).Should(BeTrue())
		v, _ = p.Output.Dequeue()
		Expect(v).Should(Equal(AnotherBlock))
		Expect(v.(*pipelinepoc.BlockImpl).Txs[0].Processed).Should(BeTrue())
		v, _ = p.Output.Dequeue()
		Expect(v).Should(Equal(ThridBlock))
		Expect(v.(*pipelinepoc.BlockImpl).Txs[0].Processed).Should(BeTrue())
	})

	It("Key Merge case 4 parallel 2", func() {
		pipelinepoc.GetLogger().Info("start Key Merge case 4 parallel 2")
		p := &pipelinepoc.Pipeline{}
		p.Init()
		ProviderImpl := &pipelinepoc.ProviderImpl{
			Pipeline: p,
		}
		go ProviderImpl.Convert()
		BlockImpl := ConstructBlock("key", 0)
		AnotherBlock := ConstructBlock("abc", 1)
		ThridBlock := ConstructBlocks("key", "abc")
		p.Comming.Enqueue(BlockImpl)
		p.Comming.Enqueue(AnotherBlock)
		p.Comming.Enqueue(ThridBlock)
		time.Sleep(500 * time.Millisecond)

		ProviderImpl.Stop()

		wg := &sync.WaitGroup{}
		c := &pipelinepoc.Consumer{
			Pipeline: p,
			Wg:       wg,
		}
		c1 := &pipelinepoc.Consumer{
			Pipeline: p,
			Wg:       wg,
		}
		Switcher := &pipelinepoc.Switcher{
			Pipeline: p,
			Wg:       wg,
			Count:    2,
		}
		Switcher.Init()
		Switcher.AddConsumer(c)
		Switcher.AddConsumer(c1)
		go Switcher.Switch()
		go c.Consume()
		go c1.Consume()
		time.Sleep(500 * time.Millisecond)

		c.Stop()
		c1.Stop()
		Switcher.Stop()

		v, _ := p.Output.Dequeue()
		Expect(v).Should(Equal(BlockImpl))
		Expect(v.(*pipelinepoc.BlockImpl).Txs[0].Processed).Should(BeTrue())
		v, _ = p.Output.Dequeue()
		Expect(v).Should(Equal(AnotherBlock))
		Expect(v.(*pipelinepoc.BlockImpl).Txs[0].Processed).Should(BeTrue())
		v, _ = p.Output.Dequeue()
		Expect(v).Should(Equal(ThridBlock))
		Expect(v.(*pipelinepoc.BlockImpl).Txs[0].Processed).Should(BeTrue())
	})

	It("Key Merge case 5 parallel 2", func() {
		pipelinepoc.GetLogger().Info("start Key Merge case 5 parallel 2")
		p := &pipelinepoc.Pipeline{}
		p.Init()
		ProviderImpl := &pipelinepoc.ProviderImpl{
			Pipeline: p,
		}
		go ProviderImpl.Convert()
		wg := &sync.WaitGroup{}
		Switcher := &pipelinepoc.Switcher{
			Pipeline: p,
			Wg:       wg,
			Count:    2,
		}
		Switcher.Init()
		c := &pipelinepoc.Consumer{
			Pipeline: p,
			Wg:       wg,
		}
		c1 := &pipelinepoc.Consumer{
			Pipeline: p,
			Wg:       wg,
		}
		Switcher.AddConsumer(c)
		Switcher.AddConsumer(c1)
		go Switcher.Switch()

		BlockImpl := ConstructBlock("key", 0)
		AnotherBlock := ConstructBlock("abc", 1)
		ThridBlock := ConstructBlocks("key", "abc")
		FourthBlock := ConstructBlock("xyz", 0)
		FiFthBlock := ConstructBlock("edf", 0)
		SixthBlock := ConstructBlocks("xyz", "edf")

		p.Comming.Enqueue(BlockImpl)
		p.Comming.Enqueue(AnotherBlock)
		p.Comming.Enqueue(ThridBlock)
		p.Comming.Enqueue(FourthBlock)
		p.Comming.Enqueue(FiFthBlock)
		p.Comming.Enqueue(SixthBlock)
		go c.Consume()
		go c1.Consume()
		time.Sleep(500 * time.Millisecond)

		ProviderImpl.Stop()
		c.Stop()
		c1.Stop()
		Switcher.Stop()

		for v, err := p.Output.Dequeue(); err != nil && v != nil; {
			Expect(v).Should(Equal(BlockImpl))
			Expect(v.(*pipelinepoc.BlockImpl).Txs[0].Processed).Should(BeTrue())
		}
		for v, err := p.Output.Dequeue(); err != nil && v != nil; {
			Expect(v).Should(Equal(AnotherBlock))
			Expect(v.(*pipelinepoc.BlockImpl).Txs[0].Processed).Should(BeTrue())
		}
		for v, err := p.Output.Dequeue(); err != nil && v != nil; {
			Expect(v).Should(Equal(ThridBlock))
			Expect(v.(*pipelinepoc.BlockImpl).Txs[0].Processed).Should(BeTrue())
		}
		for v, err := p.Output.Dequeue(); err != nil && v != nil; {
			Expect(v).Should(Equal(FourthBlock))
			Expect(v.(*pipelinepoc.BlockImpl).Txs[0].Processed).Should(BeTrue())
		}
		for v, err := p.Output.Dequeue(); err != nil && v != nil; {
			Expect(v).Should(Equal(FiFthBlock))
			Expect(v.(*pipelinepoc.BlockImpl).Txs[0].Processed).Should(BeTrue())
		}
		for v, err := p.Output.Dequeue(); err != nil && v != nil; {
			Expect(v).Should(Equal(SixthBlock))
			Expect(v.(*pipelinepoc.BlockImpl).Txs[0].Processed).Should(BeTrue())
		}
	})

	It("Key Merge case 5", func() {
		pipelinepoc.GetLogger().Info("start Key Merge case 5")
		p := &pipelinepoc.Pipeline{}
		p.Init()
		ProviderImpl := &pipelinepoc.ProviderImpl{
			Pipeline: p,
		}
		go ProviderImpl.Convert()
		wg := &sync.WaitGroup{}
		Switcher := &pipelinepoc.Switcher{
			Pipeline: p,
			Wg:       wg,
			Count:    1,
		}
		Switcher.Init()
		c := &pipelinepoc.Consumer{
			Pipeline: p,
			Wg:       wg,
		}
		Switcher.AddConsumer(c)
		go Switcher.Switch()
		BlockImpl := ConstructBlock("key", 0)
		AnotherBlock := ConstructBlock("abc", 1)
		ThridBlock := ConstructBlocks("key", "abc")
		FourthBlock := ConstructBlock("xyz", 0)
		FiFthBlock := ConstructBlock("edf", 0)
		SixthBlock := ConstructBlocks("xyz", "edf")

		p.Comming.Enqueue(BlockImpl)
		p.Comming.Enqueue(AnotherBlock)
		p.Comming.Enqueue(ThridBlock)
		p.Comming.Enqueue(FourthBlock)
		p.Comming.Enqueue(FiFthBlock)
		p.Comming.Enqueue(SixthBlock)
		go c.Consume()
		time.Sleep(500 * time.Millisecond)

		defer ProviderImpl.Stop()
		defer c.Stop()
		defer Switcher.Stop()

		for v, err := p.Output.Dequeue(); err != nil; {
			Expect(v).Should(Equal(BlockImpl))
			Expect(v.(*pipelinepoc.BlockImpl).Txs[0].Processed).Should(BeTrue())
		}
		for v, err := p.Output.Dequeue(); err != nil; {
			Expect(v).Should(Equal(AnotherBlock))
			Expect(v.(*pipelinepoc.BlockImpl).Txs[0].Processed).Should(BeTrue())
		}
		for v, err := p.Output.Dequeue(); err != nil; {
			Expect(v).Should(Equal(ThridBlock))
			Expect(v.(*pipelinepoc.BlockImpl).Txs[0].Processed).Should(BeTrue())
		}
		for v, err := p.Output.Dequeue(); err != nil; {
			Expect(v).Should(Equal(FourthBlock))
			Expect(v.(*pipelinepoc.BlockImpl).Txs[0].Processed).Should(BeTrue())
		}
		for v, err := p.Output.Dequeue(); err != nil; {
			Expect(v).Should(Equal(FiFthBlock))
			Expect(v.(*pipelinepoc.BlockImpl).Txs[0].Processed).Should(BeTrue())
		}
		for v, err := p.Output.Dequeue(); err != nil; {
			Expect(v).Should(Equal(SixthBlock))
			Expect(v.(*pipelinepoc.BlockImpl).Txs[0].Processed).Should(BeTrue())
		}
	})
})
