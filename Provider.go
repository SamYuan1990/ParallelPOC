package pipelinepoc

import (
	"strings"
)

// TODO
// Type Provider interface{}

type ProviderImpl struct {
	Pipeline *Pipeline
}

func (impl *ProviderImpl) Convert(b *BlockImpl) {
	// For each tx in block
	for _, tx := range b.GetTxs() {
		// Convert block into tx
		// Put Node into pipeline
		node := impl.ConvertTxToNode(&tx)
		keyStr := node.GetKeys()
		Keys := impl.Pipeline.Current.Keys()
		found := 0
		var nodeInLRU *Node
		for _, key := range Keys {
			if strings.Contains(key.(string), keyStr) || strings.Contains(keyStr, key.(string)) {
				v, ok := impl.Pipeline.Current.Get(key)
				if ok {
					nodeInLRU = v.(*Node)
				}
				found = found + 1
			}
		}
		if found == 1 {
			impl.Pipeline.Current.Remove(nodeInLRU.GetKeys())
			nodeInLRU.Add(node)
			impl.Pipeline.Add(nodeInLRU.MergeKey(keyStr), nodeInLRU)
		} else {
			// TODO
			impl.Pipeline.Add(keyStr, node)
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
