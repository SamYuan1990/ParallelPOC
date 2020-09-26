package pipelinepoc_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPipelinepoc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pipelinepoc Suite")
}
