package pipelinepoc_test

import (
	"github.com/SamYuan1990/pipelinepoc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pipeline", func() {
	It("Init", func() {
		p := &pipelinepoc.Pipeline{}
		err := p.Init()
		Expect(err).NotTo(HaveOccurred())
		Expect(p).NotTo(BeNil())
		Expect(p.FIFO).NotTo(BeNil())
		Expect(p.Current).NotTo(BeNil())
		Expect(p.Current.Len()).Should(Equal(0))
	})

	It("Add item", func() {
		p := &pipelinepoc.Pipeline{}
		err := p.Init()
		Expect(err).NotTo(HaveOccurred())
		node := ConstructNodeWithSingleKey("key", 0)
		p.Add(node.GetKeys(), node)
		Expect(p.Current.Len()).Should(Equal(1))
		_, value, ok := p.Current.RemoveOldest()
		Expect(ok).Should(Equal(true))
		Expect(value.(*pipelinepoc.Node)).Should(Equal(node))
	})

	It("Add item with new lru", func() {
		p := &pipelinepoc.Pipeline{}
		err := p.Init()
		Expect(err).NotTo(HaveOccurred())
		node := ConstructNodeWithSingleKey("key", 0)
		p.AddWithNewLRU(node.GetKeys(), node)
		Expect(p.FIFO.GetLen()).Should(Equal(2))
		Expect(p.Current.Len()).Should(Equal(1))
		_, value, ok := p.Current.RemoveOldest()
		Expect(ok).Should(Equal(true))
		Expect(value.(*pipelinepoc.Node)).Should(Equal(node))
		p.FIFO.Dequeue()
		v, err := p.FIFO.Dequeue()
		Expect(err).NotTo(HaveOccurred())
		Expect(v).Should(Equal(p.Current))
	})

	It("RemoveOldest Item", func() {
		p := &pipelinepoc.Pipeline{}
		err := p.Init()
		Expect(err).NotTo(HaveOccurred())
		node := ConstructNodeWithSingleKey("key", 0)
		p.Add(node.GetKeys(), node)
		_, value, ok := p.RemoveOldest()
		Expect(ok).Should(Equal(true))
		Expect(value.(*pipelinepoc.Node)).Should(Equal(node))
		Expect(p.Current.Len()).Should(Equal(0))
		Expect(p.FIFO.GetLen()).Should(Equal(1))
	})

	It("RemoveOldest empty Item", func() {
		p := &pipelinepoc.Pipeline{}
		err := p.Init()
		Expect(err).NotTo(HaveOccurred())
		_, _, ok := p.RemoveOldest()
		Expect(ok).Should(Equal(false))
		Expect(p.Current.Len()).Should(Equal(0))
		Expect(p.FIFO.GetLen()).Should(Equal(1))
	})

	It("RemoveOldest and Remove LRU", func() {
		p := &pipelinepoc.Pipeline{}
		err := p.Init()
		Expect(err).NotTo(HaveOccurred())
		node := ConstructNodeWithSingleKey("key", 0)
		p.AddWithNewLRU(node.GetKeys(), node)
		_, value, ok := p.RemoveOldest()
		Expect(ok).Should(Equal(true))
		Expect(value.(*pipelinepoc.Node)).Should(Equal(node))
		Expect(p.Current.Len()).Should(Equal(0))
		Expect(p.FIFO.GetLen()).Should(Equal(1))
	})
})
