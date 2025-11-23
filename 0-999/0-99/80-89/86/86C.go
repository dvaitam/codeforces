package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD = 1000000009

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, m int
	fmt.Fscan(reader, &n, &m)

	patterns := make([]string, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &patterns[i])
	}

	type Node struct {
		next [4]int
		fail int
		maxL int
	}

	trie := make([]Node, 110)
	nodesCount := 1

	toIdx := func(c byte) int {
		switch c {
		case 'A': return 0
		case 'C': return 1
		case 'G': return 2
		case 'T': return 3
		}
		return 0
	}

	for _, s := range patterns {
		u := 0
		for i := 0; i < len(s); i++ {
			c := toIdx(s[i])
			if trie[u].next[c] == 0 {
				trie[u].next[c] = nodesCount
				nodesCount++
			}
			u = trie[u].next[c]
		}
		if len(s) > trie[u].maxL {
			trie[u].maxL = len(s)
		}
	}

	queue := make([]int, 0, nodesCount)
	for i := 0; i < 4; i++ {
		if trie[0].next[i] != 0 {
			queue = append(queue, trie[0].next[i])
		}
	}

	head := 0
	for head < len(queue) {
		u := queue[head]
		head++
		if trie[trie[u].fail].maxL > trie[u].maxL {
			trie[u].maxL = trie[trie[u].fail].maxL
		}
		for i := 0; i < 4; i++ {
			if trie[u].next[i] != 0 {
				trie[trie[u].next[i]].fail = trie[trie[u].fail].next[i]
				queue = append(queue, trie[u].next[i])
			} else {
				trie[u].next[i] = trie[trie[u].fail].next[i]
			}
		}
	}

	dp := [2][110][1024]int{}
	dp[0][0][1023] = 1

	cur := 0
	for i := 0; i < n; i++ {
		nxt := 1 - cur
		for u := 0; u < nodesCount; u++ {
			for mask := 0; mask < 1024; mask++ {
				dp[nxt][u][mask] = 0
			}
		}

		for u := 0; u < nodesCount; u++ {
			for mask := 512; mask < 1024; mask++ {
				if dp[cur][u][mask] == 0 {
					continue
				}
				count := dp[cur][u][mask]
				for c := 0; c < 4; c++ {
					v := trie[u].next[c]
					lenMatch := trie[v].maxL
					newMask := (mask << 1) & 1023
					if lenMatch > 0 {
						newMask |= (1 << lenMatch) - 1
					}
					
					val := dp[nxt][v][newMask] + count
					if val >= MOD {
						val -= MOD
					}
					dp[nxt][v][newMask] = val
				}
			}
		}
		cur = nxt
	}

	ans := 0
	for u := 0; u < nodesCount; u++ {
		ans = (ans + dp[cur][u][1023]) % MOD
	}
	fmt.Println(ans)
}