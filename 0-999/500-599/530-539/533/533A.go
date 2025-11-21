package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const inf int64 = int64(^uint64(0) >> 2)

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
		c, err = fs.r.ReadByte()
		if err != nil {
			return 0
		}
	}
	if c == '-' {
		sign = -1
		c, err = fs.r.ReadByte()
		if err != nil {
			return 0
		}
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int(c-'0')
		c, err = fs.r.ReadByte()
		if err != nil {
			break
		}
	}
	return val * sign
}

func (fs *fastScanner) nextInt64() int64 {
	sign := int64(1)
	var val int64
	c, err := fs.r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		c, err = fs.r.ReadByte()
		if err != nil {
			return 0
		}
	}
	if c == '-' {
		sign = -1
		c, err = fs.r.ReadByte()
		if err != nil {
			return 0
		}
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int64(c-'0')
		c, err = fs.r.ReadByte()
		if err != nil {
			break
		}
	}
	return val * sign
}

func countGE(desc []int64, val int64) int {
	return sort.Search(len(desc), func(i int) bool {
		return desc[i] < val
	})
}

type state struct {
	bestVal   int64
	bestCount int
	bestNode  int
	secondVal int64
}

type requirement struct {
	cnt       int
	threshold int64
}

func main() {
	fs := newFastScanner()
	n := fs.nextInt()
	if n == 0 {
		fmt.Println(0)
		return
	}

	h := make([]int64, n)
	for i := 0; i < n; i++ {
		h[i] = fs.nextInt64()
	}

	adj := make([][]int, n)
	for i := 0; i < n-1; i++ {
		a := fs.nextInt() - 1
		b := fs.nextInt() - 1
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
	}

	k := fs.nextInt()
	miners := make([]int64, k)
	for i := 0; i < k; i++ {
		miners[i] = fs.nextInt64()
	}

	parent := make([]int, n)
	for i := range parent {
		parent[i] = -1
	}
	order := make([]int, 0, n)
	stack := []int{0}
	parent[0] = -2
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, v)
		for _, to := range adj[v] {
			if to == parent[v] {
				continue
			}
			parent[to] = v
			stack = append(stack, to)
		}
	}

	bestVal := make([]int64, n)
	bestCount := make([]int, n)
	bestNode := make([]int, n)
	secondVal := make([]int64, n)
	uniqueOwner := make([]int, n)
	for i := range uniqueOwner {
		uniqueOwner[i] = -1
	}
	pathMin := make([]int64, n)

	for _, u := range order {
		var cur state
		if parent[u] == -2 {
			cur = state{bestVal: inf, secondVal: inf, bestNode: -1}
		} else {
			p := parent[u]
			cur = state{
				bestVal:   bestVal[p],
				bestCount: bestCount[p],
				bestNode:  bestNode[p],
				secondVal: secondVal[p],
			}
		}

		val := h[u]
		if val < cur.bestVal {
			cur.secondVal = cur.bestVal
			cur.bestVal = val
			cur.bestCount = 1
			cur.bestNode = u
		} else if val == cur.bestVal {
			cur.bestCount++
			cur.bestNode = -1
		} else if val < cur.secondVal {
			cur.secondVal = val
		}

		bestVal[u] = cur.bestVal
		bestCount[u] = cur.bestCount
		bestNode[u] = cur.bestNode
		secondVal[u] = cur.secondVal
		pathMin[u] = cur.bestVal
		if cur.bestCount == 1 {
			uniqueOwner[u] = cur.bestNode
		}
	}

	// For each node that acts as the unique minimum on some root-to-node path,
	// collect the second minimum values of those paths (the maximal height they
	// can reach once the bottleneck is lifted).
	groups := make([][]int64, n)
	for u := 0; u < n; u++ {
		if owner := uniqueOwner[u]; owner != -1 {
			groups[owner] = append(groups[owner], secondVal[u])
		}
	}

	nodesDesc := make([]int64, n)
	copy(nodesDesc, pathMin)
	sort.Slice(nodesDesc, func(i, j int) bool { return nodesDesc[i] > nodesDesc[j] })

	minersDesc := make([]int64, k)
	copy(minersDesc, miners)
	sort.Slice(minersDesc, func(i, j int) bool { return minersDesc[i] > minersDesc[j] })

	uniqueTs := make([]int64, 0)
	reqCounts := make([]int, 0)
	i := 0
	for i < k {
		t := minersDesc[i]
		j := i
		for j < k && minersDesc[j] == t {
			j++
		}
		uniqueTs = append(uniqueTs, t)
		reqCounts = append(reqCounts, j)
		i = j
	}

	D := make([]int, len(uniqueTs))
	for idx, t := range uniqueTs {
		nodesCnt := countGE(nodesDesc, t)
		D[idx] = reqCounts[idx] - nodesCnt
	}

	// reqPairs stores prefix-max deficits: to satisfy all miners we need
	// at least req.cnt nodes whose new minima reach req.threshold.
	reqPairs := make([]requirement, 0)
	required := 0
	for idx, t := range uniqueTs {
		if D[idx] > required {
			required = D[idx]
			if required > 0 {
				reqPairs = append(reqPairs, requirement{cnt: required, threshold: t})
			}
		}
	}

	if len(reqPairs) == 0 {
		fmt.Println(0)
		return
	}

	tMax := reqPairs[0].threshold
	maxNeed := reqPairs[len(reqPairs)-1].cnt

	for v := 0; v < n; v++ {
		if len(groups[v]) == 0 {
			continue
		}
		sort.Slice(groups[v], func(i, j int) bool { return groups[v][i] > groups[v][j] })
		if len(groups[v]) > maxNeed {
			groups[v] = groups[v][:maxNeed]
		}
	}

	answer := int64(-1)
	for v := 0; v < n; v++ {
		if h[v] >= tMax {
			continue
		}
		arr := groups[v]
		if len(arr) < maxNeed {
			continue
		}
		ok := true
		for _, req := range reqPairs {
			if arr[req.cnt-1] < req.threshold {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}
		delta := tMax - h[v]
		if answer == -1 || delta < answer {
			answer = delta
		}
	}

	if answer == -1 {
		fmt.Println(-1)
	} else {
		fmt.Println(answer)
	}
}
