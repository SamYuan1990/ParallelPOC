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

	Update_A := &pipelinepoc.Node{
		RKey:        "",
		RKeyVersion: -1,
		WKey:        "a",
		WKeyVersion: 1,
		Input:       0,
		Left:        nil,
		Right:       nil,
	}

	RA_WB := &pipelinepoc.Node{
		RKey:        "a",
		RKeyVersion: 1,
		WKey:        "b",
		WKeyVersion: 1,
		Input:       0,
		Left:        nil,
		Right:       nil,
	}

	Context("All Read and Query should have count down lock", func() {
		It("single process", func() {
			// wa wb ua
			Tx := make([]*pipelinepoc.Node, 0)
			Tx = append(Tx, Create_A)
			Tx = append(Tx, Create_B)
			Tx = append(Tx, Update_A)
			current := &pipelinepoc.ToBeProcessQueue{}
			current = current.Init()
			// output as wa/ua wb
			current = pipelinepoc.TreeMaking(current, Tx)
			verify := pipelinepoc.Process(current)
			Expect("a").Should(Equal(verify[0].WKey))
			Expect(0).Should(Equal(verify[0].WKeyVersion))
			Expect("a").Should(Equal(verify[1].WKey))
			Expect(1).Should(Equal(verify[1].WKeyVersion))
			Expect("b").Should(Equal(verify[2].WKey))
		})
	})

	Context("input !0 break", func() {
		// single processors
		It("single process", func() {
			// wa wb us ra_wb
			Tx := make([]*pipelinepoc.Node, 0)
			Tx = append(Tx, Create_A)
			Tx = append(Tx, Create_B)
			Tx = append(Tx, Update_A)
			Tx = append(Tx, RA_WB)
			current := &pipelinepoc.ToBeProcessQueue{}
			current = current.Init()
			// output as wa/ua wb/RA_WB
			current = pipelinepoc.TreeMaking(current, Tx)
			verify := pipelinepoc.Process(current)
			Expect("a").Should(Equal(verify[0].WKey))
			Expect(0).Should(Equal(verify[0].WKeyVersion))
			Expect("a").Should(Equal(verify[1].WKey))
			Expect(1).Should(Equal(verify[1].WKeyVersion))
			Expect("b").Should(Equal(verify[2].WKey))
			Expect("b").Should(Equal(verify[3].WKey))
			Expect("a").Should(Equal(verify[3].RKey))
		})
	})
	/*
		Context("u break parallel", func() {
			// two processors
		})*/
})
