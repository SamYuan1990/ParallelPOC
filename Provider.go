package pipelinepoc

import (
	"strings"
)

// TODO
// Type Provider interface{}

type ProviderImpl struct {
	Pipeline *Pipeline
	stopped  bool
}

func (impl *ProviderImpl) Convert() {
	for !impl.stopped {
		v, err := impl.Pipeline.Comming.Dequeue()
		if err != nil {
			return
		}
		if v == nil {
			continue
		}
		b := v.(*BlockImpl)
		if b == nil {
			continue
		}
		// For each tx in block
		for _, tx := range b.GetTxs() {
			// Convert block into tx
			// Put Node into pipeline
			node := impl.ConvertTxToNode(tx)
			keyStr := node.GetKeys()
			Keys := impl.Pipeline.PCurrent.Keys()
			found := 0
			var nodeInLRU *Node
			for _, key := range Keys {
				if strings.Contains(key.(string), keyStr) || strings.Contains(keyStr, key.(string)) {
					v, ok := impl.Pipeline.PCurrent.Get(key)
					if ok {
						if v != nil {
							nodeInLRU = v.(*Node)
							found = found + 1
						}
					}
				}
			}
			if found == 1 {
				GetLogger().Info("update existing key and node", keyStr, node)
				impl.Pipeline.PCurrent.Remove(nodeInLRU.GetKeys())
				nodeInLRU.Add(node)
				impl.Pipeline.PCurrent.PeekOrAdd(nodeInLRU.MergeKey(keyStr), nodeInLRU)
			}
			if found == 0 {
				GetLogger().Info("add new key", keyStr, node)
				impl.Pipeline.PCurrent.Add(keyStr, node)
			}
			if found > 1 {
				GetLogger().Info("will hold and waiting for switch")
				for !impl.Pipeline.SwitchAble() { // hold

				}
				impl.Pipeline.SwitchP()
				impl.Pipeline.PCurrent.Add(keyStr, node)
			}
		}
		GetLogger().Info("Put V to output queue ", v)
		impl.Pipeline.Output.Enqueue(v)
	}
}

// TODO: This is an interal function
// Export for test porpuse
func (impl *ProviderImpl) ConvertTxToNode(Tx *TxImpl) *Node {
	return &Node{
		Tx: Tx,
	}
}

func (impl *ProviderImpl) Stop() {
	impl.stopped = true
}
