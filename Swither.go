package pipelinepoc

import "sync"

type Switcher struct {
	Pipeline  *Pipeline
	Wg        *sync.WaitGroup
	Consumers []*Consumer
	Count     int
	stopped   bool
}

func (s *Switcher) Init() {
	s.Wg.Add(s.Count)
}

func (s *Switcher) AddConsumer(c *Consumer) {
	if s.Consumers == nil {
		s.Consumers = make([]*Consumer, 0)
	}
	s.Consumers = append(s.Consumers, c)
}

func (s *Switcher) Switch() {
	for !s.stopped {
		s.Wg.Wait()
		s.Pipeline.SwitchC()
		for _, consumer := range s.Consumers {
			consumer.Resume()
		}
		s.Wg.Add(s.Count)
	}
}

func (s *Switcher) Stop() {
	s.stopped = true
}
