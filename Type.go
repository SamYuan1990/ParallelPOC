package pipelinepoc

import (
	"strings"
	"sync"
)

type Node struct {
	//Right *Node
	//Left  *Node
	Next *Node
	Tx   *TxImpl
	lock sync.Mutex
	//Dependency int
}

func (n *Node) Add(node *Node) {
	current := n
	for current != nil && current.Next != nil {
		current = current.Next
	}
	if current != nil {
		current.Next = node
	}
}

// TODO
func (n *Node) Match(keys string) bool {
	for _, key := range n.Tx.Wkeys {
		if strings.Contains(keys, key.Name) {
			return true
		}
	}
	for _, key := range n.Tx.Rkeys {
		if strings.Contains(keys, key.Name) {
			return true
		}
	}
	return false
}

// TODO
func (n *Node) MergeKey(key string) string {
	return n.GetKeys() + key
}

// TODO
func (n *Node) GetKeys() string {
	var s string
	if n != nil {
		for _, key := range n.Tx.Wkeys {
			s = s + key.Name
		}
		for _, key := range n.Tx.Rkeys {
			s = s + key.Name
		}
	}
	return s
}

type Key struct {
	Name    string
	Version int
}

/* ToDo
type Block interface {
	GetTxs() []Tx
}

type Tx interface {
	GetWkeys() []Key
	GetRkeys() []Key
}
*/
type BlockImpl struct {
	Txs []*TxImpl
}

type TxImpl struct {
	Wkeys     []Key
	Rkeys     []Key
	Processed bool
}

func (b BlockImpl) GetTxs() []*TxImpl {
	return b.Txs
}

func (tx TxImpl) GetWkeys() []Key {
	return tx.Wkeys
}

func (tx TxImpl) GetRkeys() []Key {
	return tx.Rkeys
}
