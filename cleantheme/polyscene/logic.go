package polyscene

import (
	"github.com/vron/v2d"
)

func (c *Polyscene) SetBounds(bnd v2d.Rect) {
	c.mu.Lock()
	c.bnd = bnd
	c.mu.Unlock()
}

func (c *Polyscene) Bounds() (bnd v2d.Rect) {
	c.mu.Lock()
	bnd = c.bnd
	c.mu.Unlock()
	return
}
