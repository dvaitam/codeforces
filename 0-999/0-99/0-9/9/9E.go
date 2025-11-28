package main

import (
	"bufio"
	"fmt"
	"os"
)

type DSU struct {
	p  []int
	sz []int // size of component (vertices)
}

func NewDSU(n int) *DSU {
	d := &DSU{
		p:  make([]int, n+1),
		sz: make([]int, n+1),
	}
	for i := 1; i <= n; i++ {
		d.p[i] = i
		d.sz[i] = 1
	}
	return d
}

func (d *DSU) Find(x int) int {
	if d.p[x] != x {
		d.p[x] = d.Find(d.p[x])
	}
	return d.p[x]
}

func (d *DSU) Union(x, y int) bool {
	rx, ry := d.Find(x), d.Find(y)
	if rx != ry {
		// simple union logic, no rank needed for N=50
		d.p[ry] = rx
		d.sz[rx] += d.sz[ry]
		return true
	}
	return false
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}

	deg := make([]int, n+1)
	dsu := NewDSU(n)
	// We also need to track edge counts per component to detect existing cycles
	compEdges := make(map[int]int)

	// Initial check for single node N=1 and M=0 -> need to add loop
	// N=1 M=1 (1 1) -> valid.

	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		deg[u]++
		deg[v]++
		
		// Use a temp find to track edges before union changes roots
		ru := dsu.Find(u)
		rv := dsu.Find(v)
		
		if ru != rv {
			// merging two components
			edgesSum := compEdges[ru] + compEdges[rv] + 1
			dsu.Union(u, v)
			newRoot := dsu.Find(u)
			compEdges[newRoot] = edgesSum
			delete(compEdges, ru)
			delete(compEdges, rv)
			// Re-add newRoot if we deleted both? 
			// If ru was merged into rv (or vice versa), one ID remains. 
			// My Union implementation makes `ry` child of `rx`. So `rx` is new root.
			// Correct logic:
			// root of union is rx.
			// compEdges[rx] = edgesSum.
			// compEdges[ry] is effectively gone/stale. (Map cleanup)
		} else {
			// Adding edge within component
			compEdges[ru]++
		}
	}

	// 1. Validate Degrees
	for i := 1; i <= n; i++ {
		if deg[i] > 2 {
			fmt.Println("NO")
			return
		}
	}

	// 2. Validate Components (No premature cycles)
	// Update map keys to ensure we check current roots
	// Rebuild compEdges map based on current roots, just to be safe against stale keys
	// Actually easier: recalculate edges sum
	// But we need to distinguish "edges within component".
	// The tracking above is tricky if roots change.
	// Let's just recount edges. Iterate all edges? No, we didn't store them.
	// But we know `sum(deg) = 2 * edges`.
	// So sum of degrees in a component / 2 = edges in component.
	
	compDegSum := make(map[int]int)
	for i := 1; i <= n; i++ {
		root := dsu.Find(i)
		compDegSum[root] += deg[i]
	}
	
	for i := 1; i <= n; i++ {
		if dsu.p[i] == i { // is root
			edges := compDegSum[i] / 2
			verts := dsu.sz[i]
			
			if edges > verts {
				// More edges than vertices => >1 cycle or complex graph. Impossible for degree<=2.
				// Degree constraint <=2 implies max edges = verts (cycle). 
				// edges > verts implies some node has deg > 2, already checked.
				// So this branch might be redundant but safe.
				fmt.Println("NO")
				return
			}
			if edges == verts {
				// Contains a cycle
				if verts < n {
					// Premature cycle (does not include all nodes)
					fmt.Println("NO")
					return
				}
				// If verts == n, it's a full cycle. 
				// If m edges added and graph is valid full cycle, we are done?
				// We might need to output 0 edges.
			}
		}
	}

	var addedEdges [][2]int

	// 3. Greedy Addition
	// Repeat until we have exactly N edges total (and 1 component of size N).
	// Or simpler: repeat until all degrees are 2.
	
	for {
		// Check if finished
		allTwo := true
		for i := 1; i <= n; i++ {
			if deg[i] != 2 {
				allTwo = false
				break
			}
		}
		if allTwo {
			// Verify connectivity (should be 1 component)
			root := dsu.Find(1)
			if dsu.sz[root] != n {
				// This case implies we formed multiple disjoint cycles.
				// My greedy logic below prevents this (only close cycle if size == N).
				// So this shouldn't happen if logic is correct.
				fmt.Println("NO")
				return
			}
			break
		}

		added := false
		// Find lexicographically first valid edge
		for u := 1; u <= n; u++ {
			if deg[u] >= 2 {
				continue
			}
			// j starts at u for self-loop check (only if N=1)
			startV := u + 1
			if n == 1 {
				startV = u
			}
			
			for v := startV; v <= n; v++ {
				if deg[v] >= 2 {
					continue
				}
				
				rootU := dsu.Find(u)
				rootV := dsu.Find(v)
				
				// Valid conditions:
				// 1. Merging two different components (paths). Always allowed.
				// 2. Closing a cycle (same component). Only allowed if it forms the FINAL N-cycle.
				//    i.e., the component size must be N.
				
				valid := false
				if rootU != rootV {
					valid = true
				} else {
					// Same component. Check if this component contains all nodes.
					if dsu.sz[rootU] == n {
						valid = true
					}
				}
				
				if valid {
					deg[u]++
					deg[v]++
					dsu.Union(u, v)
					addedEdges = append(addedEdges, [2]int{u, v})
					added = true
					break
				}
			}
			if added {
				break
			}
		}
		
		if !added {
			// Could not add edge but degrees not all 2?
			// Should not happen for valid inputs solvable by this logic.
			fmt.Println("NO")
			return
		}
	}

	fmt.Println("YES")
	fmt.Println(len(addedEdges))
	for _, e := range addedEdges {
		fmt.Printf("%d %d\n", e[0], e[1])
	}
}

