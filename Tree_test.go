package pipelinepoc_test

import (
	"github.com/SamYuan1990/pipelinepoc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Main", func() {

	Create_A := &pipelinepoc.Node{
		RKey:        "",
		RKeyVersion: -1,
		WKey:        "a",
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

	RA_WD := &pipelinepoc.Node{
		RKey:        "a",
		RKeyVersion: 1,
		WKey:        "d",
		WKeyVersion: 1,
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

	Context("All action convert into tree mode", func() {
		// W A W B
		It("should supports WA and WB into 2 nodes", func() {
			Tx := make([]*pipelinepoc.Node, 0)
			Tx = append(Tx, Create_A)
			Tx = append(Tx, Create_B)
			current := &pipelinepoc.ToBeProcessQueue{}
			current = current.Init()
			current = pipelinepoc.TreeMaking(current, Tx)
			Expect("a").Should(Equal(current.GetTop().Val[0].WKey))
			Expect("b").Should(Equal(current.GetTop().Val[1].WKey))
		})
		// W A W B W A
		It("should supports WA and WB into 2 nodes", func() {
			Tx := make([]*pipelinepoc.Node, 0)
			Tx = append(Tx, Create_A)
			Tx = append(Tx, Create_B)
			Tx = append(Tx, Update_A)
			current := &pipelinepoc.ToBeProcessQueue{}
			current = current.Init()
			current = pipelinepoc.TreeMaking(current, Tx)
			Expect("a").Should(Equal(current.GetTop().Val[0].WKey))
			Expect(0).Should(Equal(current.GetTop().Val[0].WKeyVersion))
			Expect("a").Should(Equal(current.GetTop().Val[0].Right.WKey))
			Expect(1).Should(Equal(current.GetTop().Val[0].Right.WKeyVersion))
			Expect("b").Should(Equal(current.GetTop().Val[1].WKey))
		})
		// For second block comming, W A, W B, W A/W C, W D, R A W D
		It("should supports appending next block", func() {
			Tx := make([]*pipelinepoc.Node, 0)
			Tx = append(Tx, Create_A)
			Tx = append(Tx, Create_B)
			Tx = append(Tx, Update_A)
			current := &pipelinepoc.ToBeProcessQueue{}
			current = current.Init()
			current = pipelinepoc.TreeMaking(current, Tx)
			Expect("a").Should(Equal(current.GetTop().Val[0].WKey))
			Expect(0).Should(Equal(current.GetTop().Val[0].WKeyVersion))
			Expect("a").Should(Equal(current.GetTop().Val[0].Right.WKey))
			Expect(1).Should(Equal(current.GetTop().Val[0].Right.WKeyVersion))
			Expect("b").Should(Equal(current.GetTop().Val[1].WKey))
			Tx2 := make([]*pipelinepoc.Node, 0)
			Tx2 = append(Tx2, Create_C)
			Tx2 = append(Tx2, Create_D)
			Tx2 = append(Tx2, RA_WD)
			current = pipelinepoc.TreeMaking(current, Tx2)
			Expect("a").Should(Equal(current.GetTop().Val[0].WKey))
			Expect(0).Should(Equal(current.GetTop().Val[0].WKeyVersion))
			Expect("a").Should(Equal(current.GetTop().Val[0].Right.WKey))
			Expect(1).Should(Equal(current.GetTop().Val[0].Right.WKeyVersion))
			Expect("d").Should(Equal(current.GetTop().Val[0].Right.Right.WKey))
			Expect("a").Should(Equal(current.GetTop().Val[0].Right.Right.RKey))
			Expect(1).Should(Equal(current.GetTop().Val[0].Right.Right.RKeyVersion))
			Expect(1).Should(Equal(current.GetTop().Val[0].Right.Right.Input))
			Expect("b").Should(Equal(current.GetTop().Val[1].WKey))
			Expect("c").Should(Equal(current.GetTop().Val[2].WKey))
			Expect("d").Should(Equal(current.GetTop().Val[3].WKey))
			Expect(1).Should(Equal(current.GetTop().Val[3].Right.Input))
		})
	})
})
