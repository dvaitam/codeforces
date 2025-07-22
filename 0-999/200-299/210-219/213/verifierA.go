package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type caseA struct {
	n   int
	c   []int
	pre [][]int
}

func genCaseA(rng *rand.Rand) caseA {
	n := rng.Intn(8) + 2 // 2..9 to keep runtime small
	c := make([]int, n+1)
	for i := 1; i <= n; i++ {
		c[i] = rng.Intn(3) + 1
	}
	pre := make([][]int, n+1)
	for i := 1; i <= n; i++ {
		k := rng.Intn(i) // prerequisites only from <i
		used := map[int]bool{}
		for j := 0; j < k; j++ {
			for {
				u := rng.Intn(i)
				if u == 0 || used[u] {
					continue
				}
				used[u] = true
				pre[i] = append(pre[i], u)
				break
			}
		}
	}
	return caseA{n, c, pre}
}

func solveA(tc caseA) int {
	n := tc.n
	c := tc.c
	adj := make([][]int, n+1)
	indeg0 := make([]int, n+1)
	for i := 1; i <= n; i++ {
		for _, u := range tc.pre[i] {
			adj[u] = append(adj[u], i)
			indeg0[i]++
		}
	}
	moveCost := [4][4]int{{}, {0, 0, 1, 2}, {0, 2, 0, 1}, {0, 1, 2, 0}}
	const INF = int(1 << 30)
	res := INF
	for start := 1; start <= 3; start++ {
		indeg := make([]int, n+1)
		copy(indeg, indeg0)
		avail := make([][]int, 4)
		for i := 1; i <= n; i++ {
			if indeg[i] == 0 {
				avail[c[i]] = append(avail[c[i]], i)
			}
		}
		curr := start
		cost := 0
		done := 0
		for done < n {
			if len(avail[curr]) > 0 {
				u := avail[curr][0]
				avail[curr] = avail[curr][1:]
				cost += 1
				done++
				for _, v := range adj[u] {
					indeg[v]--
					if indeg[v] == 0 {
						avail[c[v]] = append(avail[c[v]], v)
					}
				}
			} else {
				bestM := 0
				best := INF
				for m := 1; m <= 3; m++ {
					if len(avail[m]) > 0 && moveCost[curr][m] < best {
						best = moveCost[curr][m]
						bestM = m
					}
				}
				if bestM == 0 {
					break
				}
				cost += best
				curr = bestM
			}
		}
		if done == n && cost < res {
			res = cost
		}
	}
	return res
}

func runA(bin string, tc caseA) error {
	var sb strings.Builder
	fmt.Fprintln(&sb, tc.n)
	for i := 1; i <= tc.n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		fmt.Fprint(&sb, tc.c[i])
	}
	sb.WriteByte('\n')
	for i := 1; i <= tc.n; i++ {
		fmt.Fprint(&sb, len(tc.pre[i]))
		for _, v := range tc.pre[i] {
			fmt.Fprintf(&sb, " %d", v)
		}
		sb.WriteByte('\n')
	}
	input := sb.String()
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	exp := solveA(tc)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		tc := genCaseA(rng)
		if err := runA(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
