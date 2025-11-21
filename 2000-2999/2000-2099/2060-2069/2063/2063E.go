package main

import (
	"bufio"
	"fmt"
	"os"
)

type treap struct {
	key         int
	val, sum    int64
	pr          uint32
	left, right *treap
}

var rng uint32 = 123456789

func nextRand() uint32 {
	rng ^= rng << 13
	rng ^= rng >> 17
	rng ^= rng << 5
	return rng
}

func treapSum(t *treap) int64 {
	if t == nil {
		return 0
	}
	return t.sum
}

func treapUpdate(t *treap) {
	if t != nil {
		t.sum = t.val + treapSum(t.left) + treapSum(t.right)
	}
}

func treapRotateRight(t *treap) *treap {
	l := t.left
	t.left = l.right
	l.right = t
	treapUpdate(t)
	treapUpdate(l)
	return l
}

func treapRotateLeft(t *treap) *treap {
	r := t.right
	t.right = r.left
	r.left = t
	treapUpdate(t)
	treapUpdate(r)
	return r
}

func treapInsert(t *treap, key int, val int64) *treap {
	if t == nil {
		return &treap{key: key, val: val, sum: val, pr: nextRand()}
	}
	if key == t.key {
		t.val += val
	} else if key < t.key {
		t.left = treapInsert(t.left, key, val)
		if t.left.pr < t.pr {
			t = treapRotateRight(t)
		}
	} else {
		t.right = treapInsert(t.right, key, val)
		if t.right.pr < t.pr {
			t = treapRotateLeft(t)
		}
	}
	treapUpdate(t)
	return t
}

func treapQueryGE(t *treap, key int) int64 {
	if t == nil {
		return 0
	}
	if t.key < key {
		return treapQueryGE(t.right, key)
	}
	return t.val + treapSum(t.right) + treapQueryGE(t.left, key)
}

type fastReader struct {
	r *bufio.Reader
}

func newFastReader() *fastReader {
	return &fastReader{r: bufio.NewReader(os.Stdin)}
}

func (fr *fastReader) nextInt() int {
	sign := 1
	val := 0
	c, _ := fr.r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		c, _ = fr.r.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, _ = fr.r.ReadByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int(c-'0')
		c, _ = fr.r.ReadByte()
	}
	return sign * val
}

func main() {
	fr := newFastReader()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	t := fr.nextInt()
	for ; t > 0; t-- {
		n := fr.nextInt()
		adj := make([][]int, n+1)
		for i := 0; i < n-1; i++ {
			u := fr.nextInt()
			v := fr.nextInt()
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}

		parent := make([]int, n+1)
		depth := make([]int, n+1)
		order := make([]int, 0, n)
		stack := []int{1}
		parent[1] = 0
		for len(stack) > 0 {
			u := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			order = append(order, u)
			for _, v := range adj[u] {
				if v == parent[u] {
					continue
				}
				parent[v] = u
				depth[v] = depth[u] + 1
				stack = append(stack, v)
			}
		}

		dp := make([]map[int]int64, n+1)
		subSize := make([]int64, n+1)
		var answer int64

		for idx := len(order) - 1; idx >= 0; idx-- {
			u := order[idx]
			children := make([]int, 0)
			for _, v := range adj[u] {
				if parent[v] == u {
					children = append(children, v)
				}
			}

			type childInfo struct {
				mp   map[int]int64
				size int64
			}
			childMaps := make([]childInfo, 0, len(children))
			for _, v := range children {
				childMaps = append(childMaps, childInfo{dp[v], subSize[v]})
			}

			if len(childMaps) >= 2 {
				var tr *treap
				var totalPairs, totalSize, sumMin int64
				for _, ch := range childMaps {
					size := ch.size
					totalPairs += size * totalSize
					totalSize += size
					for depthKey, val := range ch.mp {
						other := treapQueryGE(tr, depthKey)
						sumMin += int64(depthKey-depth[u]) * val * other
					}
					for depthKey, val := range ch.mp {
						tr = treapInsert(tr, depthKey, val)
					}
				}
				answer += 2*sumMin - totalPairs
			}

			var base map[int]int64
			if len(childMaps) > 0 {
				idxMax := 0
				for i := 1; i < len(childMaps); i++ {
					if len(childMaps[i].mp) > len(childMaps[idxMax].mp) {
						idxMax = i
					}
				}
				base = childMaps[idxMax].mp
				for i := 0; i < len(childMaps); i++ {
					if i == idxMax {
						continue
					}
					for depthKey, val := range childMaps[i].mp {
						base[depthKey] += val
					}
				}
			} else {
				base = make(map[int]int64)
			}
			base[depth[u]]++
			dp[u] = base

			subSize[u] = 1
			for _, v := range children {
				subSize[u] += subSize[v]
			}
		}

		fmt.Fprintln(out, answer)
	}
}
