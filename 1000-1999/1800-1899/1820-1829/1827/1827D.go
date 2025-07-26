package main

import (
	"bufio"
	"fmt"
	"os"
)

const LOG = 20

type Fenwick struct {
	n    int
	tree []int
}

func NewFenwick(n int) *Fenwick {
	return &Fenwick{n: n, tree: make([]int, n+2)}
}

func (f *Fenwick) Add(idx, val int) {
	for idx <= f.n {
		f.tree[idx] += val
		idx += idx & -idx
	}
}

func (f *Fenwick) Sum(idx int) int {
	res := 0
	for idx > 0 {
		res += f.tree[idx]
		idx -= idx & -idx
	}
	return res
}

func (f *Fenwick) RangeSum(l, r int) int {
	if l > r {
		return 0
	}
	return f.Sum(r) - f.Sum(l-1)
}

var (
	parent   []int
	children [][]int
	tin      []int
	tout     []int
	depth    []int
	up       [][]int
	timer    int
)

func dfs(v int) {
	timer++
	tin[v] = timer
	for _, to := range children[v] {
		depth[to] = depth[v] + 1
		up[to][0] = v
		for j := 1; j < LOG; j++ {
			up[to][j] = up[up[to][j-1]][j-1]
		}
		dfs(to)
	}
	tout[v] = timer
}

func isAncestor(u, v int) bool {
	return tin[u] <= tin[v] && tout[v] <= tout[u]
}

func getChild(u, v int) int {
	// assumes u is ancestor of v and u != v
	d := depth[v] - depth[u] - 1
	x := v
	for k := LOG - 1; k >= 0; k-- {
		if d&(1<<k) != 0 {
			x = up[x][k]
		}
	}
	return x
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		parent = make([]int, n+1)
		children = make([][]int, n+1)
		for i := 2; i <= n; i++ {
			fmt.Fscan(reader, &parent[i])
			children[parent[i]] = append(children[parent[i]], i)
		}

		tin = make([]int, n+1)
		tout = make([]int, n+1)
		depth = make([]int, n+1)
		up = make([][]int, n+1)
		for i := range up {
			up[i] = make([]int, LOG)
		}
		timer = 0
		dfs(1)

		bit := NewFenwick(n)
		bit.Add(tin[1], 1)
		active := 1
		centroid := 1
		// compute initial max child size
		maxChild := 0
		for _, ch := range children[centroid] {
			sz := bit.RangeSum(tin[ch], tout[ch])
			if sz > maxChild {
				maxChild = sz
			}
		}

		answers := make([]int, n-1)
		for i := 2; i <= n; i++ {
			x := i
			bit.Add(tin[x], 1)
			active++

			// update child size info if centroid is ancestor of x
			if isAncestor(centroid, x) && centroid != x {
				ch := getChild(centroid, x)
				sz := bit.RangeSum(tin[ch], tout[ch])
				if sz > maxChild {
					maxChild = sz
				}
			}

			// move centroid if needed
			for {
				moved := false
				sizeCentroid := bit.RangeSum(tin[centroid], tout[centroid])
				if active-sizeCentroid > active/2 {
					centroid = parent[centroid]
					// rebuild maxChild
					maxChild = 0
					for _, ch := range children[centroid] {
						sz := bit.RangeSum(tin[ch], tout[ch])
						if sz > maxChild {
							maxChild = sz
						}
					}
					moved = true
				} else {
					if isAncestor(centroid, x) && centroid != x {
						ch := getChild(centroid, x)
						sz := bit.RangeSum(tin[ch], tout[ch])
						if sz > active/2 {
							centroid = ch
							// rebuild
							maxChild = 0
							for _, cc := range children[centroid] {
								szz := bit.RangeSum(tin[cc], tout[cc])
								if szz > maxChild {
									maxChild = szz
								}
							}
							moved = true
						} else if sz > maxChild {
							maxChild = sz
						}
					}
				}
				if !moved {
					break
				}
			}

			sizeCentroid := bit.RangeSum(tin[centroid], tout[centroid])
			diffParent := abs(active - 2*sizeCentroid)
			diffChild := active - 2*maxChild
			if diffParent < diffChild {
				answers[i-2] = diffParent
			} else {
				answers[i-2] = diffChild
			}
		}

		for i, v := range answers {
			if i+2 == n+1 {
				fmt.Fprint(writer, v)
			} else {
				fmt.Fprintf(writer, "%d ", v)
			}
		}
		fmt.Fprintln(writer)
	}
}
