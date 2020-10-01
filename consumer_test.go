package pipelinepoc_test

import (
	"github.com/SamYuan1990/pipelinepoc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Main", func() {
	Context("All action convert into tree mode", func() {
		// W A W B
		It("should supports WA and WB into 2 nodes", func() {
			Create_A := &pipelinepoc.Node{
				RKey:        "",
				RKeyVersion: -1,
				WKey:        "a",
				WKeyVersion: 0,
				Input:       0,
				Left:        nil,
				Right:       nil,
			}

			Create_B := &pipelinepoc.Node{
				RKey:        "",
				RKeyVersion: -1,
				WKey:        "b",
				WKeyVersion: 0,
				Input:       0,
				Left:        nil,
				Right:       nil,
			}

			Tx := make([]*pipelinepoc.Node, 0)
			Tx = append(Tx, Create_A)
			Tx = append(Tx, Create_B)
			consumer := &pipelinepoc.Consumer{
				Nodes: make([]*pipelinepoc.Node, 0),
			}
			consumer.TreeMaking(Tx)
			Expect("a").Should(Equal(consumer.Nodes[0].WKey))
			Expect("b").Should(Equal(consumer.Nodes[1].WKey))
		})
		// W A W B W A
		It("should supports WA and WB into 2 nodes", func() {
			Create_A := &pipelinepoc.Node{
				RKey:        "",
				RKeyVersion: -1,
				WKey:        "a",
				WKeyVersion: 0,
				Input:       0,
				Left:        nil,
				Right:       nil,
			}

			Create_B := &pipelinepoc.Node{
				RKey:        "",
				RKeyVersion: -1,
				WKey:        "b",
				WKeyVersion: 0,
				Input:       0,
				Left:        nil,
				Right:       nil,
			}

			Update_A := &pipelinepoc.Node{
				RKey:        "",
				RKeyVersion: -1,
				WKey:        "a",
				WKeyVersion: 1,
				Input:       0,
				Left:        nil,
				Right:       nil,
			}
			Tx := make([]*pipelinepoc.Node, 0)
			Tx = append(Tx, Create_A)
			Tx = append(Tx, Create_B)
			Tx = append(Tx, Update_A)
			consumer := &pipelinepoc.Consumer{
				Nodes: make([]*pipelinepoc.Node, 0),
			}
			consumer.TreeMaking(Tx)
			Expect("a").Should(Equal(consumer.Nodes[0].WKey))
			Expect(0).Should(Equal(consumer.Nodes[0].WKeyVersion))
			Expect("a").Should(Equal(consumer.Nodes[0].Right.WKey))
			Expect(1).Should(Equal(consumer.Nodes[0].Right.WKeyVersion))
			Expect("b").Should(Equal(consumer.Nodes[1].WKey))
		})

		// For second block comming, W A, W B, W A/W C, W D, R A W D
		It("should supports appending next block", func() {
			Create_A := &pipelinepoc.Node{
				RKey:        "",
				RKeyVersion: -1,
				WKey:        "a",
				WKeyVersion: 0,
				Input:       0,
				Left:        nil,
				Right:       nil,
			}

			Create_B := &pipelinepoc.Node{
				RKey:        "",
				RKeyVersion: -1,
				WKey:        "b",
				WKeyVersion: 0,
				Input:       0,
				Left:        nil,
				Right:       nil,
			}

			Update_A := &pipelinepoc.Node{
				RKey:        "",
				RKeyVersion: -1,
				WKey:        "a",
				WKeyVersion: 1,
				Input:       0,
				Left:        nil,
				Right:       nil,
			}

			Create_C := &pipelinepoc.Node{
				RKey:        "",
				RKeyVersion: -1,
				WKey:        "c",
				WKeyVersion: 0,
				Input:       0,
				Left:        nil,
				Right:       nil,
			}

			Create_D := &pipelinepoc.Node{
				RKey:        "",
				RKeyVersion: -1,
				WKey:        "d",
				WKeyVersion: 0,
				Input:       0,
				Left:        nil,
				Right:       nil,
			}

			RA_WD := &pipelinepoc.Node{
				RKey:        "a",
				RKeyVersion: 1,
				WKey:        "d",
				WKeyVersion: 1,
				Input:       0,
				Left:        nil,
				Right:       nil,
			}

			Tx := make([]*pipelinepoc.Node, 0)
			Tx = append(Tx, Create_A)
			Tx = append(Tx, Create_B)
			Tx = append(Tx, Update_A)
			consumer := &pipelinepoc.Consumer{
				Nodes: make([]*pipelinepoc.Node, 0),
			}
			consumer.TreeMaking(Tx)
			Expect("a").Should(Equal(consumer.Nodes[0].WKey))
			Expect(0).Should(Equal(consumer.Nodes[0].WKeyVersion))
			Expect("a").Should(Equal(consumer.Nodes[0].Right.WKey))
			Expect(1).Should(Equal(consumer.Nodes[0].Right.WKeyVersion))
			Expect("b").Should(Equal(consumer.Nodes[1].WKey))
			Tx2 := make([]*pipelinepoc.Node, 0)
			Tx2 = append(Tx2, Create_C)
			Tx2 = append(Tx2, Create_D)
			Tx2 = append(Tx2, RA_WD)
			consumer.TreeMaking(Tx2)
			Expect("a").Should(Equal(consumer.Nodes[0].WKey))
			Expect(0).Should(Equal(consumer.Nodes[0].WKeyVersion))
			Expect("a").Should(Equal(consumer.Nodes[0].Right.WKey))
			Expect(1).Should(Equal(consumer.Nodes[0].Right.WKeyVersion))
			Expect("d").Should(Equal(consumer.Nodes[0].Right.Right.WKey))
			Expect("a").Should(Equal(consumer.Nodes[0].Right.Right.RKey))
			Expect(1).Should(Equal(consumer.Nodes[0].Right.Right.RKeyVersion))
			Expect(1).Should(Equal(consumer.Nodes[0].Right.Right.Input))
			Expect("b").Should(Equal(consumer.Nodes[1].WKey))
			Expect("c").Should(Equal(consumer.Nodes[2].WKey))
			Expect("d").Should(Equal(consumer.Nodes[3].WKey))
			Expect(1).Should(Equal(consumer.Nodes[3].Right.Input))
		})
	})

	Context("All Read supports", func() {
		//R A, R B
		It("should supports RA and RB into 2 nodes", func() {
			Read_A := &pipelinepoc.Node{
				RKey:        "a",
				RKeyVersion: 0,
				WKey:        "",
				WKeyVersion: -1,
				Input:       0,
				Left:        nil,
				Right:       nil,
			}

			Read_B := &pipelinepoc.Node{
				RKey:        "b",
				RKeyVersion: 0,
				WKey:        "",
				WKeyVersion: -1,
				Input:       0,
				Left:        nil,
				Right:       nil,
			}

			Tx := make([]*pipelinepoc.Node, 0)
			Tx = append(Tx, Read_A)
			Tx = append(Tx, Read_B)
			consumer := &pipelinepoc.Consumer{
				Nodes: make([]*pipelinepoc.Node, 0),
			}
			consumer.TreeMaking(Tx)
			Expect("a").Should(Equal(consumer.Nodes[0].RKey))
			Expect("b").Should(Equal(consumer.Nodes[1].RKey))
		})
		// R A, R A
		It("should supports RA only nodes", func() {

			Read_A := &pipelinepoc.Node{
				RKey:        "a",
				RKeyVersion: 0,
				WKey:        "",
				WKeyVersion: -1,
				Input:       0,
				Left:        nil,
				Right:       nil,
			}

			Tx := make([]*pipelinepoc.Node, 0)
			Tx = append(Tx, Read_A)
			Tx = append(Tx, Read_A)
			consumer := &pipelinepoc.Consumer{
				Nodes: make([]*pipelinepoc.Node, 0),
			}
			consumer.TreeMaking(Tx)
			Expect("a").Should(Equal(consumer.Nodes[0].RKey))
			Expect("a").Should(Equal(consumer.Nodes[0].Right.RKey))
		})

		It("should WR in same nodes", func() {
			Read_A := &pipelinepoc.Node{
				RKey:        "a",
				RKeyVersion: 0,
				WKey:        "",
				WKeyVersion: -1,
				Input:       0,
				Left:        nil,
				Right:       nil,
			}

			Create_A := &pipelinepoc.Node{
				RKey:        "",
				RKeyVersion: -1,
				WKey:        "a",
				WKeyVersion: 0,
				Input:       0,
				Left:        nil,
				Right:       nil,
			}
			Tx := make([]*pipelinepoc.Node, 0)
			Tx = append(Tx, Create_A)
			Tx = append(Tx, Read_A)
			consumer := &pipelinepoc.Consumer{
				Nodes: make([]*pipelinepoc.Node, 0),
			}
			consumer.TreeMaking(Tx)
			Expect("a").Should(Equal(consumer.Nodes[0].WKey))
			Expect("a").Should(Equal(consumer.Nodes[0].Right.RKey))
		})
	})

	Context("U can't parallel", func() {
		Read_A := &pipelinepoc.Node{
			RKey:        "a",
			RKeyVersion: 0,
			WKey:        "",
			WKeyVersion: -1,
			Input:       0,
			Left:        nil,
			Right:       nil,
		}

		Create_A := &pipelinepoc.Node{
			RKey:        "",
			RKeyVersion: -1,
			WKey:        "a",
			WKeyVersion: 0,
			Input:       0,
			Left:        nil,
			Right:       nil,
		}
		U := &pipelinepoc.Node{
			RKey:        "",
			RKeyVersion: 0,
			WKey:        "",
			WKeyVersion: -1,
			Input:       0,
			Left:        nil,
			Right:       nil,
			UFlag:       true,
		}
		Create_B := &pipelinepoc.Node{
			RKey:        "",
			RKeyVersion: -1,
			WKey:        "b",
			WKeyVersion: 0,
			Input:       0,
			Left:        nil,
			Right:       nil,
		}
		// WA , WB, U, RA
		It("U can't parallel", func() {
			Tx := make([]*pipelinepoc.Node, 0)
			Tx = append(Tx, Create_A)
			Tx = append(Tx, Create_B)
			Tx = append(Tx, U)
			Tx = append(Tx, Read_A)
			consumer := &pipelinepoc.Consumer{
				Nodes: make([]*pipelinepoc.Node, 0),
			}
			consumer.TreeMaking(Tx)
			Expect("a").Should(Equal(consumer.Nodes[0].WKey))
			Expect("b").Should(Equal(consumer.Nodes[1].WKey))
			Expect(true).Should(Equal(consumer.Nodes[2].UFlag))
			Expect("a").Should(Equal(consumer.Nodes[3].RKey))
		})
	})
})
