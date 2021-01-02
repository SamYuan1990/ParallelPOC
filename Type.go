package pipelinepoc

import "strings"

type Node struct {
	//Right *Node
	//Left  *Node
	Next  *Node
	Wkeys []Key
	Rkeys []Key
	//Dependency int
}

func (n *Node) Add(node *Node) {
	current := n
	for current.Next != nil {
		current = current.Next
	}
	current.Next = node
}

// TODO
func (n *Node) Match(keys string) bool {
	for _, key := range n.Wkeys {
		if strings.Contains(keys, key.Name) {
			//strings.Contains(keys,key.Name）{
			return true
		}
	}
	for _, key := range n.Rkeys {
		if strings.Contains(keys, key.Name) {
			return true
		}
	}
	return false
}

// TODO
func (n *Node) MergeKey(key string) string {
	return key + n.GetKeys()
}

// TODO
func (n *Node) GetKeys() string {
	var s string
	for _, key := range n.Wkeys {
		s = s + key.Name
	}
	for _, key := range n.Rkeys {
		s = s + key.Name
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
	Txs []TxImpl
}

type TxImpl struct {
	Wkeys []Key
	Rkeys []Key
}

func (b BlockImpl) GetTxs() []TxImpl {
	return b.Txs
}

func (tx TxImpl) GetWkeys() []Key {
	return tx.Wkeys
}

func (tx TxImpl) GetRkeys() []Key {
	return tx.Rkeys
}
