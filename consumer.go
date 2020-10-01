package pipelinepoc

type Consumer struct {
	Nodes []*Node
	NextU int
}

// Dequeue // a,b,c
func (c *Consumer) TreeMaking(txs []*Node) {
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
