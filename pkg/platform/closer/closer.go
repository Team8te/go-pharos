package closer

import (
	"os"
	"os/signal"
	"sync"
)

type Closer struct {
	done chan struct{}
	once sync.Once
}

func New(sig ...os.Signal) *Closer {
	c := &Closer{done: make(chan struct{})}
	if len(sig) > 0 {
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, sig...)
			<-ch
			signal.Stop(ch)
			c.Close()
		}()
	}
	return c
}

// Wait blocks until all closer functions are done
func (c *Closer) Wait() {
	<-c.done
}

func (c *Closer) Close() {
	c.once.Do(func() {
		close(c.done)
	})
}
