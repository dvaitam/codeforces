package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type pair struct {
	val int
	idx int
}

func solveCase(a []int) (int64, int, []int) {
	n := len(a)
	pairs := make([]pair, n)
	for i := 0; i < n; i++ {
		pairs[i] = pair{val: a[i], idx: i}
	}
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].val < pairs[j].val })
	const INF int64 = 1<<63 - 1
	dp := make([]int64, n+1)
	prev := make([]int, n+1)
	for i := range dp {
		dp[i] = INF
	}
	dp[0] = 0
	for i := 3; i <= n; i++ {
		for s := 3; s <= 5; s++ {
			if i-s < 0 {
				continue
			}
			diff := pairs[i-1].val - pairs[i-s].val
			cost := dp[i-s] + int64(diff)
			if cost < dp[i] {
				dp[i] = cost
				prev[i] = s
			}
		}
	}
	teams := make([]int, n)
	teamCount := 0
	for i := n; i > 0; {
		s := prev[i]
		teamCount++
		for j := i - s; j < i; j++ {
			teams[pairs[j].idx] = teamCount
		}
		i -= s
	}
	return dp[n], teamCount, teams
}

func generateCase() (string, string) {
	n := rand.Intn(10) + 3
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rand.Intn(2000) + 1
	}
	cost, teamCount, teams := solveCase(arr)
	var in strings.Builder
	fmt.Fprintf(&in, "%d\n", n)
	for i := 0; i < n; i++ {
		if i+1 == n {
			fmt.Fprintf(&in, "%d\n", arr[i])
		} else {
			fmt.Fprintf(&in, "%d ", arr[i])
		}
	}
	var out strings.Builder
	fmt.Fprintf(&out, "%d %d\n", cost, teamCount)
	for i := 0; i < n; i++ {
		if i+1 == n {
			fmt.Fprintf(&out, "%d\n", teams[i])
		} else {
			fmt.Fprintf(&out, "%d ", teams[i])
		}
	}
	return in.String(), out.String()
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, buf.String())
	}
	return strings.TrimSpace(buf.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(5)
	for i := 0; i < 100; i++ {
		in, exp := generateCase()
		got, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected:\n%s\ngot:\n%s\n", i+1, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
