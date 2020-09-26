package pipelinepoc_test

import (
	"github.com/SamYuan1990/pipelinepoc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Main", func() {

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

	Context("All W should before R", func() {
		//R A, R B
		It("should supports RA and RB into 2 nodes", func() {
			Tx := make([]*pipelinepoc.Node, 0)
			Tx = append(Tx, Read_A)
			Tx = append(Tx, Read_B)
			current := &pipelinepoc.ToBeProcessQueue{}
			current = current.Init()
			current = pipelinepoc.TreeMaking(current, Tx)
			Expect("a").Should(Equal(current.GetTop().Val[0].RKey))
			Expect("b").Should(Equal(current.GetTop().Val[1].RKey))
		})
		// R A, R A
		It("should supports RA only nodes", func() {
			Tx := make([]*pipelinepoc.Node, 0)
			Tx = append(Tx, Read_A)
			Tx = append(Tx, Read_A)
			current := &pipelinepoc.ToBeProcessQueue{}
			current = current.Init()
			current = pipelinepoc.TreeMaking(current, Tx)
			Expect("a").Should(Equal(current.GetTop().Val[0].RKey))
			Expect("a").Should(Equal(current.GetTop().Val[1].Right.RKey))
		})
	})
})
