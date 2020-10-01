package pipelinepoc

// Dequeue // a,b,c
func TreeMaking(current *ToBeProcessQueue, txs []*Node) *ToBeProcessQueue {
	for _, v := range txs {
		find := false
		if v.UFlag {
			current = current.AppendU(v)
			continue
		}
		for _, vc := range current.GetButtom().Val {
			if vc.WKey == v.WKey && vc.WKey != "" {
				vc.AddWriteKey(v)
				find = true
			}
			if (vc.WKey == v.RKey || vc.RKey == v.RKey) && v.RKey != "" {
				vc.AddReadKey(v)
				v.InputAdd()
			}
		}
		if !find {
			current = current.Append(v)
		}
	}
	return current
}
