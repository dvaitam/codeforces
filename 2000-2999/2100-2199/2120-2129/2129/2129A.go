package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type dsu struct {
	p []int
	r []int
}

func newDSU(n int) *dsu {
	p := make([]int, n)
	r := make([]int, n)
	for i := range p {
		p[i] = i
	}
	return &dsu{p: p, r: r}
}

func (d *dsu) find(x int) int {
	if d.p[x] != x {
		d.p[x] = d.find(d.p[x])
	}
	return d.p[x]
}

func (d *dsu) unite(a, b int) bool {
	a = d.find(a)
	b = d.find(b)
	if a == b {
		return false
	}
	if d.r[a] < d.r[b] {
		a, b = b, a
	}
	d.p[b] = a
	if d.r[a] == d.r[b] {
		d.r[a]++
	}
	return true
}

type edge struct {
	a, b      int
	idx       int
	mandatory bool
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		if _, err := fmt.Fscan(in, &n); err != nil {
			return
		}
		edges := make([]edge, n)
		maxCoord := 2 * n
		coverage := make([]int, maxCoord+2) // counts for unit segments [x, x+1]
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &edges[i].a, &edges[i].b)
			edges[i].idx = i + 1
			for x := edges[i].a; x < edges[i].b; x++ {
				coverage[x]++
			}
		}

		// Identify mandatory edges (unique coverage on some unit)
		for i := range edges {
			a, b := edges[i].a, edges[i].b
			mandatory := false
			for x := a; x < b; x++ {
				if coverage[x] == 1 {
					mandatory = true
					break
				}
			}
			edges[i].mandatory = mandatory
		}

		covered := make([]bool, maxCoord+2)
		uncovered := 0
		for x := 1; x <= maxCoord; x++ {
			if coverage[x] > 0 {
				uncovered++
			}
		}

		selected := make([]bool, n)
		d := newDSU(maxCoord + 1)

		// Add mandatory edges first.
		for i, e := range edges {
			if !e.mandatory {
				continue
			}
			selected[i] = true
			d.unite(e.a, e.b)
			for x := e.a; x < e.b; x++ {
				if !covered[x] {
					covered[x] = true
					uncovered--
				}
			}
		}

		// Collect optional edges with initial contribution for ordering.
		type opt struct {
			id   int
			gain int
			a, b int
		}
		opts := make([]opt, 0, n)
		for i, e := range edges {
			if e.mandatory {
				continue
			}
			gain := 0
			for x := e.a; x < e.b; x++ {
				if !covered[x] {
					gain++
				}
			}
			opts = append(opts, opt{id: i, gain: gain, a: e.a, b: e.b})
		}

		sort.Slice(opts, func(i, j int) bool {
			if opts[i].gain == opts[j].gain {
				// longer interval later so tie-break by shorter index to diversify
				return (opts[i].b - opts[i].a) > (opts[j].b - opts[j].a)
			}
			return opts[i].gain > opts[j].gain
		})

		// First pass: add edges that help cover uncovered positions without creating cycles.
		for _, o := range opts {
			if uncovered == 0 {
				break
			}
			if selected[o.id] {
				continue
			}
			actual := 0
			for x := o.a; x < o.b; x++ {
				if !covered[x] {
					actual++
				}
			}
			if actual == 0 {
				continue
			}
			if d.unite(edges[o.id].a, edges[o.id].b) {
				selected[o.id] = true
				for x := o.a; x < o.b; x++ {
					if !covered[x] {
						covered[x] = true
						uncovered--
					}
				}
			}
		}

		// Second pass: if something is still uncovered, add edges even if they create cycles.
		if uncovered > 0 {
			for _, o := range opts {
				if uncovered == 0 {
					break
				}
				if selected[o.id] {
					continue
				}
				need := false
				for x := o.a; x < o.b; x++ {
					if !covered[x] {
						need = true
						break
					}
				}
				if !need {
					continue
				}
				selected[o.id] = true
				for x := o.a; x < o.b; x++ {
					if !covered[x] {
						covered[x] = true
						uncovered--
					}
				}
			}
		}

		// Collect result indices.
		res := make([]int, 0, n)
		for i, sel := range selected {
			if sel {
				res = append(res, edges[i].idx)
			}
		}
		fmt.Fprintln(out, len(res))
		for i, v := range res {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, v)
		}
		if len(res) > 0 {
			fmt.Fprintln(out)
		}
	}
}
