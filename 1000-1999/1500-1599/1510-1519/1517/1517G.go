package main

import (
	"bufio"
	"fmt"
	"os"
)

type Point struct {
	x, y      int64
	w         int64
	important bool
}

type Pattern struct {
	a, b, c, d int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	pts := make([]Point, n)
	pos := make(map[[2]int64]int)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &pts[i].x, &pts[i].y, &pts[i].w)
		if pts[i].x%2 == 0 && pts[i].y%2 == 0 {
			pts[i].important = true
		}
		pos[[2]int64{pts[i].x, pts[i].y}] = i
	}

	var patterns []Pattern
	for i, p := range pts {
		if !p.important {
			continue
		}
		dirs := []int{-1, 1}
		for _, dx := range dirs {
			bx := p.x + int64(dx)
			by := p.y
			bIdx, okB := pos[[2]int64{bx, by}]
			if !okB {
				continue
			}
			var vecs [][2]int
			if dx == 1 {
				vecs = [][2]int{{0, 1}, {0, -1}, {-1, 1}, {-1, -1}}
			} else {
				vecs = [][2]int{{0, 1}, {0, -1}, {1, 1}, {1, -1}}
			}
			for _, v := range vecs {
				cx := p.x + int64(v[0])
				cy := p.y + int64(v[1])
				dx2 := bx + int64(v[0])
				dy2 := by + int64(v[1])
				cIdx, okC := pos[[2]int64{cx, cy}]
				dIdx, okD := pos[[2]int64{dx2, dy2}]
				if okC && okD {
					patterns = append(patterns, Pattern{i, bIdx, cIdx, dIdx})
				}
			}
		}
	}

	// build adjacency graph
	adj := make([][]int, n)
	for _, pat := range patterns {
		arr := []int{pat.a, pat.b, pat.c, pat.d}
		for i := 0; i < 4; i++ {
			for j := i + 1; j < 4; j++ {
				u := arr[i]
				v := arr[j]
				adj[u] = append(adj[u], v)
				adj[v] = append(adj[v], u)
			}
		}
	}

	// compute components
	compID := make([]int, n)
	for i := range compID {
		compID[i] = -1
	}
	var comps [][]int
	for i := 0; i < n; i++ {
		if compID[i] >= 0 {
			continue
		}
		stack := []int{i}
		compID[i] = len(comps)
		var comp []int
		for len(stack) > 0 {
			u := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			comp = append(comp, u)
			for _, v := range adj[u] {
				if compID[v] == -1 {
					compID[v] = compID[i]
					stack = append(stack, v)
				}
			}
		}
		comps = append(comps, comp)
	}

	// patterns per component
	pcomp := make([][]Pattern, len(comps))
	for _, pat := range patterns {
		id := compID[pat.a]
		pcomp[id] = append(pcomp[id], pat)
	}

	var total int64
	for idx, comp := range comps {
		m := len(comp)
		mapping := make(map[int]int, m)
		for i, v := range comp {
			mapping[v] = i
		}
		best := int64(0)
		patterns := pcomp[idx]
		if m <= 20 {
			limit := 1 << m
			for mask := 0; mask < limit; mask++ {
				ok := true
				for _, pat := range patterns {
					if ((mask>>mapping[pat.a])&1) == 1 &&
						((mask>>mapping[pat.b])&1) == 1 &&
						((mask>>mapping[pat.c])&1) == 1 &&
						((mask>>mapping[pat.d])&1) == 1 {
						ok = false
						break
					}
				}
				if !ok {
					continue
				}
				var sum int64
				for i, v := range comp {
					if (mask>>i)&1 == 1 {
						sum += pts[v].w
					}
				}
				if sum > best {
					best = sum
				}
			}
		} else {
			// fallback greedy: keep all, remove conflicting patterns iteratively
			keep := make([]bool, m)
			for i := range keep {
				keep[i] = true
			}
			for _, pat := range patterns {
				if keep[mapping[pat.a]] && keep[mapping[pat.b]] && keep[mapping[pat.c]] && keep[mapping[pat.d]] {
					// remove vertex with minimal weight
					ids := []int{pat.a, pat.b, pat.c, pat.d}
					minIdx := ids[0]
					for _, id := range ids[1:] {
						if pts[id].w < pts[minIdx].w {
							minIdx = id
						}
					}
					keep[mapping[minIdx]] = false
				}
			}
			var sum int64
			for i, v := range comp {
				if keep[i] {
					sum += pts[v].w
				}
			}
			best = sum
		}
		total += best
	}

	fmt.Fprintln(out, total)
}
