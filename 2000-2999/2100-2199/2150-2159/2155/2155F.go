package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

type fastScanner struct {
	r *bufio.Reader
}

func newFastScanner() *fastScanner {
	return &fastScanner{r: bufio.NewReaderSize(os.Stdin, 1<<20)}
}

func (fs *fastScanner) nextInt() int {
	sign := 1
	val := 0
	c, err := fs.r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		if err != nil {
			return 0
		}
		c, err = fs.r.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, err = fs.r.ReadByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int(c-'0')
		c, err = fs.r.ReadByte()
		if err != nil {
			break
		}
	}
	return sign * val
}

type dsu struct {
	parent []int
	size   []int
}

func newDSU(n int) *dsu {
	parent := make([]int, n)
	size := make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
		size[i] = 1
	}
	return &dsu{parent: parent, size: size}
}

func (d *dsu) find(x int) int {
	for d.parent[x] != x {
		d.parent[x] = d.parent[d.parent[x]]
		x = d.parent[x]
	}
	return x
}

func (d *dsu) union(a, b int) {
	pa := d.find(a)
	pb := d.find(b)
	if pa == pb {
		return
	}
	if d.size[pa] < d.size[pb] {
		pa, pb = pb, pa
	}
	d.parent[pb] = pa
	d.size[pa] += d.size[pb]
}

func main() {
	fs := newFastScanner()
	out := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer out.Flush()

	t := fs.nextInt()
	for ; t > 0; t-- {
		n := fs.nextInt()
		k := fs.nextInt()
		s := fs.nextInt()
		q := fs.nextInt()

		edges := make([][2]int, n-1)
		for i := 0; i < n-1; i++ {
			u := fs.nextInt()
			v := fs.nextInt()
			edges[i][0] = u
			edges[i][1] = v
		}

		colors := make([][]int, n+1)
		nodesOfColor := make([][]int, k+1)
		for i := 0; i < s; i++ {
			v := fs.nextInt()
			c := fs.nextInt()
			colors[v] = append(colors[v], c)
			nodesOfColor[c] = append(nodesOfColor[c], v)
		}

		colorIndex := make([]map[int]int, n+1)
		compLeader := make([][]int, n+1)
		for node := 1; node <= n; node++ {
			if len(colors[node]) == 0 {
				continue
			}
			mp := make(map[int]int, len(colors[node]))
			for idx, color := range colors[node] {
				mp[color] = idx
			}
			colorIndex[node] = mp
			compLeader[node] = make([]int, len(colors[node]))
		}

		nodePos := make([]map[int]int, k+1)
		colorDSU := make([]*dsu, k+1)
		for color := 1; color <= k; color++ {
			nodes := nodesOfColor[color]
			if len(nodes) == 0 {
				continue
			}
			mp := make(map[int]int, len(nodes))
			for idx, node := range nodes {
				mp[node] = idx
			}
			nodePos[color] = mp
			colorDSU[color] = newDSU(len(nodes))
		}

		for _, e := range edges {
			u := e[0]
			v := e[1]
			if len(colors[u]) > len(colors[v]) {
				u, v = v, u
			}
			if len(colors[u]) == 0 || len(colors[v]) == 0 {
				continue
			}
			cmap := colorIndex[v]
			for _, color := range colors[u] {
				if _, ok := cmap[color]; !ok {
					continue
				}
				d := colorDSU[color]
				if d == nil {
					continue
				}
				mp := nodePos[color]
				d.union(mp[u], mp[v])
			}
		}

		for color := 1; color <= k; color++ {
			nodes := nodesOfColor[color]
			if len(nodes) == 0 {
				continue
			}
			d := colorDSU[color]
			if d == nil {
				continue
			}
			for idx, node := range nodes {
				leader := d.find(idx)
				pos := colorIndex[node][color]
				compLeader[node][pos] = leader
			}
		}

		threshold := int(math.Sqrt(float64(s)))
		if threshold < 1 {
			threshold = 1
		}

		heavyIdx := make([]int, n+1)
		for i := range heavyIdx {
			heavyIdx[i] = -1
		}
		heavyNodes := make([]int, 0)
		for node := 1; node <= n; node++ {
			if len(colors[node]) > threshold {
				heavyIdx[node] = len(heavyNodes)
				heavyNodes = append(heavyNodes, node)
			}
		}

		heavyCnt := len(heavyNodes)
		var heavyAns [][]int
		if heavyCnt > 0 {
			heavyAns = make([][]int, heavyCnt)
			for i := 0; i < heavyCnt; i++ {
				heavyAns[i] = make([]int, heavyCnt)
			}
			keyMap := make(map[uint64][]int, len(heavyNodes)*threshold)
			for idx, node := range heavyNodes {
				for pos, color := range colors[node] {
					leader := compLeader[node][pos]
					key := (uint64(color) << 32) | uint64(leader)
					keyMap[key] = append(keyMap[key], idx)
				}
			}
			for _, list := range keyMap {
				for i := 0; i < len(list); i++ {
					hi := list[i]
					for j := i + 1; j < len(list); j++ {
						hj := list[j]
						heavyAns[hi][hj]++
						heavyAns[hj][hi]++
					}
				}
			}
		}

		results := make([]int, q)
		for i := 0; i < q; i++ {
			u := fs.nextInt()
			v := fs.nextInt()
			if u == v {
				results[i] = len(colors[u])
				continue
			}
			iu := heavyIdx[u]
			iv := heavyIdx[v]
			if iu != -1 && iv != -1 {
				results[i] = heavyAns[iu][iv]
				continue
			}
			if len(colors[u]) > len(colors[v]) {
				u, v = v, u
			}
			ans := 0
			if len(colors[u]) > 0 && len(colors[v]) > 0 {
				cmap := colorIndex[v]
				compV := compLeader[v]
				compU := compLeader[u]
				for idx, color := range colors[u] {
					if otherIdx, ok := cmap[color]; ok {
						if compU[idx] == compV[otherIdx] {
							ans++
						}
					}
				}
			}
			results[i] = ans
		}

		for i, val := range results {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, val)
		}
		fmt.Fprintln(out)
	}
}
