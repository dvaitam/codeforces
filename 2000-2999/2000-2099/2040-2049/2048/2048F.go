package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	maxK = 60                  // 2^60 > 1e18, enough operations along any path
	inf  = uint64(1<<62) - 1e6 // large sentinel to avoid overflow
)

func ceilDiv(a, b uint64) uint64 {
	if a >= inf {
		return inf
	}
	return (a + b - 1) / b
}

func value(arr []uint64, idx int) uint64 {
	if idx < len(arr) {
		return arr[idx]
	}
	return 1
}

func ensureTailOne(arr []uint64) []uint64 {
	if len(arr) == 0 || arr[len(arr)-1] != 1 {
		arr = append(arr, 1)
	}
	return arr
}

func trimTrailingOnes(arr []uint64) []uint64 {
	for len(arr) > 1 && arr[len(arr)-1] == 1 && arr[len(arr)-2] == 1 {
		arr = arr[:len(arr)-1]
	}
	return arr
}

func effectiveLen(arr []uint64) int {
	if len(arr) == 0 {
		return 0
	}
	if arr[len(arr)-1] == 1 {
		return len(arr) - 1
	}
	return len(arr)
}

// combine children needs: res[m] = min_{i+j=m} max(a[i], b[j])
// a, b are non-increasing; res inherits non-increasing property.
func combine(a, b []uint64) []uint64 {
	effA := effectiveLen(a)
	effB := effectiveLen(b)
	limit := effA + effB
	if limit > maxK {
		limit = maxK
	}
	res := make([]uint64, limit+1)
	p := 0
	for m := 0; m <= limit; m++ {
		if p > m {
			p = m
		}
		best := value(a, p)
		vb := value(b, m-p)
		if vb > best {
			best = vb
		}
		for p < m {
			cand := value(a, p+1)
			v2 := value(b, m-p-1)
			if v2 > cand {
				cand = v2
			}
			if cand <= best {
				p++
				best = cand
			} else {
				break
			}
		}
		res[m] = best
	}
	res = ensureTailOne(res)
	res = trimTrailingOnes(res)
	return res
}

func buildPows(base uint64) [maxK + 1]uint64 {
	var p [maxK + 1]uint64
	p[0] = 1
	for i := 1; i <= maxK; i++ {
		if p[i-1] >= inf/base {
			p[i] = inf
		} else {
			p[i] = p[i-1] * base
			if p[i] > inf {
				p[i] = inf
			}
		}
	}
	return p
}

// compute need array for node with value aVal, base bVal and combined child need.
func computeNeed(aVal, bVal uint64, child []uint64) []uint64 {
	pows := buildPows(bVal)
	need := make([]uint64, maxK+1)
	for i := range need {
		need[i] = inf
	}

	maxStart := aVal
	c0 := value(child, 0)
	if c0 > maxStart {
		maxStart = c0
	}

	childEff := effectiveLen(child)
	for t := 0; t <= maxK; t++ {
		factor := pows[t]
		if factor >= maxStart {
			for k := t; k <= maxK; k++ {
				need[k] = 1
			}
			break
		}

		limitRem := maxK - t
		if childEff < limitRem {
			limitRem = childEff
		}
		for rem := 0; rem <= limitRem; rem++ {
			cur := value(child, rem)
			if aVal > cur {
				cur = aVal
			}
			req := ceilDiv(cur, factor)
			k := t + rem
			if req < need[k] {
				need[k] = req
			}
			if req == 1 {
				for kk := k; kk <= maxK && need[kk] > 1; kk++ {
					need[kk] = 1
				}
				break
			}
		}
	}

	for k := 1; k <= maxK; k++ {
		if need[k] > need[k-1] {
			need[k] = need[k-1]
		}
	}
	need = ensureTailOne(need)
	need = trimTrailingOnes(need)
	return need
}

// build min-cartesian tree for b array; returns root index and children arrays.
func buildCartesian(b []uint64) (int, []int, []int) {
	n := len(b)
	parent := make([]int, n)
	for i := range parent {
		parent[i] = -1
	}
	stack := make([]int, 0, n)
	for i := 0; i < n; i++ {
		last := -1
		for len(stack) > 0 && b[stack[len(stack)-1]] > b[i] {
			last = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
		}
		if len(stack) > 0 {
			parent[i] = stack[len(stack)-1]
		}
		if last != -1 {
			parent[last] = i
		}
		stack = append(stack, i)
	}

	left := make([]int, n)
	right := make([]int, n)
	for i := 0; i < n; i++ {
		left[i], right[i] = -1, -1
	}
	root := -1
	for i, p := range parent {
		if p == -1 {
			root = i
			continue
		}
		if i < p {
			left[p] = i
		} else {
			right[p] = i
		}
	}
	return root, left, right
}

type fastScanner struct {
	r *bufio.Reader
}

func newScanner() *fastScanner {
	return &fastScanner{r: bufio.NewReader(os.Stdin)}
}

func (fs *fastScanner) next() (uint64, bool) {
	var val uint64
	c, err := fs.r.ReadByte()
	for err == nil && (c == ' ' || c == '\n' || c == '\r' || c == '\t') {
		c, err = fs.r.ReadByte()
	}
	if err != nil {
		return 0, false
	}
	for err == nil && c >= '0' && c <= '9' {
		val = val*10 + uint64(c-'0')
		c, err = fs.r.ReadByte()
	}
	if err == nil {
		fs.r.UnreadByte()
	}
	return val, true
}

func main() {
	in := newScanner()
	tVal, ok := in.next()
	if !ok {
		return
	}
	t := int(tVal)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	for ; t > 0; t-- {
		nVal, _ := in.next()
		n := int(nVal)
		a := make([]uint64, n)
		b := make([]uint64, n)
		for i := 0; i < n; i++ {
			x, _ := in.next()
			a[i] = x
		}
		for i := 0; i < n; i++ {
			x, _ := in.next()
			b[i] = x
		}

		root, left, right := buildCartesian(b)
		order := make([]int, 0, n)
		stack := []int{root}
		for len(stack) > 0 {
			v := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			order = append(order, v)
			if left[v] != -1 {
				stack = append(stack, left[v])
			}
			if right[v] != -1 {
				stack = append(stack, right[v])
			}
		}

		needs := make([][]uint64, n)
		for i := len(order) - 1; i >= 0; i-- {
			v := order[i]
			ln := []uint64{1}
			if left[v] != -1 {
				ln = needs[left[v]]
			}
			rn := []uint64{1}
			if right[v] != -1 {
				rn = needs[right[v]]
			}
			child := combine(ln, rn)
			needs[v] = computeNeed(a[v], b[v], child)
		}

		ans := 0
		rootNeed := needs[root]
		for ans < len(rootNeed) && rootNeed[ans] > 1 {
			ans++
		}
		fmt.Fprintln(writer, ans)
	}
}
