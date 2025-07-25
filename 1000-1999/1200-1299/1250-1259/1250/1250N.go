package main

import (
	"bufio"
	"fmt"
	"os"
)

// DSU structure for union find
type DSU struct {
	parent []int
	rank   []int
}

func NewDSU(n int) *DSU {
	d := &DSU{
		parent: make([]int, n),
		rank:   make([]int, n),
	}
	for i := 0; i < n; i++ {
		d.parent[i] = i
	}
	return d
}

func (d *DSU) Find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.Find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) Union(a, b int) {
	a = d.Find(a)
	b = d.Find(b)
	if a == b {
		return
	}
	if d.rank[a] < d.rank[b] {
		a, b = b, a
	}
	d.parent[b] = a
	if d.rank[a] == d.rank[b] {
		d.rank[a]++
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		x := make([]int64, n)
		y := make([]int64, n)

		contactFirst := make(map[int64]int)
		dsu := NewDSU(n)
		contactMap := make(map[int64][]int)

		for i := 0; i < n; i++ {
			fmt.Fscan(in, &x[i], &y[i])
			if j, ok := contactFirst[x[i]]; ok {
				dsu.Union(i, j)
			} else {
				contactFirst[x[i]] = i
			}
			if j, ok := contactFirst[y[i]]; ok {
				dsu.Union(i, j)
			} else {
				contactFirst[y[i]] = i
			}
			contactMap[x[i]] = append(contactMap[x[i]], i)
			contactMap[y[i]] = append(contactMap[y[i]], i)
		}

		// group wires by component
		compWires := make(map[int][]int)
		for i := 0; i < n; i++ {
			r := dsu.Find(i)
			compWires[r] = append(compWires[r], i)
		}

		// choose root component (component of wire 0)
		rootRep := dsu.Find(0)
		rootContact := x[0]

		operations := make([][3]int64, 0)

		for rep, wires := range compWires {
			if rep == rootRep {
				continue
			}
			// find a wire and endpoint to modify
			chosenWire := wires[0]
			oldContact := x[chosenWire]

			// search for contact with degree 1
			found := false
			for _, w := range wires {
				if len(contactMap[x[w]]) == 1 {
					chosenWire = w
					oldContact = x[w]
					found = true
					break
				}
				if len(contactMap[y[w]]) == 1 {
					chosenWire = w
					oldContact = y[w]
					found = true
					break
				}
			}
			if !found {
				// if no unique contact, change x of chosen wire
				chosenWire = wires[0]
				oldContact = x[chosenWire]
			}

			operations = append(operations, [3]int64{int64(chosenWire + 1), oldContact, rootContact})
		}

		fmt.Fprintln(out, len(operations))
		for _, op := range operations {
			fmt.Fprintf(out, "%d %d %d\n", op[0], op[1], op[2])
		}
	}
}
