package pipelinepoc_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Main", func() {
	Context("All Read and Query should have count down lock", func() {
		// single processor
	})

	Context("w and w can parallel", func() {
		// two processors
	})

	Context("u break parallel", func() {
		// two processors
	})
})
