package pipelinepoc_test

import (
	"strings"

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
			Expect(node.Wkeys).Should(Equal(Wkeys))
			Expect(node.Rkeys).Should(Equal(Rkeys))
		})
	})

	Context("Convert", func() {
		It("Basic testing", func() {
			p := &pipelinepoc.Pipeline{}
			p.Init()
			ProviderImpl := &pipelinepoc.ProviderImpl{
				Pipeline: p,
			}
			Wkeys := make([]pipelinepoc.Key, 0)
			Rkeys := make([]pipelinepoc.Key, 0)

			tximpl := pipelinepoc.TxImpl{
				Wkeys: Wkeys,
				Rkeys: Rkeys,
			}
			txs := make([]pipelinepoc.TxImpl, 0)
			txs = append(txs, tximpl)

			BlockImpl := &pipelinepoc.BlockImpl{
				Txs: txs,
			}
			ProviderImpl.Convert(BlockImpl)
			node := ProviderImpl.ConvertTxToNode(&tximpl)
			Expect(len(p.Current.Keys())).Should(Equal(1))
			key, value, ok := p.RemoveOldest()
			Expect(ok).Should(Equal(true))
			Expect(key).Should(Equal(node.GetKeys()))
			Expect(value.(*pipelinepoc.Node).Wkeys).Should(Equal(Wkeys))
			Expect(value.(*pipelinepoc.Node).Rkeys).Should(Equal(Rkeys))
			Expect(len(p.Current.Keys())).Should(Equal(0))
		})

		It("Key Merge", func() {
			p := &pipelinepoc.Pipeline{}
			p.Init()
			ProviderImpl := &pipelinepoc.ProviderImpl{
				Pipeline: p,
			}
			BlockImpl := ConstructBlock("key", 0)
			AnotherBlock := ConstructBlock("key", 1)
			ProviderImpl.Convert(BlockImpl)
			ProviderImpl.Convert(AnotherBlock)
			Expect(len(p.Current.Keys())).Should(Equal(1))
			_, value, ok := p.RemoveOldest()
			Expect(ok).Should(Equal(true))
			node := value.(*pipelinepoc.Node)
			Expect(node.Next).NotTo(BeNil())
			Expect(len(p.Current.Keys())).Should(Equal(0))
		})

		It("Key Merge 2", func() {
			p := &pipelinepoc.Pipeline{}
			p.Init()
			ProviderImpl := &pipelinepoc.ProviderImpl{
				Pipeline: p,
			}
			BlockImpl := ConstructBlocks("key", "abc")
			AnotherBlock := ConstructBlock("key", 1)
			ProviderImpl.Convert(BlockImpl)
			ProviderImpl.Convert(AnotherBlock)
			Expect(len(p.Current.Keys())).Should(Equal(1))
			_, value, ok := p.RemoveOldest()
			Expect(ok).Should(Equal(true))
			node := value.(*pipelinepoc.Node)
			Expect(node.Next).NotTo(BeNil())
			Expect(strings.Contains(node.GetKeys(), "key")).Should(Equal(true))
			Expect(strings.Contains(node.GetKeys(), "abc")).Should(Equal(true))
			Expect(len(p.Current.Keys())).Should(Equal(0))
		})

		It("Key Merge 2", func() {
			p := &pipelinepoc.Pipeline{}
			p.Init()
			ProviderImpl := &pipelinepoc.ProviderImpl{
				Pipeline: p,
			}
			BlockImpl := ConstructBlock("key", 1)
			AnotherBlock := ConstructBlocks("key", "abc")
			ProviderImpl.Convert(BlockImpl)
			ProviderImpl.Convert(AnotherBlock)
			Expect(len(p.Current.Keys())).Should(Equal(1))
			key, value, ok := p.RemoveOldest()
			Expect(ok).Should(Equal(true))
			node := value.(*pipelinepoc.Node)
			Expect(node.Next).NotTo(BeNil())
			Expect(strings.Contains(key.(string), "key")).Should(Equal(true))
			Expect(strings.Contains(key.(string), "abc")).Should(Equal(true))
			Expect(len(p.Current.Keys())).Should(Equal(0))
		})

		It("Key Merge 3", func() {
			p := &pipelinepoc.Pipeline{}
			p.Init()
			ProviderImpl := &pipelinepoc.ProviderImpl{
				Pipeline: p,
			}
			BlockImpl := ConstructBlock("key", 0)
			AnotherBlock := ConstructBlock("abc", 1)
			ThridBlock := ConstructBlocks("key", "abc")
			ProviderImpl.Convert(BlockImpl)
			ProviderImpl.Convert(AnotherBlock)
			ProviderImpl.Convert(ThridBlock)
			Expect(len(p.Current.Keys())).Should(Equal(3))
			key, value, ok := p.RemoveOldest()
			Expect(ok).Should(Equal(true))
			node := value.(*pipelinepoc.Node)
			Expect(node.Next).To(BeNil())
			Expect(node.GetKeys()).Should(Equal("key"))
			Expect(key.(string)).Should(Equal("key"))
			key, value, ok = p.RemoveOldest()
			Expect(ok).Should(Equal(true))
			node = value.(*pipelinepoc.Node)
			Expect(node.Next).To(BeNil())
			Expect(node.GetKeys()).Should(Equal("abc"))
			Expect(key.(string)).Should(Equal("abc"))
			key, value, ok = p.RemoveOldest()
			Expect(ok).Should(Equal(true))
			node = value.(*pipelinepoc.Node)
			Expect(node.Next).To(BeNil())
			Expect(node.GetKeys()).Should(Equal("keyabc"))
			Expect(key.(string)).Should(Equal("keyabc"))
		})
	})

})
