package pipelinepoc

import (
	"sync"
)

type Consumer struct {
	Pipeline *Pipeline
	Wg       *sync.WaitGroup
	stopped  bool
	suspend  bool
}

func (c *Consumer) Consume() {
	for !c.stopped {
		if c.suspend {
			continue
		}
		// get from pipeline
		_, value, ok := c.Pipeline.CCurrent.RemoveOldest()
		if ok {
			tmp := value.(*Node)
			//c.Pipeline.CCurrent.Remove(key)
			for tmp != nil {
				//time.Sleep(5 * time.Millisecond)
				tmp.Tx.Processed = true
				tmp = tmp.Next
			}
		}
		if c.Pipeline.CCurrent.Len() == 0 && c.Pipeline.CNext.Len() > 0 {
			//wait group here
			c.Wg.Done()
			c.suspend = true
			//c.Wg.Wait()
		}
	}
}

func (c *Consumer) Stop() {
	c.stopped = true
}

func (c *Consumer) Resume() {
	c.suspend = false
}
