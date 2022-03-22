package pipeline

import (
	"github.com/nilsbu/lastfm/pkg/charts"
)

type graph struct {
	root    *node
	counter int
	caches  [][]string
	limit   int
}

type node struct {
	children   map[string]*node
	charts     charts.Charts
	lastAccess int
}

func newGraph(limit int) *graph {
	return &graph{
		root:   &node{children: map[string]*node{}},
		caches: [][]string{},
		limit:  limit,
	}
}

func (c *graph) get(steps []string) charts.Charts {
	c.counter++
	n := c.find(steps, c.counter)
	if n != nil {
		return n.charts
	} else {
		return nil
	}
}

func (c *graph) set(steps []string, charts charts.Charts) charts.Charts {
	c.counter++ // TODO don't set counter on set()
	n := c.find(steps[:len(steps)-1], c.counter)
	if n != nil {
		for i := 0; i < len(c.caches); {
			if isPredecessor(steps, c.caches[i]) && len(steps) < len(c.caches[i]) {
				c.removeCache(i)
			} else {
				i++
			}
		}
		n.children[steps[len(steps)-1]] = &node{
			map[string]*node{},
			charts,
			c.counter,
		}

		if len(steps) == 1 || steps[len(steps)-1] == "cache" {
			c.caches = append(c.caches, steps)
			c.prune()
		}

		return charts
	} else {
		return nil
	}
}

func (c *graph) find(steps []string, counter int) *node {
	n := c.root
	if counter > 0 {
		n.lastAccess = counter
	}
	for _, step := range steps {
		if next, ok := n.children[step]; ok {
			n = next
			if counter > 0 {
				n.lastAccess = counter
			}
		} else {
			return nil
		}
	}
	return n
}

func (c *graph) prune() {
	for len(c.caches) > c.limit {
		oldestAge, oldestId := -1, -1
		for i, cache := range c.caches {
			if isPredecessor(cache, c.caches[len(c.caches)-1]) {
				continue
			}
			age := c.find(cache, -1).lastAccess
			if oldestAge < age ||
				(oldestAge == age && len(cache) > len(c.caches[oldestId])) {
				oldestAge = age
				oldestId = i
			}
		}
		if oldestId != -1 {
			c.removeCache(oldestId)
		} else {
			break
		}
	}
}

func (c *graph) removeCache(cacheId int) {
	// remove from nodes (includes children)
	steps := c.caches[cacheId]
	parent := c.find(steps[:len(steps)-1], -1)
	delete(parent.children, steps[len(steps)-1])

	// remove from caches list
	newList := make([][]string, 0)
	newList = append(newList, c.caches[:cacheId]...)
	newList = append(newList, c.caches[cacheId+1:]...)
	c.caches = newList
}

func isPredecessor(steps, comp []string) bool {
	if len(steps) > len(comp) {
		return false
	}
	for i, step := range steps {
		if step != comp[i] {
			return false
		}
	}
	return true
}
