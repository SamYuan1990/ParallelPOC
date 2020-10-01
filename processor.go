package pipelinepoc

func Process(tbpq *ToBeProcessQueue) []*Node {
	rs := make([]*Node, 0)
	currentTxs := tbpq.GetTop().Val
	for _, cur := range currentTxs { //paralle， replace by channel
		// go thread logic
		// middle, right
		ns := &NodeStack{
			Stack: make([]*Node, 10),
			Size:  0,
		}
		ns.Push(cur)
		for ns.Size > 0 {
			cur = ns.Pop()
			cur.Processed = true
			rs = append(rs, cur)
			if cur.Right != nil {
				if cur.Right.Input > 0 {
					cur.Right.ReduceInput()
					//append
					/*
						currentTxs.append cur.Right.Right
						currentTxs.append cur.Right.Right
					*/
					break
				}
				ns.Push(cur.Right)
			}
		}
	}
	//todo:dequeue
	return rs
}
