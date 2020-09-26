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

	Create_B := &pipelinepoc.Node{
		RKey:        "",
		RKeyVersion: -1,
		WKey:        "b",
		WKeyVersion: 0,
		Input:       0,
		Left:        nil,
		Right:       nil,
	}

	Read_A := &pipelinepoc.Node{
		RKey:        "a",
		RKeyVersion: 0,
		WKey:        "",
		WKeyVersion: -1,
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

	Context("U can't parallel", func() {
		// WA , WB, U, RA
		It("U can't parallel", func() {
			Tx := make([]*pipelinepoc.Node, 0)
			Tx = append(Tx, Create_A)
			Tx = append(Tx, Create_B)
			Tx = append(Tx, U)
			Tx = append(Tx, Read_A)
			current := &pipelinepoc.ToBeProcessQueue{}
			current = current.Init()
			current = pipelinepoc.TreeMaking(current, Tx)
			Expect("a").Should(Equal(current.GetTop().Val[0].WKey))
			Expect("b").Should(Equal(current.GetTop().Val[1].WKey))
			Expect(true).Should(Equal(current.GetTop().Val[2].UFlag))
			Expect("a").Should(Equal(current.GetButtom().Val[0].RKey))
		})
	})

})
