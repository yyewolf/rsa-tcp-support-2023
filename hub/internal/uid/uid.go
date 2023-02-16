package uid

import "sync"

type Generator struct {
	last int
	*sync.Mutex
}

func NewGenerator() *Generator {
	return &Generator{
		Mutex: &sync.Mutex{},
	}
}

var g = NewGenerator()

func init() {
	g.last = 0
}

func New() int {
	g.Lock()
	defer g.Unlock()

	g.last++
	return g.last
}
