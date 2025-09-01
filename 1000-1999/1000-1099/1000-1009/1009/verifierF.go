package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// runProg executes a binary with given input and returns trimmed stdout or an error including stderr.
func runProg(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCaseF(rng *rand.Rand) (int, [][2]int) {
	n := rng.Intn(15) + 1
	if n == 1 {
		return n, nil
	}
	edges := make([][2]int, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges[i-2] = [2]int{p, i}
	}
	return n, edges
}

// parseInts parses any whitespace-separated list of integers.
func parseInts(s string) ([]int, error) {
	fields := strings.Fields(strings.TrimSpace(s))
	res := make([]int, len(fields))
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("invalid int %q", f)
		}
		res[i] = v
	}
	return res, nil
}

func intsToString(a []int) string {
	if len(a) == 0 {
		return ""
	}
	var sb strings.Builder
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	return sb.String()
}

// expectedAns1009F computes answers: for each node u, the smallest d maximizing the
// count of nodes at distance d within u's subtree (rooted at 1).
func expectedAns1009F(n int, edges [][2]int) []int {
	adj := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	depth := make([]int, n+1)
	tin := make([]int, n+1)
	tout := make([]int, n+1)
	timeD := 0
	var dfs func(u, p, d int)
	dfs = func(u, p, d int) {
		timeD++
		tin[u] = timeD
		depth[u] = d
		for _, v := range adj[u] {
			if v == p {
				continue
			}
			dfs(v, u, d+1)
		}
		tout[u] = timeD
	}
	dfs(1, 0, 0)

	ans := make([]int, n+1)
	for u := 1; u <= n; u++ {
		freq := make(map[int]int)
		maxc, best := 0, 0
		for v := 1; v <= n; v++ {
			if tin[u] <= tin[v] && tin[v] <= tout[u] {
				d := depth[v] - depth[u]
				freq[d]++
				c := freq[d]
				if c > maxc || (c == maxc && d < best) {
					maxc, best = c, d
				}
			}
		}
		ans[u] = best
	}
	return ans
}

func runCaseF(bin string, n int, edges [][2]int) error {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(n))
	sb.WriteByte('\n')
	for i, e := range edges {
		if i > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(strconv.Itoa(e[0]))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(e[1]))
	}
	if n > 1 {
		sb.WriteByte('\n')
	}
	input := sb.String()
	got, err := runProg(bin, input)
	if err != nil {
		return err
	}
	gotList, err := parseInts(got)
	if err != nil {
		return fmt.Errorf("failed to parse candidate output: %v (output=%q)", err, got)
	}
	if len(gotList) != n {
		return fmt.Errorf("length mismatch: expected %d values got %d", n, len(gotList))
	}
	expListFull := expectedAns1009F(n, edges)
	expList := expListFull[1:]
	for i := 0; i < n; i++ {
		if gotList[i] != expList[i] {
			return fmt.Errorf("mismatch at pos %d: expected %d got %d\nexpected %s got %s", i+1, expList[i], gotList[i], intsToString(expList), intsToString(gotList))
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, edges := genCaseF(rng)
		if err := runCaseF(bin, n, edges); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
