package main

import (
	"bufio"
	"fmt"
	"os"
)

// shuffleInverse computes the preimage of the balanced shuffle operation.
func shuffleInverse(s string) string {
	n := len(s)
	type entry struct{ pos, open int }
	memo := make(map[[2]int]entry)

	var dfs func(int, int) bool
	dfs = func(pos, need int) bool {
		key := [2]int{pos, need}
		if e, ok := memo[key]; ok {
			return e.pos != -1
		}
		if pos == n {
			if need == 0 {
				memo[key] = entry{pos, 0}
				return true
			}
			memo[key] = entry{-1, 0}
			return false
		}
		openCnt, closeCnt := 0, 0
		best := entry{-1, 0}
		for i := pos; i < n; i++ {
			if s[i] == '(' {
				openCnt++
			} else {
				closeCnt++
				if closeCnt > need {
					break
				}
			}
			if closeCnt == need {
				if dfs(i+1, openCnt) {
					best = entry{i + 1, openCnt}
				}
			}
		}
		memo[key] = best
		return best.pos != -1
	}

	if !dfs(0, 0) {
		return ""
	}

	// reconstruct groups in reverse order
	pos, need := 0, 0
	var groups [][]byte
	for pos < n {
		e := memo[[2]int{pos, need}]
		groups = append(groups, []byte(s[pos:e.pos]))
		pos = e.pos
		need = e.open
	}

	// reverse each group to get original buckets
	buckets := make([][]byte, len(groups))
	for i, g := range groups {
		b := make([]byte, len(g))
		for j := range g {
			b[j] = g[len(g)-1-j]
		}
		buckets[i] = b
	}

	// reconstruct the original sequence
	res := make([]byte, n)
	idx := make([]int, len(buckets))
	bal := 0
	for i := 0; i < n; i++ {
		ch := buckets[bal][idx[bal]]
		idx[bal]++
		res[i] = ch
		if ch == '(' {
			bal++
		} else {
			bal--
		}
	}
	return string(res)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}
	fmt.Fprintln(writer, shuffleInverse(s))
}
