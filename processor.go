package pipelinepoc

import (
	"fmt"
	"time"
)

type Processor struct {
	Name string
}

func (p *Processor) Process(ch chan *Node) []*Node {
	rs := make([]*Node, 0)
	fmt.Println("Start Process")
	for {
		cur, ok := <-ch //paralle， replace by channel
		if !ok {
			break
		}
		if cur.FinalFlag {
			break
		}
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
			fmt.Println(p.Name, " Process ", cur)
			if cur.RKey != "" {
				time.Sleep(time.Second * 3)
			}
			if cur.WKey != "" {
				time.Sleep(time.Second * 5)
			}
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
