package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCaseA struct {
	n   int
	pre [][]int
}

func generateCaseA(rng *rand.Rand) testCaseA {
	n := rng.Intn(7) + 1 // 1..8
	pre := make([][]int, n)
	for i := 0; i < n; i++ {
		k := rng.Intn(n)
		if k > n-1 {
			k = n - 1
		}
		deps := make([]int, 0, k)
		used := make(map[int]bool)
		for len(deps) < k {
			x := rng.Intn(n)
			if x == i || used[x] {
				continue
			}
			used[x] = true
			deps = append(deps, x)
		}
		pre[i] = deps
	}
	return testCaseA{n: n, pre: pre}
}

func buildInputA(t testCaseA) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprint(t.n))
	sb.WriteByte('\n')
	for i := 0; i < t.n; i++ {
		sb.WriteString(fmt.Sprint(len(t.pre[i])))
		for _, v := range t.pre[i] {
			sb.WriteByte(' ')
			sb.WriteString(fmt.Sprint(v + 1))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func solveA(reader *bufio.Reader) string {
	var T int
	fmt.Fscan(reader, &T)
	var out strings.Builder
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		adj := make([][]int, n)
		pre := make([][]int, n)
		indeg := make([]int, n)
		for i := 0; i < n; i++ {
			var k int
			fmt.Fscan(reader, &k)
			pre[i] = make([]int, k)
			for j := 0; j < k; j++ {
				var x int
				fmt.Fscan(reader, &x)
				x--
				pre[i][j] = x
				adj[x] = append(adj[x], i)
				indeg[i]++
			}
		}
		queue := make([]int, 0, n)
		for i := 0; i < n; i++ {
			if indeg[i] == 0 {
				queue = append(queue, i)
			}
		}
		order := make([]int, 0, n)
		for len(queue) > 0 {
			u := queue[0]
			queue = queue[1:]
			order = append(order, u)
			for _, v := range adj[u] {
				indeg[v]--
				if indeg[v] == 0 {
					queue = append(queue, v)
				}
			}
		}
		if len(order) != n {
			out.WriteString("-1\n")
			continue
		}
		dp := make([]int, n)
		ans := 0
		for _, u := range order {
			dp[u] = 1
			for _, v := range pre[u] {
				cand := dp[v]
				if v > u {
					cand++
				}
				if cand > dp[u] {
					dp[u] = cand
				}
			}
			if dp[u] > ans {
				ans = dp[u]
			}
		}
		out.WriteString(fmt.Sprintf("%d\n", ans))
	}
	return strings.TrimSpace(out.String())
}

func expectedA(t testCaseA) string {
	input := buildInputA(t)
	return solveA(bufio.NewReader(strings.NewReader(input)))
}

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return out.String(), fmt.Errorf("timeout")
	}
	if err != nil {
		return out.String(), err
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseA(rng)
		input := buildInputA(tc)
		expect := expectedA(tc)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\nOutput:%s", i+1, err, out)
			os.Exit(1)
		}
		got := strings.TrimSpace(out)
		exp := strings.TrimSpace(expect)
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
