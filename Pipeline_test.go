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
		Expect(p.Comming).NotTo(BeNil())
		Expect(p.Current).NotTo(BeNil())
		Expect(p.Current.Len()).Should(Equal(0))
		Expect(p.Next).NotTo(BeNil())
		Expect(p.Next.Len()).Should(Equal(0))
		Expect(p.Output).NotTo(BeNil())
		Expect(p.SwitchAble()).Should(BeTrue())
	})

	It("SwitchP", func() {
		p := &pipelinepoc.Pipeline{}
		err := p.Init()
		Expect(err).NotTo(HaveOccurred())
		Expect(p).NotTo(BeNil())
		Expect(p.Comming).NotTo(BeNil())
		current := p.PCurrent
		Expect(p.Current).NotTo(BeNil())
		Expect(p.Current.Len()).Should(Equal(0))
		next := p.PNext
		Expect(p.Next).NotTo(BeNil())
		Expect(p.Next.Len()).Should(Equal(0))
		Expect(p.Output).NotTo(BeNil())
		p.SwitchP()
		Expect(p.PNext).Should(Equal(current))
		Expect(p.PCurrent).Should(Equal(next))
		Expect(p.SwitchAble()).Should(BeFalse())
	})

	It("SwitchC", func() {
		p := &pipelinepoc.Pipeline{}
		err := p.Init()
		Expect(err).NotTo(HaveOccurred())
		Expect(p).NotTo(BeNil())
		Expect(p.Comming).NotTo(BeNil())
		current := p.CCurrent
		Expect(p.Current).NotTo(BeNil())
		Expect(p.Current.Len()).Should(Equal(0))
		next := p.CNext
		Expect(p.Next).NotTo(BeNil())
		Expect(p.Next.Len()).Should(Equal(0))
		Expect(p.Output).NotTo(BeNil())
		p.SwitchC()
		Expect(p.CNext).Should(Equal(current))
		Expect(p.CCurrent).Should(Equal(next))
		Expect(p.SwitchAble()).Should(BeFalse())
	})
})
