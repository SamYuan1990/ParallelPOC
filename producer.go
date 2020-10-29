package pipelinepoc

type Producer struct {
	Nodes []*Node
	NextU int
}

// Dequeue // a,b,c
func (c *Producer) TreeMaking(txs []*Node) {
	for _, v := range txs {
		find := false
		if v.UFlag {
			c.NextU = len(c.Nodes)
			c.Nodes = append(c.Nodes, v)
			continue
		}
		for i, vc := range c.Nodes {
			if i >= c.NextU {
				if vc.WKey == v.WKey && vc.WKey != "" {
					vc.AddWriteKey(v)
					find = true
				}
				if (vc.WKey == v.RKey || vc.RKey == v.RKey) && v.RKey != "" {
					vc.AddReadKey(v)
					v.InputAdd()
					find = true
				}
			}
		}
		if !find {
			c.Nodes = append(c.Nodes, v)
		}
	}
}

func (c *Producer) IntoChan(ch chan *Node) {
	for _, v := range c.Nodes {
		ch <- v
	}
	ch <- &Node{FinalFlag: true}
	//close(ch)
}