package main

import (
	"bufio"
	"fmt"
	"os"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var s string
	fmt.Fscan(reader, &s)
	n := len(s)

	nxt := make([][]int, n+1)
	for i := range nxt {
		nxt[i] = make([]int, 26)
	}
	for c := 0; c < 26; c++ {
		nxt[n][c] = n
	}
	for i := n - 1; i >= 0; i-- {
		copy(nxt[i], nxt[i+1])
		nxt[i][int(s[i]-'a')] = i
	}

	var keyBuf []byte
	buildKey := func(pos []int, idx int) string {
		keyBuf = keyBuf[:0]
		for _, v := range pos {
			keyBuf = append(keyBuf, byte(v))
		}
		keyBuf = append(keyBuf, byte(idx))
		return string(keyBuf)
	}

	var can func(k, m int) bool
	can = func(k, m int) bool {
		pos := make([]int, k)
		memo := map[string]bool{}
		var dfs func(int) bool
		dfs = func(idx int) bool {
			if idx == m {
				return true
			}
			key := buildKey(pos, idx)
			if v, ok := memo[key]; ok {
				return v
			}
			for c := 0; c < 26; c++ {
				p := nxt[pos[0]][c]
				if p == n {
					continue
				}
				newPos := make([]int, k)
				newPos[0] = p + 1
				valid := true
				for t := 1; t < k; t++ {
					p2 := nxt[max(pos[t], newPos[t-1])][c]
					if p2 == n {
						valid = false
						break
					}
					newPos[t] = p2 + 1
				}
				if !valid {
					continue
				}
				old := make([]int, k)
				copy(old, pos)
				copy(pos, newPos)
				if dfs(idx + 1) {
					memo[key] = true
					copy(pos, old)
					return true
				}
				copy(pos, old)
			}
			memo[key] = false
			return false
		}
		return dfs(0)
	}

	ans := 0
	for L := n; L >= 2; L-- {
		found := false
		for k := 2; k <= L; k++ {
			if L%k != 0 {
				continue
			}
			m := L / k
			if can(k, m) {
				ans = L
				found = true
				break
			}
		}
		if found {
			break
		}
	}
	fmt.Println(ans)
}
