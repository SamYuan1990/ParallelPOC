package pipelinepoc

import (
	"time"
)

type Consumer struct {
	Pipeline *Pipeline
	stop     bool
}

func (c *Consumer) Consume(output chan *Node) {
	for {
		_, value, ok := c.Pipeline.RemoveOldest()
		if ok {
			node := value.(*Node)
			time.Sleep(time.Millisecond * 30)
			//fmt.Println(node)
			for node != nil {
				output <- node
				node = node.Next
			}
		}
		if c.stop {
			return
		}
	}
}

func (c *Consumer) Stop() {
	c.stop = true
}
