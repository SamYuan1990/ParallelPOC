package pipelinepoc

import (
	"strings"

	lru "github.com/hashicorp/golang-lru"
)

// TODO
// Type Provider interface{}

type ProviderImpl struct {
	LRU *lru.Cache
}

func (impl *ProviderImpl) Convert(b *BlockImpl) {
	// For each tx in block
	for _, tx := range b.GetTxs() {
		// Convert block into tx
		// Put Node into pipeline
		node := impl.ConvertTxToNode(&tx)
		keyStr := node.GetKeys()
		Keys := impl.LRU.Keys()
		found := 0
		var nodeInLRU *Node
		for _, key := range Keys {
			if strings.Contains(key.(string), keyStr) || strings.Contains(keyStr, key.(string)) {
				v, ok := impl.LRU.Get(key)
				if ok {
					nodeInLRU = v.(*Node)
				}
				found = found + 1
			}
		}
		if found == 1 {
			impl.LRU.Remove(nodeInLRU.GetKeys())
			nodeInLRU.Add(node)
			impl.LRU.Add(nodeInLRU.MergeKey(keyStr), nodeInLRU)

		} else {
			impl.LRU.Add(keyStr, node)
		}
	}
}

// TODO: This is an interal function
// Export for test porpuse
func (impl *ProviderImpl) ConvertTxToNode(Tx *TxImpl) *Node {
	Wkeys := Tx.GetWkeys()
	Rkeys := Tx.GetRkeys()
	return &Node{
		Wkeys: Wkeys,
		Rkeys: Rkeys,
	}
}
