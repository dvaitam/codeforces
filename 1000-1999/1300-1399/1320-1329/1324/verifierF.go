package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solve(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var n int
	fmt.Fscan(in, &n)
	val := make([]int, n+1)
	for i := 1; i <= n; i++ {
		var x int
		fmt.Fscan(in, &x)
		if x == 1 {
			val[i] = 1
		} else {
			val[i] = -1
		}
	}
	adj := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	parent := make([]int, n+1)
	order := make([]int, 0, n)
	stack := []int{1}
	parent[1] = 0
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, v)
		for _, u := range adj[v] {
			if u == parent[v] {
				continue
			}
			parent[u] = v
			stack = append(stack, u)
		}
	}
	dp := make([]int, n+1)
	for i := len(order) - 1; i >= 0; i-- {
		v := order[i]
		dp[v] = val[v]
		for _, u := range adj[v] {
			if u == parent[v] {
				continue
			}
			if dp[u] > 0 {
				dp[v] += dp[u]
			}
		}
	}
	up := make([]int, n+1)
	res := make([]int, n+1)
	stack = []int{1}
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		res[v] = dp[v] + up[v]
		for _, u := range adj[v] {
			if u == parent[v] {
				continue
			}
			tmp := up[v] + dp[v]
			if dp[u] > 0 {
				tmp -= dp[u]
			}
			if tmp < 0 {
				tmp = 0
			}
			up[u] = tmp
			stack = append(stack, u)
		}
	}
	var sb strings.Builder
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(res[i]))
	}
	return sb.String()
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(10) + 2
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(rng.Intn(2)))
	}
	sb.WriteByte('\n')
	for i := 2; i <= n; i++ {
		u := rng.Intn(i-1) + 1
		fmt.Fprintf(&sb, "%d %d\n", u, i)
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := generateCase(rng)
		exp := solve(in)
		got, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
