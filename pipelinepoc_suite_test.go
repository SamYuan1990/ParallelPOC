package pipelinepoc_test

import (
	"testing"

	"github.com/SamYuan1990/pipelinepoc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPipelinepoc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pipelinepoc Suite")
}

func NodeCreation(RKey string, RKeyVersion int, WKey string, WKeyVersion int) *pipelinepoc.Node {
	return &pipelinepoc.Node{
		RKey:        RKey,
		RKeyVersion: RKeyVersion,
		WKey:        WKey,
		WKeyVersion: WKeyVersion,
		Input:       0,
		Left:        nil,
		Right:       nil,
	}
}
