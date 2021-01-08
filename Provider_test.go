package pipelinepoc_test

import (
	"strings"
	"time"

	"github.com/SamYuan1990/pipelinepoc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Provider", func() {
	Context("ConvertTxToNode", func() {
		It("Basic testing", func() {
			ProviderImpl := &pipelinepoc.ProviderImpl{}
			Wkeys := make([]pipelinepoc.Key, 0)
			Rkeys := make([]pipelinepoc.Key, 0)
			tximpl := &pipelinepoc.TxImpl{
				Wkeys: Wkeys,
				Rkeys: Rkeys,
			}
			node := ProviderImpl.ConvertTxToNode(tximpl)
			Expect(node.Tx.Wkeys).Should(Equal(Wkeys))
			Expect(node.Tx.Rkeys).Should(Equal(Rkeys))
		})
	})

	Context("Convert", func() {
		It("Basic testing", func() {
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
			p.Comming.Enqueue(BlockImpl)
			time.Sleep(100 * time.Millisecond)
			node := ProviderImpl.ConvertTxToNode(tximpl)
			Expect(len(p.PCurrent.Keys())).Should(Equal(1))
			key, value, ok := p.PCurrent.RemoveOldest()
			Expect(ok).Should(Equal(true))
			Expect(key).Should(Equal(node.GetKeys()))
			Expect(value.(*pipelinepoc.Node).Tx.Wkeys).Should(Equal(Wkeys))
			Expect(value.(*pipelinepoc.Node).Tx.Rkeys).Should(Equal(Rkeys))
			Expect(len(p.PCurrent.Keys())).Should(Equal(0))
			Expect(p.Comming.GetLen()).Should(Equal(0))
			Expect(p.Output.GetLen()).Should(Equal(1))
			v, _ := p.Output.Dequeue()
			Expect(v).Should(Equal(BlockImpl))
			ProviderImpl.Stop()
			//Expect(goleak.Find(goleak.IgnoreTopFunction("github.com/onsi/ginkgo/internal/specrunner.(*SpecRunner).registerForInterrupts"))).NotTo(HaveOccurred())
		})

		It("Key Merge case 1", func() {
			p := &pipelinepoc.Pipeline{}
			p.Init()
			ProviderImpl := &pipelinepoc.ProviderImpl{
				Pipeline: p,
			}
			go ProviderImpl.Convert()
			BlockImpl := ConstructBlock("key", 0)
			AnotherBlock := ConstructBlock("key", 1)
			p.Comming.Enqueue(BlockImpl)
			p.Comming.Enqueue(AnotherBlock)
			time.Sleep(100 * time.Millisecond)

			Expect(len(p.PCurrent.Keys())).Should(Equal(1))
			key, value, ok := p.PCurrent.RemoveOldest()
			Expect(ok).Should(Equal(true))
			// TODO as key merge returns keykey but not key
			Expect(key).Should(Equal("keykey"))
			node := value.(*pipelinepoc.Node)
			Expect(node.Next).NotTo(BeNil())
			Expect(len(p.PCurrent.Keys())).Should(Equal(0))
			Expect(p.Comming.GetLen()).Should(Equal(0))
			Expect(p.Output.GetLen()).Should(Equal(2))
			v, _ := p.Output.Dequeue()
			Expect(v).Should(Equal(BlockImpl))
			v, _ = p.Output.Dequeue()
			Expect(v).Should(Equal(AnotherBlock))
			ProviderImpl.Stop()
		})

		It("Key Merge case 2", func() {
			p := &pipelinepoc.Pipeline{}
			p.Init()
			ProviderImpl := &pipelinepoc.ProviderImpl{
				Pipeline: p,
			}
			go ProviderImpl.Convert()
			BlockImpl := ConstructBlocks("key", "abc")
			AnotherBlock := ConstructBlock("key", 1)
			p.Comming.Enqueue(BlockImpl)
			p.Comming.Enqueue(AnotherBlock)
			time.Sleep(100 * time.Millisecond)
			ProviderImpl.Stop()

			Expect(len(p.PCurrent.Keys())).Should(Equal(1))
			key, value, ok := p.PCurrent.RemoveOldest()
			Expect(ok).Should(Equal(true))
			Expect(key).Should(Equal("keyabckey"))

			node := value.(*pipelinepoc.Node)
			Expect(node.Next).NotTo(BeNil())
			Expect(strings.Contains(node.GetKeys(), "key")).Should(Equal(true))
			Expect(strings.Contains(node.GetKeys(), "abc")).Should(Equal(true))
			Expect(len(p.PCurrent.Keys())).Should(Equal(0))
			Expect(p.Output.GetLen()).Should(Equal(2))
			Expect(p.Comming.GetLen()).Should(Equal(0))
			v, _ := p.Output.Dequeue()
			Expect(v).Should(Equal(BlockImpl))
			v, _ = p.Output.Dequeue()
			Expect(v).Should(Equal(AnotherBlock))
		})

		It("Key Merge case 3", func() {
			p := &pipelinepoc.Pipeline{}
			p.Init()
			ProviderImpl := &pipelinepoc.ProviderImpl{
				Pipeline: p,
			}
			go ProviderImpl.Convert()
			BlockImpl := ConstructBlock("key", 1)
			AnotherBlock := ConstructBlocks("key", "abc")
			p.Comming.Enqueue(BlockImpl)
			p.Comming.Enqueue(AnotherBlock)
			time.Sleep(100 * time.Millisecond)
			ProviderImpl.Stop()

			Expect(len(p.PCurrent.Keys())).Should(Equal(1))
			key, value, ok := p.PCurrent.RemoveOldest()
			Expect(ok).Should(Equal(true))
			Expect(key).Should(Equal("keykeyabc"))

			node := value.(*pipelinepoc.Node)
			Expect(node.Next).NotTo(BeNil())
			Expect(len(p.PCurrent.Keys())).Should(Equal(0))
			Expect(p.Output.GetLen()).Should(Equal(2))
			Expect(p.Comming.GetLen()).Should(Equal(0))
			v, _ := p.Output.Dequeue()
			Expect(v).Should(Equal(BlockImpl))
			v, _ = p.Output.Dequeue()
			Expect(v).Should(Equal(AnotherBlock))
		})

		It("Key Merge case 4", func() {
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
			time.Sleep(100 * time.Millisecond)

			Expect(len(p.PCurrent.Keys())).Should(Equal(1))
			Expect(len(p.CCurrent.Keys())).Should(Equal(2))

			Expect(p.Output.GetLen()).Should(Equal(3))
			Expect(p.Comming.GetLen()).Should(Equal(0))
			v, _ := p.Output.Dequeue()
			Expect(v).Should(Equal(BlockImpl))
			v, _ = p.Output.Dequeue()
			Expect(v).Should(Equal(AnotherBlock))
			v, _ = p.Output.Dequeue()
			Expect(v).Should(Equal(ThridBlock))

			key, value, ok := p.PCurrent.RemoveOldest()
			Expect(ok).Should(Equal(true))
			node := value.(*pipelinepoc.Node)
			Expect(node.Next).To(BeNil())
			Expect(node.GetKeys()).Should(Equal("keyabc"))
			Expect(key.(string)).Should(Equal("keyabc"))
			key, value, ok = p.CCurrent.RemoveOldest()
			Expect(ok).Should(Equal(true))
			node = value.(*pipelinepoc.Node)
			Expect(node.Next).To(BeNil())
			Expect(node.GetKeys()).Should(Equal("key"))
			Expect(key.(string)).Should(Equal("key"))
			key, value, ok = p.CCurrent.RemoveOldest()
			Expect(ok).Should(Equal(true))
			node = value.(*pipelinepoc.Node)
			Expect(node.Next).To(BeNil())
			Expect(node.GetKeys()).Should(Equal("abc"))
			Expect(key.(string)).Should(Equal("abc"))
		})

		It("Key Merge case 5", func() {
			p := &pipelinepoc.Pipeline{}
			p.Init()
			ProviderImpl := &pipelinepoc.ProviderImpl{
				Pipeline: p,
			}
			go ProviderImpl.Convert()
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
			time.Sleep(100 * time.Millisecond)

			Expect(len(p.PCurrent.Keys())).Should(Equal(3))
			Expect(len(p.CCurrent.Keys())).Should(Equal(2))

			Expect(p.Output.GetLen()).Should(Equal(5))
			Expect(p.Comming.GetLen()).Should(Equal(0))
			v, _ := p.Output.Dequeue()
			Expect(v).Should(Equal(BlockImpl))
			v, _ = p.Output.Dequeue()
			Expect(v).Should(Equal(AnotherBlock))
			v, _ = p.Output.Dequeue()
			Expect(v).Should(Equal(ThridBlock))
			v, _ = p.Output.Dequeue()
			Expect(v).Should(Equal(FourthBlock))
			v, _ = p.Output.Dequeue()
			Expect(v).Should(Equal(FiFthBlock))

			key, value, ok := p.PCurrent.RemoveOldest()
			Expect(ok).Should(Equal(true))
			node := value.(*pipelinepoc.Node)
			Expect(node.Next).To(BeNil())
			Expect(node.GetKeys()).Should(Equal("keyabc"))
			Expect(key.(string)).Should(Equal("keyabc"))

			key, value, ok = p.PCurrent.RemoveOldest()
			Expect(ok).Should(Equal(true))
			node = value.(*pipelinepoc.Node)
			Expect(node.Next).To(BeNil())
			Expect(node.GetKeys()).Should(Equal("xyz"))
			Expect(key.(string)).Should(Equal("xyz"))

			key, value, ok = p.PCurrent.RemoveOldest()
			Expect(ok).Should(Equal(true))
			node = value.(*pipelinepoc.Node)
			Expect(node.Next).To(BeNil())
			Expect(node.GetKeys()).Should(Equal("edf"))
			Expect(key.(string)).Should(Equal("edf"))

			key, value, ok = p.CCurrent.RemoveOldest()
			Expect(ok).Should(Equal(true))
			node = value.(*pipelinepoc.Node)
			Expect(node.Next).To(BeNil())
			Expect(node.GetKeys()).Should(Equal("key"))
			Expect(key.(string)).Should(Equal("key"))
			key, value, ok = p.CCurrent.RemoveOldest()
			Expect(ok).Should(Equal(true))
			node = value.(*pipelinepoc.Node)
			Expect(node.Next).To(BeNil())
			Expect(node.GetKeys()).Should(Equal("abc"))
			Expect(key.(string)).Should(Equal("abc"))
			p.SwitchC()
			time.Sleep(100 * time.Millisecond)
			Expect(p.Output.GetLen()).Should(Equal(1))
			v, _ = p.Output.Dequeue()
			Expect(v).Should(Equal(SixthBlock))
			key, value, ok = p.PCurrent.RemoveOldest()
			Expect(ok).Should(Equal(true))
			node = value.(*pipelinepoc.Node)
			Expect(node.Next).To(BeNil())
			Expect(node.GetKeys()).Should(Equal("xyzedf"))
			Expect(key.(string)).Should(Equal("xyzedf"))
		})
	})

})
