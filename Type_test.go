package pipelinepoc_test

import (
	"github.com/SamYuan1990/pipelinepoc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Type", func() {

	Context("Node", func() {
		It("add node", func() {
			n1 := ConstructNodeWithSingleKey("key1", 0)
			n2 := ConstructNodeWithSingleKey("key1", 0)
			n3 := ConstructNodeWithSingleKey("key1", 0)
			n1.Add(n2)
			n1.Add(n3)
			Expect(n2).Should(Equal(n1.Next))
			Expect(n3).Should(Equal(n2.Next))
		})
	})

	Context("type", func() {
		It("New node", func() {
			data := &pipelinepoc.Node{}
			Expect(data.Next).Should(BeNil())
		})

		It("Next", func() {
			Next := &pipelinepoc.Node{}
			data := &pipelinepoc.Node{
				Next: Next,
			}
			Expect(data.Next).Should(Equal(Next))
		})

		It("Wkeys", func() {
			Wkeys := make([]pipelinepoc.Key, 0)
			tximpl := pipelinepoc.TxImpl{
				Wkeys: Wkeys,
			}
			data := &pipelinepoc.Node{
				Tx: &tximpl,
			}
			Expect(data.Tx.Wkeys).Should(Equal(Wkeys))
		})
	})

	Context("Block Sample", func() {
		It("BlockImp", func() {
			txs := make([]*pipelinepoc.TxImpl, 0)
			BlockImpl := &pipelinepoc.BlockImpl{
				Txs: txs,
			}
			Expect(BlockImpl.GetTxs()).Should(Equal(txs))
		})

		It("TxImpl", func() {
			Wkeys := make([]pipelinepoc.Key, 0)
			Rkeys := make([]pipelinepoc.Key, 0)

			tximpl := &pipelinepoc.TxImpl{
				Wkeys: Wkeys,
				Rkeys: Rkeys,
			}
			Expect(tximpl.GetWkeys()).Should(Equal(Wkeys))
			Expect(tximpl.GetRkeys()).Should(Equal(Rkeys))
		})
	})
})
