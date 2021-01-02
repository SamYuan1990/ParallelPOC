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

func ConstructBlocks(key string, anotherkey string) *pipelinepoc.BlockImpl {
	Wkeys := make([]pipelinepoc.Key, 0)
	keys := ConstructKey(key, 0)
	anotherkeys := ConstructKey(anotherkey, 0)
	Wkeys = append(Wkeys, keys, anotherkeys)
	tximpl := pipelinepoc.TxImpl{
		Wkeys: Wkeys,
	}
	txs := make([]pipelinepoc.TxImpl, 0)
	txs = append(txs, tximpl)

	BlockImpl := &pipelinepoc.BlockImpl{
		Txs: txs,
	}
	return BlockImpl
}

func ConstructBlock(key string, value int) *pipelinepoc.BlockImpl {
	Wkeys := make([]pipelinepoc.Key, 0)
	keys := ConstructKey(key, value)
	Wkeys = append(Wkeys, keys)
	tximpl := pipelinepoc.TxImpl{
		Wkeys: Wkeys,
	}
	txs := make([]pipelinepoc.TxImpl, 0)
	txs = append(txs, tximpl)

	BlockImpl := &pipelinepoc.BlockImpl{
		Txs: txs,
	}
	return BlockImpl
}

func ConstructKey(key string, value int) pipelinepoc.Key {
	return pipelinepoc.Key{
		Name:    key,
		Version: value,
	}
}

func ConstructNodeWithSingleKey(key string, value int) *pipelinepoc.Node {
	Wkeys0 := make([]pipelinepoc.Key, 0)
	key0 := ConstructKey(key, value)
	Wkeys0 = append(Wkeys0, key0)
	return &pipelinepoc.Node{
		Wkeys: Wkeys0,
	}
}

func ConstructNodeWithTwoKey(key string, key1 string) *pipelinepoc.Node {
	Wkeys0 := make([]pipelinepoc.Key, 0)
	key_0 := ConstructKey(key, 0)
	key_1 := ConstructKey(key1, 0)
	Wkeys0 = append(Wkeys0, key_0, key_1)
	return &pipelinepoc.Node{
		Wkeys: Wkeys0,
	}
}
