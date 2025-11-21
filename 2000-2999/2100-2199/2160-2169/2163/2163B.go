package main

import (
	"bufio"
	"fmt"
	"os"
)

type fastScanner struct {
	r *bufio.Reader
}

func newFastScanner() *fastScanner {
	return &fastScanner{r: bufio.NewReader(os.Stdin)}
}

func (fs *fastScanner) nextInt() int {
	sign, val := 1, 0
	c, _ := fs.r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		c, _ = fs.r.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, _ = fs.r.ReadByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int(c-'0')
		c, _ = fs.r.ReadByte()
	}
	return sign * val
}

func (fs *fastScanner) nextString() string {
	c, _ := fs.r.ReadByte()
	for c <= ' ' {
		c, _ = fs.r.ReadByte()
	}
	buf := make([]byte, 0, 16)
	for c > ' ' {
		buf = append(buf, c)
		c, _ = fs.r.ReadByte()
	}
	return string(buf)
}

type pair struct {
	l, r int
}

func computeNearest(p []int, less bool, left bool) []int {
	n := len(p)
	res := make([]int, n)
	for i := range res {
		res[i] = -1
	}
	stack := make([]int, 0, n)
	if left {
		for i := 0; i < n; i++ {
			val := p[i]
			for len(stack) > 0 {
				top := stack[len(stack)-1]
				if less {
					if p[top] < val {
						break
					}
				} else {
					if p[top] > val {
						break
					}
				}
				stack = stack[:len(stack)-1]
			}
			if len(stack) > 0 {
				res[i] = stack[len(stack)-1]
			}
			stack = append(stack, i)
		}
	} else {
		for i := n - 1; i >= 0; i-- {
			val := p[i]
			for len(stack) > 0 {
				top := stack[len(stack)-1]
				if less {
					if p[top] < val {
						break
					}
				} else {
					if p[top] > val {
						break
					}
				}
				stack = stack[:len(stack)-1]
			}
			if len(stack) > 0 {
				res[i] = stack[len(stack)-1]
			}
			stack = append(stack, i)
		}
	}
	return res
}

func main() {
	fs := newFastScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	t := fs.nextInt()
	for ; t > 0; t-- {
		n := fs.nextInt()
		p := make([]int, n)
		for i := 0; i < n; i++ {
			p[i] = fs.nextInt()
		}
		x := fs.nextString()

		needPos := make([]int, 0)
		for i := 0; i < n; i++ {
			if x[i] == '1' {
				needPos = append(needPos, i)
			}
		}

		if len(needPos) == 0 {
			fmt.Fprintln(out, 0)
			continue
		}

		prevLess := computeNearest(p, true, true)
		prevGreater := computeNearest(p, false, true)
		nextLess := computeNearest(p, true, false)
		nextGreater := computeNearest(p, false, false)

		m := len(needPos)
		candidates := make([][]int64, m)
		opInfo := make(map[int64]pair)
		failed := false
		mult := int64(n + 1)
		for idx, pos := range needPos {
			ops := make([]int64, 0, 2)
			if prevLess[pos] != -1 && nextGreater[pos] != -1 {
				key := int64(prevLess[pos])*mult + int64(nextGreater[pos])
				if _, ok := opInfo[key]; !ok {
					opInfo[key] = pair{prevLess[pos], nextGreater[pos]}
				}
				ops = append(ops, key)
			}
			if prevGreater[pos] != -1 && nextLess[pos] != -1 {
				key := int64(prevGreater[pos])*mult + int64(nextLess[pos])
				if _, ok := opInfo[key]; !ok {
					opInfo[key] = pair{prevGreater[pos], nextLess[pos]}
				}
				// avoid duplicate entry when both keys same
				if len(ops) == 0 || ops[len(ops)-1] != key {
					ops = append(ops, key)
				}
			}
			if len(ops) == 0 {
				fmt.Fprintln(out, -1)
				failed = true
				break
			}
			candidates[idx] = ops
		}
		if failed {
			continue
		}

		covered := make([]bool, m)
		opCoverage := make(map[int64][]int)
		curOps := make([]pair, 0, 5)
		var answer []pair
		coveredCount := 0

		var getCoverage func(int64) []int
		getCoverage = func(key int64) []int {
			if v, ok := opCoverage[key]; ok {
				return v
			}
			op := opInfo[key]
			l, r := op.l, op.r
			lo, hi := p[l], p[r]
			if lo > hi {
				lo, hi = hi, lo
			}
			arr := make([]int, 0)
			for idx, pos := range needPos {
				if l < pos && pos < r {
					val := p[pos]
					if lo < val && val < hi {
						arr = append(arr, idx)
					}
				}
			}
			opCoverage[key] = arr
			return arr
		}

		var dfs func() bool
		dfs = func() bool {
			if coveredCount == m {
				answer = append([]pair(nil), curOps...)
				return true
			}
			if len(curOps) == 5 {
				return false
			}
			u := -1
			for i := 0; i < m; i++ {
				if !covered[i] {
					u = i
					break
				}
			}
			if u == -1 {
				answer = append([]pair(nil), curOps...)
				return true
			}
			for _, key := range candidates[u] {
				cov := getCoverage(key)
				newly := make([]int, 0, len(cov))
				for _, idx := range cov {
					if !covered[idx] {
						covered[idx] = true
						newly = append(newly, idx)
					}
				}
				if len(newly) == 0 {
					continue
				}
				coveredCount += len(newly)
				curOps = append(curOps, opInfo[key])
				if dfs() {
					return true
				}
				curOps = curOps[:len(curOps)-1]
				coveredCount -= len(newly)
				for _, idx := range newly {
					covered[idx] = false
				}
			}
			return false
		}

		if !dfs() {
			fmt.Fprintln(out, -1)
			continue
		}

		fmt.Fprintln(out, len(answer))
		for _, op := range answer {
			fmt.Fprintf(out, "%d %d\n", op.l+1, op.r+1)
		}
	}
}
