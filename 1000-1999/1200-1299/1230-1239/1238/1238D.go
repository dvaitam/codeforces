package main

import (
	"bufio"
	"fmt"
	"os"
)

type Fenwick struct {
	n   int
	bit []int
}

func NewFenwick(n int) *Fenwick { return &Fenwick{n, make([]int, n+2)} }
func (f *Fenwick) Add(idx, val int) {
	for idx <= f.n {
		f.bit[idx] += val
		idx += idx & -idx
	}
}
func (f *Fenwick) Sum(idx int) int {
	res := 0
	for idx > 0 {
		res += f.bit[idx]
		idx -= idx & -idx
	}
	return res
}

func countGood(s string) int64 {
	n := len(s)
	next := make([]int, n)
	prev := make([]int, n)
	last := map[byte]int{'A': -1, 'B': -1}
	for i := 0; i < n; i++ { // prev same index
		prev[i] = last[s[i]]
		last[s[i]] = i
	}
	last['A'] = n
	last['B'] = n
	for i := n - 1; i >= 0; i-- { // next same index
		next[i] = last[s[i]]
		last[s[i]] = i
	}
	buckets := make([][]int, n+1) // buckets[r] -> list of l with next[l]=r
	for l := 0; l < n; l++ {
		if next[l] < n {
			buckets[next[l]] = append(buckets[next[l]], l+1) // convert to 1-index
		}
	}
	ft := NewFenwick(n)
	var ans int64
	for r := 0; r < n; r++ { // r is 0-index
		// add all l whose next[l] == r
		for _, idx := range buckets[r] {
			ft.Add(idx, 1)
		}
		if prev[r] >= 0 {
			ans += int64(ft.Sum(prev[r] + 1))
		}
	}
	return ans
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	var s string
	if _, err := fmt.Fscan(in, &n, &s); err != nil {
		return
	}
	fmt.Println(countGood(s))
}
