package pipelinepoc

import (
	"fmt"
	"time"
)

type Processor struct {
	Name string
}

func (p *Processor) Process(ch chan *Node, done chan bool) []*Node {
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
				time.Sleep(time.Nanosecond * 3)
			}
			if cur.WKey != "" {
				time.Sleep(time.Nanosecond * 5)
			}
			rs = append(rs, cur)
			if cur.Right != nil {
				if cur.Right.Input > 0 {
					cur.Right.ReduceInput()
					if cur.Right.Right != nil {
						fmt.Println(p.Name, " move RR ", cur.Right.Right)
						ch <- cur.Right.Right
					}
					if cur.Right.Left != nil {
						fmt.Println("move RL")
						ch <- cur.Right.Left
					}
					break
				}
				ns.Push(cur.Right)
			}
		}
	}
	fmt.Println("exit ", p.Name)
	done <- true
	//todo:dequeue
	return rs
}
