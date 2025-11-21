package main

import (
	"bufio"
	"fmt"
	"os"
)

const inf int64 = 1<<62 - 1

type fastScanner struct {
	r *bufio.Reader
}

func newFastScanner() *fastScanner {
	return &fastScanner{r: bufio.NewReader(os.Stdin)}
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

func find(par []int, x int) int {
	for par[x] != x {
		par[x] = par[par[x]]
		x = par[x]
	}
	return x
}

type stackItem struct {
	node    int
	visited bool
}

func main() {
	fs := newFastScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	T := fs.nextInt()
	for ; T > 0; T-- {
		n := fs.nextInt()
		m := fs.nextInt()
		size := 2*n + 5

		parent := make([]int, size)
		left := make([]int, size)
		right := make([]int, size)
		label := make([]int64, size)
		bestUp := make([]int64, size)
		dpCost := make([]int64, size)
		remain := make([]int, size)
		parity := make([]int, n+1)

		for i := 1; i < size; i++ {
			parent[i] = i
			label[i] = inf
		}

		sumWeights := int64(0)
		tot := n

		for i := 0; i < m; i++ {
			u := fs.nextInt()
			v := fs.nextInt()
			w := fs.nextInt()
			weight := int64(w)
			sumWeights += weight
			if u != v {
				parity[u] ^= 1
				parity[v] ^= 1
			}

			ru := find(parent, u)
			rv := find(parent, v)
			if ru != rv {
				tot++
				left[tot] = ru
				right[tot] = rv
				label[tot] = weight
				parent[ru] = tot
				parent[rv] = tot
				parent[tot] = tot
			} else {
				if label[ru] > weight {
					label[ru] = weight
				}
			}
		}

		root := find(parent, 1)
		stack := []int{root}
		bestUp[root] = label[root]
		for len(stack) > 0 {
			node := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if node <= n {
				continue
			}
			lc := left[node]
			rc := right[node]
			if lc != 0 {
				bestUp[lc] = label[lc]
				if bestUp[node] < bestUp[lc] {
					bestUp[lc] = bestUp[node]
				}
				stack = append(stack, lc)
			}
			if rc != 0 {
				bestUp[rc] = label[rc]
				if bestUp[node] < bestUp[rc] {
					bestUp[rc] = bestUp[node]
				}
				stack = append(stack, rc)
			}
		}

		post := []stackItem{{node: root, visited: false}}
		for len(post) > 0 {
			item := post[len(post)-1]
			post = post[:len(post)-1]
			node := item.node
			if node <= n {
				remain[node] = parity[node]
				dpCost[node] = 0
				continue
			}
			if !item.visited {
				post = append(post, stackItem{node: node, visited: true})
				post = append(post, stackItem{node: right[node], visited: false})
				post = append(post, stackItem{node: left[node], visited: false})
				continue
			}
			lc := left[node]
			rc := right[node]
			cntL := remain[lc]
			cntR := remain[rc]
			pairs := cntL
			if cntR < pairs {
				pairs = cntR
			}
			dpCost[node] = dpCost[lc] + dpCost[rc] + int64(pairs)*bestUp[node]
			remain[node] = cntL + cntR - 2*pairs
		}

		answer := sumWeights + dpCost[root]
		fmt.Fprintln(out, answer)
	}
}
