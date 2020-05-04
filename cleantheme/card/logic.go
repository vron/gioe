package card

import (
	"math"
)

// UpdateValue sets the value the widget should show
func (c *Card) UpdateValue(val float64) {
	c.mu.Lock()
	if val > c.max {
		c.max = val
	}
	if val < c.min {
		c.min = val
	}
	c.val = val
	c.mu.Unlock()
}

func (c *Card) values() (val, min, max float64) {
	c.mu.Lock()
	val, min, max = c.val, c.min, c.max
	c.mu.Unlock()
	if min == math.MaxFloat64 {
		min = 0
	}
	if max == -math.MaxFloat64 {
		max = 0
	}
	return
}
