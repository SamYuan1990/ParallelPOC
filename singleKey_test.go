package pipelinepoc_test

import (
	"github.com/SamYuan1990/pipelinepoc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SingleKey", func() {
	Context("Single Key", func() {
		It("Create should in a node", func() {
			CreateA := NodeCreation("", -1, "a", 0)
			Tx := make([]*pipelinepoc.Node, 0)
			Tx = append(Tx, CreateA)
			Producer := &pipelinepoc.Producer{
				Nodes: make([]*pipelinepoc.Node, 0),
			}
			Producer.TreeMaking(Tx)
			Expect("a").Should(Equal(Producer.Nodes[0].WKey))
			Expect(0).Should(Equal(Producer.Nodes[0].WKeyVersion))
		})

		It("Read should in a node", func() {
			ReadA := NodeCreation("a", 0, "", -1)
			Tx := make([]*pipelinepoc.Node, 0)
			Tx = append(Tx, ReadA)
			Producer := &pipelinepoc.Producer{
				Nodes: make([]*pipelinepoc.Node, 0),
			}
			Producer.TreeMaking(Tx)
			Expect("a").Should(Equal(Producer.Nodes[0].RKey))
			Expect(0).Should(Equal(Producer.Nodes[0].RKeyVersion))
		})

		It("Update should in a node", func() {
			WriteA := NodeCreation("", 0, "a", 1)
			Tx := make([]*pipelinepoc.Node, 0)
			Tx = append(Tx, WriteA)
			Producer := &pipelinepoc.Producer{
				Nodes: make([]*pipelinepoc.Node, 0),
			}
			Producer.TreeMaking(Tx)
			Expect("a").Should(Equal(Producer.Nodes[0].WKey))
			Expect(1).Should(Equal(Producer.Nodes[0].WKeyVersion))
		})
	})

	Context("Single Key without mvcc", func() {
		// notice the create, read , update can be in different blocks
		// but follow create with as tx0, read as tx1, update as tx2
		// ex block 0(tx0), block 1(tx1), block 2(tx2)
		// if only one tx can be allowed in block
		It("Create, read, update", func() {
			CreateA := NodeCreation("", -1, "a", 0)
			ReadA := NodeCreation("a", 0, "", -1)
			WriteA := NodeCreation("", -1, "a", 1)

			Tx := make([]*pipelinepoc.Node, 0)
			Tx = append(Tx, CreateA)
			Tx = append(Tx, ReadA)
			Tx = append(Tx, WriteA)

			Producer := &pipelinepoc.Producer{
				Nodes: make([]*pipelinepoc.Node, 0),
			}
			Producer.TreeMaking(Tx)
			// creation should be root
			Expect("a").Should(Equal(Producer.Nodes[0].WKey))
			Expect(0).Should(Equal(Producer.Nodes[0].WKeyVersion))
			// then read
			Expect("a").Should(Equal(Producer.Nodes[0].Right.RKey))
			Expect(0).Should(Equal(Producer.Nodes[0].Right.RKeyVersion))
			// then write
			Expect("a").Should(Equal(Producer.Nodes[0].Right.Right.WKey))
			Expect(1).Should(Equal(Producer.Nodes[0].Right.Right.WKeyVersion))
		})

		It("Create, update, read", func() {
			CreateA := NodeCreation("", -1, "a", 0)
			ReadA := NodeCreation("a", 1, "", -1)
			WriteA := NodeCreation("", -1, "a", 1)

			Tx := make([]*pipelinepoc.Node, 0)
			Tx = append(Tx, CreateA)
			Tx = append(Tx, WriteA)
			Tx = append(Tx, ReadA)

			Producer := &pipelinepoc.Producer{
				Nodes: make([]*pipelinepoc.Node, 0),
			}
			Producer.TreeMaking(Tx)
			// creation should be root
			Expect("a").Should(Equal(Producer.Nodes[0].WKey))
			Expect(0).Should(Equal(Producer.Nodes[0].WKeyVersion))
			// then read
			Expect("a").Should(Equal(Producer.Nodes[0].Right.WKey))
			Expect(1).Should(Equal(Producer.Nodes[0].Right.WKeyVersion))
			// then write
			Expect("a").Should(Equal(Producer.Nodes[0].Right.Right.RKey))
			Expect(1).Should(Equal(Producer.Nodes[0].Right.Right.RKeyVersion))
		})

		It("Create, read, read", func() {
			CreateA := NodeCreation("", -1, "a", 0)
			ReadA := NodeCreation("a", 0, "", -1)
			ReadAAnother := NodeCreation("a", 0, "", -1)

			Tx := make([]*pipelinepoc.Node, 0)
			Tx = append(Tx, CreateA)
			Tx = append(Tx, ReadA)
			Tx = append(Tx, ReadAAnother)

			Producer := &pipelinepoc.Producer{
				Nodes: make([]*pipelinepoc.Node, 0),
			}
			Producer.TreeMaking(Tx)
			// creation should be root
			Expect("a").Should(Equal(Producer.Nodes[0].WKey))
			Expect(0).Should(Equal(Producer.Nodes[0].WKeyVersion))
			// then read
			Expect("a").Should(Equal(Producer.Nodes[0].Right.RKey))
			Expect(0).Should(Equal(Producer.Nodes[0].Right.RKeyVersion))
			// then read
			Expect("a").Should(Equal(Producer.Nodes[0].Right.Right.RKey))
			Expect(0).Should(Equal(Producer.Nodes[0].Right.Right.RKeyVersion))
		})

		It("update, update another", func() {
			WriteA := NodeCreation("", -1, "a", 1)
			WriteAAnother := NodeCreation("", -1, "a", 2)

			Tx := make([]*pipelinepoc.Node, 0)
			Tx = append(Tx, WriteA)
			Tx = append(Tx, WriteAAnother)

			Producer := &pipelinepoc.Producer{
				Nodes: make([]*pipelinepoc.Node, 0),
			}
			Producer.TreeMaking(Tx)
			// then write
			Expect("a").Should(Equal(Producer.Nodes[0].WKey))
			Expect(1).Should(Equal(Producer.Nodes[0].WKeyVersion))
			// then write
			Expect("a").Should(Equal(Producer.Nodes[0].Right.WKey))
			Expect(2).Should(Equal(Producer.Nodes[0].Right.WKeyVersion))
		})
	})

	Context("Single Key with mvcc", func() {
		It("Create, Create", func() {
			CreateA := NodeCreation("", -1, "a", 0)
			CreateAAnother := NodeCreation("", -1, "a", 0)
			Tx := make([]*pipelinepoc.Node, 0)
			Tx = append(Tx, CreateA)
			Tx = append(Tx, CreateAAnother)

			Producer := &pipelinepoc.Producer{
				Nodes: make([]*pipelinepoc.Node, 0),
			}
			Producer.TreeMaking(Tx)
			Expect("a").Should(Equal(Producer.Nodes[0].WKey))
			Expect(0).Should(Equal(Producer.Nodes[0].WKeyVersion))
			Expect("a").Should(Equal(Producer.Nodes[0].Left.WKey))
			Expect(0).Should(Equal(Producer.Nodes[0].Left.WKeyVersion))
		})

		It("update, update", func() {
			WriteA := NodeCreation("", -1, "a", 1)
			WriteAAnother := NodeCreation("", -1, "a", 1)
			Tx := make([]*pipelinepoc.Node, 0)
			Tx = append(Tx, WriteA)
			Tx = append(Tx, WriteAAnother)

			Producer := &pipelinepoc.Producer{
				Nodes: make([]*pipelinepoc.Node, 0),
			}
			Producer.TreeMaking(Tx)
			Expect("a").Should(Equal(Producer.Nodes[0].WKey))
			Expect(1).Should(Equal(Producer.Nodes[0].WKeyVersion))
			Expect("a").Should(Equal(Producer.Nodes[0].Left.WKey))
			Expect(1).Should(Equal(Producer.Nodes[0].Left.WKeyVersion))
		})

		It("update, read", func() {
			WriteA := NodeCreation("", -1, "a", 1)
			ReadA := NodeCreation("a", 0, "", -1)
			Tx := make([]*pipelinepoc.Node, 0)
			Tx = append(Tx, WriteA)
			Tx = append(Tx, ReadA)

			Producer := &pipelinepoc.Producer{
				Nodes: make([]*pipelinepoc.Node, 0),
			}
			Producer.TreeMaking(Tx)
			Expect("a").Should(Equal(Producer.Nodes[0].WKey))
			Expect(1).Should(Equal(Producer.Nodes[0].WKeyVersion))
			Expect("a").Should(Equal(Producer.Nodes[0].Left.RKey))
			Expect(0).Should(Equal(Producer.Nodes[0].Left.RKeyVersion))
		})
	})
})
