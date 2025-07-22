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

type test struct {
	input    string
	expected string
}

func solve(input string) string {
	reader := strings.NewReader(input)
	var n int
	fmt.Fscan(reader, &n)
	capv := make([]int64, n+2)
	curr := make([]int64, n+2)
	parent := make([]int, n+2)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &capv[i])
		parent[i] = i
	}
	parent[n+1] = n + 1
	find := func(x int) int {
		for parent[x] != x {
			parent[x] = parent[parent[x]]
			x = parent[x]
		}
		return x
	}
	union := func(x, y int) { parent[find(x)] = find(y) }
	var m int
	fmt.Fscan(reader, &m)
	var out bytes.Buffer
	for i := 0; i < m; i++ {
		var tp int
		fmt.Fscan(reader, &tp)
		if tp == 1 {
			var p int
			var x int64
			fmt.Fscan(reader, &p, &x)
			idx := find(p)
			for idx <= n && x > 0 {
				space := capv[idx] - curr[idx]
				if space > x {
					curr[idx] += x
					x = 0
					break
				}
				curr[idx] = capv[idx]
				x -= space
				union(idx, idx+1)
				idx = find(idx)
			}
		} else {
			var k int
			fmt.Fscan(reader, &k)
			fmt.Fprintln(&out, curr[k])
		}
	}
	return strings.TrimSpace(out.String())
}

func generateTests() []test {
	rand.Seed(45)
	var tests []test
	for len(tests) < 100 {
		n := rand.Intn(5) + 1
		m := rand.Intn(10) + 1
		capv := make([]int64, n)
		for i := range capv {
			capv[i] = int64(rand.Intn(5) + 1)
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for i, v := range capv {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		fmt.Fprintf(&sb, "%d\n", m)
		for j := 0; j < m; j++ {
			if rand.Intn(2) == 0 {
				p := rand.Intn(n) + 1
				x := rand.Intn(5) + 1
				fmt.Fprintf(&sb, "1 %d %d\n", p, x)
			} else {
				k := rand.Intn(n) + 1
				fmt.Fprintf(&sb, "2 %d\n", k)
			}
		}
		inp := sb.String()
		tests = append(tests, test{inp, solve(inp)})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("time limit")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		out, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(t.expected) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%sexpected:%s\n got:%s\n", i+1, t.input, t.expected, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
