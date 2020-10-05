package pipelinepoc

/*
For any of node, it should a set of Rkey, Wkey, as tx.
With in this poc, just use int instead.

For any Left, means sharing same latest W version
As:
W A,100
R A(left)
R A(left)
W A(+10 ,left)

For any Right, means sharing differe version of latest version
As:
R A, W B,100
R A, W B,90
either 100 or 110 will successed

For any node, if Input is not 0, means having dependency
*/
type Node struct {
	RKey        string
	RKeyVersion int
	WKey        string
	WKeyVersion int
	Input       int
	Left        *Node
	Right       *Node
	UFlag       bool
	Processed   bool
	FinalFlag   bool
}

func (n *Node) InputAdd() {
	n.Input = n.Input + 1
}

func (n *Node) ReduceInput() {
	n.Input--
}

func (n *Node) AddReadKey(input *Node) {
	if n.RKeyVersion <= input.RKeyVersion {
		//right
		if n.Right == nil {
			n.Right = input
		} else {
			n.Right.AddReadKey(input)
		}
	} else {
		//left
		if n.Left == nil {
			n.Left = input
		} else {
			n.Left.AddReadKey(input)
		}
	}
}

func (n *Node) AddWriteKey(input *Node) {
	if n.WKeyVersion < input.WKeyVersion {
		//right
		if n.Right == nil {
			n.Right = input
		} else {
			n.Right.AddWriteKey(input)
		}
	} else {
		//left
		if n.Left == nil {
			n.Left = input
		} else {
			n.Left.AddWriteKey(input)
		}
	}
}

type QueueNode struct {
	Val  chan *Node
	Next *QueueNode
}

type ToBeProcessChan struct {
	Head   *QueueNode
	Buttom *QueueNode
	// [a, b, c][a,b,c,u][a,b,c]
}

func (tbpq *ToBeProcessChan) GetTop() *QueueNode {
	return tbpq.Head
}

func (tbpq *ToBeProcessChan) GetButtom() *QueueNode {
	return tbpq.Buttom
}

func (tbpq *ToBeProcessChan) Init() *ToBeProcessChan {
	tbpq.Buttom = &QueueNode{
		Val:  make(chan *Node, 0),
		Next: nil,
	}
	tbpq.Head = tbpq.Buttom
	return tbpq
}

// EnQueue // a,b,c,u
func (tbpq *ToBeProcessChan) Append(n *Node) *ToBeProcessChan {
	tbpq.Buttom.Val <- n
	return tbpq
}

// EnQueue // a,b,c,u
func (tbpq *ToBeProcessChan) AppendU(n *Node) *ToBeProcessChan {
	tmp := &QueueNode{
		Val:  make(chan *Node, 0),
		Next: nil,
	}
	tbpq.Buttom.Next = tmp
	tbpq.Buttom = tmp
	tbpq.Buttom.Val <- n
	return tbpq
}

type NodeStack struct {
	Stack []*Node
	Size  int
}

func (ns *NodeStack) Pop() *Node {
	tmp := ns.Stack[ns.Size-1]
	ns.Size--
	return tmp
}

func (ns *NodeStack) Push(input *Node) {
	ns.Stack[ns.Size] = input
	ns.Size++
}
