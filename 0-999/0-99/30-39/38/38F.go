package main

import (
	"bufio"
	"fmt"
	"os"
)

type result struct {
	win      bool
	own, opp int64
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	words := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &words[i])
	}
	// collect substrings and counts
	subMap := make(map[string]int)
	num := make([]int64, 0)
	for i := 0; i < n; i++ {
		w := words[i]
		seen := make(map[string]bool)
		L := len(w)
		for l := 0; l < L; l++ {
			for r := l + 1; r <= L; r++ {
				s := w[l:r]
				if !seen[s] {
					seen[s] = true
					if id, ok := subMap[s]; ok {
						num[id]++
					} else {
						id = len(num)
						subMap[s] = id
						num = append(num, 1)
					}
				}
			}
		}
	}
	m := len(num)
	subs := make([]string, m)
	for s, id := range subMap {
		subs[id] = s
	}
	// compute weights and group by length
	weight := make([]int64, m)
	maxLen := 0
	buckets := make(map[int][]int)
	for id, s := range subs {
		sum := int64(0)
		for i := 0; i < len(s); i++ {
			sum += int64(s[i] - 'a' + 1)
		}
		weight[id] = sum * num[id]
		L := len(s)
		if L > maxLen {
			maxLen = L
		}
		buckets[L] = append(buckets[L], id)
	}
	// build transitions
	nexts := make([][]int, m)
	for id, s := range subs {
		for c := byte('a'); c <= 'z'; c++ {
			t1 := string(c) + s
			if j, ok := subMap[t1]; ok {
				nexts[id] = append(nexts[id], j)
			}
			t2 := s + string(c)
			if j, ok := subMap[t2]; ok {
				nexts[id] = append(nexts[id], j)
			}
		}
	}
	// DP
	dp := make([]result, m)
	// lengths from maxLen down to 1
	for L := maxLen; L >= 1; L-- {
		for _, id := range buckets[L] {
			bestSet := false
			var best result
			// for each move
			for _, j := range nexts[id] {
				child := dp[j]
				w := weight[j]
				res := result{
					win: !child.win,
					own: w + child.opp,
					opp: child.own,
				}
				if !bestSet || better(res, best) {
					bestSet = true
					best = res
				}
			}
			if bestSet {
				dp[id] = best
			} else {
				dp[id] = result{win: false, own: 0, opp: 0}
			}
		}
	}
	// initial choice among length-1 substrings
	var initSet bool
	var initRes result
	for _, id := range buckets[1] {
		child := dp[id]
		w := weight[id]
		res := result{
			win: !child.win,
			own: w + child.opp,
			opp: child.own,
		}
		if !initSet || better(res, initRes) {
			initSet = true
			initRes = res
		}
	}
	if initRes.win {
		fmt.Println("First")
	} else {
		fmt.Println("Second")
	}
	fmt.Printf("%d %d\n", initRes.own, initRes.opp)
}

func better(a, b result) bool {
	if a.win != b.win {
		return a.win && !b.win
	}
	if a.own != b.own {
		return a.own > b.own
	}
	return a.opp < b.opp
}
