package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the problem described in problemC.txt.
// We build an Aho-Corasick automaton for all antibody patterns over
// digits 0 and 1. For each gene and pair of automaton states we keep
// the minimal length of a sequence produced from the gene that moves
// the automaton from the first state to the second one without ever
// visiting a forbidden state (where some pattern ends). Genes 0 and 1
// simply advance the automaton by one digit. For other genes we
// repeatedly relax these lengths using their production rules until a
// fixed point is reached. If no undetectable sequence exists, we print
// "YES"; otherwise we print "NO" together with the minimal length.

type acNode struct {
	next [2]int
	fail int
	out  bool
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var G, N, M int
	if _, err := fmt.Fscan(in, &G, &N, &M); err != nil {
		return
	}
	rules := make([][][]int, G)
	for i := 0; i < N; i++ {
		var a, k int
		fmt.Fscan(in, &a, &k)
		r := make([]int, k)
		for j := 0; j < k; j++ {
			fmt.Fscan(in, &r[j])
		}
		rules[a] = append(rules[a], r)
	}
	patterns := make([][]int, M)
	for i := 0; i < M; i++ {
		var l int
		fmt.Fscan(in, &l)
		p := make([]int, l)
		for j := 0; j < l; j++ {
			fmt.Fscan(in, &p[j])
		}
		patterns[i] = p
	}

	// Build Aho-Corasick automaton
	nodes := []acNode{{}}
	for _, p := range patterns {
		v := 0
		for _, d := range p {
			if nodes[v].next[d] == 0 {
				nodes = append(nodes, acNode{})
				nodes[v].next[d] = len(nodes) - 1
			}
			v = nodes[v].next[d]
		}
		nodes[v].out = true
	}
	queue := make([]int, 0)
	for c := 0; c < 2; c++ {
		v := nodes[0].next[c]
		if v != 0 {
			queue = append(queue, v)
		}
	}
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		f := nodes[v].fail
		nodes[v].out = nodes[v].out || nodes[f].out
		for c := 0; c < 2; c++ {
			u := nodes[v].next[c]
			if u != 0 {
				nodes[u].fail = nodes[f].next[c]
				queue = append(queue, u)
			} else {
				nodes[v].next[c] = nodes[f].next[c]
			}
		}
	}
	for c := 0; c < 2; c++ {
		if nodes[0].next[c] == 0 {
			nodes[0].next[c] = 0
		}
	}

	m := len(nodes)
	const INF int64 = 1 << 60
	dp := make([][][]int64, G)
	for i := 0; i < G; i++ {
		dp[i] = make([][]int64, m)
		for s := 0; s < m; s++ {
			dp[i][s] = make([]int64, m)
			for t := 0; t < m; t++ {
				dp[i][s][t] = INF
			}
		}
	}

	// digits 0 and 1
	for s := 0; s < m; s++ {
		t := nodes[s].next[0]
		if !nodes[t].out {
			dp[0][s][t] = 1
		}
		t = nodes[s].next[1]
		if !nodes[t].out {
			dp[1][s][t] = 1
		}
	}

	changed := true
	for changed {
		changed = false
		for a := 2; a < G; a++ {
			for s := 0; s < m; s++ {
				for _, rule := range rules[a] {
					tmp := make([]int64, m)
					for i := 0; i < m; i++ {
						tmp[i] = INF
					}
					tmp[s] = 0
					for _, sym := range rule {
						ntmp := make([]int64, m)
						for i := 0; i < m; i++ {
							ntmp[i] = INF
						}
						for st := 0; st < m; st++ {
							if tmp[st] == INF {
								continue
							}
							for t := 0; t < m; t++ {
								val := dp[sym][st][t]
								if val == INF {
									continue
								}
								tot := tmp[st] + val
								if tot < ntmp[t] {
									ntmp[t] = tot
								}
							}
						}
						tmp = ntmp
					}
					for t := 0; t < m; t++ {
						if tmp[t] < dp[a][s][t] {
							dp[a][s][t] = tmp[t]
							changed = true
						}
					}
				}
			}
		}
	}

	for a := 2; a < G; a++ {
		best := INF
		for t := 0; t < m; t++ {
			if dp[a][0][t] < best {
				best = dp[a][0][t]
			}
		}
		if best == INF {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO", best)
		}
	}
}
