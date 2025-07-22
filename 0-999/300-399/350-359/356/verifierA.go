package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func solveA(in string) string {
	reader := bufio.NewReader(strings.NewReader(in))
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return ""
	}
	ans := make([]int, n+2)
	parent := make([]int, n+2)
	for i := 1; i <= n+1; i++ {
		parent[i] = i
	}
	var find func(int) int
	find = func(x int) int {
		for parent[x] != x {
			parent[x] = parent[parent[x]]
			x = parent[x]
		}
		return x
	}
	union := func(a, b int) {
		pa := find(a)
		pb := find(b)
		parent[pa] = pb
	}
	for i := 0; i < m; i++ {
		var l, r, x int
		fmt.Fscan(reader, &l, &r, &x)
		j := find(l)
		for j <= r {
			if j == x {
				j = find(j + 1)
				continue
			}
			ans[j] = x
			union(j, j+1)
			j = find(j)
		}
	}
	var sb strings.Builder
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(ans[i]))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func genTest(r *rand.Rand) string {
	n := r.Intn(9) + 2
	m := r.Intn(2*n) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 0; i < m; i++ {
		l := r.Intn(n-1) + 1
		rgt := l + r.Intn(n-l) + 1
		x := l + r.Intn(rgt-l+1)
		fmt.Fprintf(&sb, "%d %d %d\n", l, rgt, x)
	}
	return sb.String()
}

func runBinary(path, in string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, path)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: verifierA <path-to-binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	r := rand.New(rand.NewSource(1))
	const tests = 100
	for i := 0; i < tests; i++ {
		in := genTest(r)
		expect := strings.TrimSpace(solveA(in))
		got, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, in, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
