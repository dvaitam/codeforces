package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type TestB struct {
	n       int
	parents []int
	l       []int64
	r       []int64
}

func solveB(t TestB) int {
	g := make([][]int, t.n)
	for i := 1; i < t.n; i++ {
		p := t.parents[i-1] - 1
		g[p] = append(g[p], i)
	}
	var dfs func(int) int64
	res := 0
	dfs = func(v int) int64 {
		sum := int64(0)
		for _, to := range g[v] {
			sum += dfs(to)
		}
		if sum < t.l[v] {
			res++
			return t.r[v]
		}
		if sum > t.r[v] {
			return t.r[v]
		}
		return sum
	}
	dfs(0)
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(43)
	const tCases = 100
	var sb strings.Builder
	var exp strings.Builder
	sb.WriteString(fmt.Sprintln(tCases))
	for i := 0; i < tCases; i++ {
		n := rand.Intn(9) + 2
		parents := make([]int, n-1)
		for j := 1; j < n; j++ {
			parents[j-1] = rand.Intn(j) + 1
		}
		l := make([]int64, n)
		r := make([]int64, n)
		for j := 0; j < n; j++ {
			low := rand.Intn(10) + 1
			high := low + rand.Intn(5)
			l[j] = int64(low)
			r[j] = int64(high)
		}
		sb.WriteString(fmt.Sprintln(n))
		for j, p := range parents {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", p))
		}
		sb.WriteByte('\n')
		for j := 0; j < n; j++ {
			sb.WriteString(fmt.Sprintf("%d %d\n", l[j], r[j]))
		}
		exp.WriteString(fmt.Sprintf("%d\n", solveB(TestB{n, parents, l, r})))
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error running binary: %v\noutput:\n%s", err, out.String())
		os.Exit(1)
	}
	got := strings.TrimSpace(out.String())
	want := strings.TrimSpace(exp.String())
	if got != want {
		fmt.Fprintf(os.Stderr, "wrong answer\nexpected:\n%s\ngot:\n%s\n", want, got)
		os.Exit(1)
	}
	fmt.Println("all tests passed")
}
