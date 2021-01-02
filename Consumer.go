package pipelinepoc

import (
	"time"

	lru "github.com/hashicorp/golang-lru"
)

type Consumer struct {
	LRU  *lru.Cache
	stop bool
}

func (c *Consumer) Consume(output chan *Node) {
	for {
		_, value, ok := c.LRU.RemoveOldest()
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
