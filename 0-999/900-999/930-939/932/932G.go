package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 1000000007

// node represents a palindrome in the eertree
type node struct {
	next  [26]int
	link  int
	len   int
	diff  int
	slink int
}

// countPartitions counts ways to partition string t into even length palindromes
func countPartitions(t []byte) int {
	n := len(t)
	s := make([]byte, n+1)
	copy(s[1:], t)

	nodes := make([]node, n+3)
	nodes[1].len = -1
	nodes[1].link = 1
	nodes[1].slink = 1
	nodes[2].len = 0
	nodes[2].link = 1
	nodes[2].slink = 1
	last, tot := 2, 2

	dp := make([]int, n+1)
	series := make([]int, n+3)
	dp[0] = 1

	for i := 1; i <= n; i++ {
		c := int(s[i] - 'a')
		cur := last
		for s[i-nodes[cur].len-1] != s[i] {
			cur = nodes[cur].link
		}
		if nodes[cur].next[c] == 0 {
			tot++
			nodes[tot].len = nodes[cur].len + 2
			x := nodes[cur].link
			for s[i-nodes[x].len-1] != s[i] {
				x = nodes[x].link
			}
			linkNode := nodes[x].next[c]
			if linkNode == 0 {
				linkNode = 2
			}
			nodes[tot].link = linkNode
			nodes[cur].next[c] = tot
			nodes[tot].diff = nodes[tot].len - nodes[linkNode].len
			if nodes[tot].diff == nodes[linkNode].diff {
				nodes[tot].slink = nodes[linkNode].slink
			} else {
				nodes[tot].slink = linkNode
			}
		}
		last = nodes[cur].next[c]
		dp[i] = 0
		v := last
		for nodes[v].len > 0 {
			series[v] = dp[i-(nodes[nodes[v].slink].len+nodes[v].diff)]
			if nodes[v].diff == nodes[nodes[v].link].diff {
				series[v] = (series[v] + series[nodes[v].link]) % mod
			}
			if i%2 == 0 {
				dp[i] = (dp[i] + series[v]) % mod
			}
			v = nodes[v].slink
		}
	}
	return dp[n]
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var s string
	fmt.Fscan(reader, &s)
	n := len(s)
	t := make([]byte, n)
	for i := 0; i < n/2; i++ {
		t[2*i] = s[i]
		t[2*i+1] = s[n-1-i]
	}
	fmt.Println(countPartitions(t))
}
