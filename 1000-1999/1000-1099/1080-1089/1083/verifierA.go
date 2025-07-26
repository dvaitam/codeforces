package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type test struct {
	input    string
	expected string
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func solveA(input string) string {
	rdr := bufio.NewReader(strings.NewReader(input))
	readInt := func() int {
		var x int
		var sign int = 1
		c, _ := rdr.ReadByte()
		for (c < '0' || c > '9') && c != '-' {
			c, _ = rdr.ReadByte()
		}
		if c == '-' {
			sign = -1
			c, _ = rdr.ReadByte()
		}
		for c >= '0' && c <= '9' {
			x = x*10 + int(c-'0')
			c, _ = rdr.ReadByte()
		}
		return x * sign
	}
	n := readInt()
	val := make([]int, n+1)
	f := make([]int64, n+1)
	g := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		val[i] = readInt()
		f[i] = int64(val[i])
		g[i] = int64(val[i])
	}
	type Edge struct{ to, w int }
	adj := make([][]Edge, n+1)
	for i := 2; i <= n; i++ {
		u := readInt()
		v := readInt()
		w := readInt()
		adj[u] = append(adj[u], Edge{v, w})
		adj[v] = append(adj[v], Edge{u, w})
	}
	type Frame struct{ u, parent, idx, wUp int }
	stack := make([]Frame, 0, n)
	stack = append(stack, Frame{1, 0, 0, 0})
	for len(stack) > 0 {
		top := &stack[len(stack)-1]
		u := top.u
		if top.idx < len(adj[u]) {
			e := adj[u][top.idx]
			top.idx++
			v := e.to
			w := e.w
			if v == top.parent {
				continue
			}
			if g[u] >= int64(w) {
				cand := g[u] - int64(w) + int64(val[v])
				if cand > g[v] {
					g[v] = cand
				}
			}
			stack = append(stack, Frame{v, u, 0, w})
		} else {
			cur := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if cur.parent != 0 {
				if f[cur.u] >= int64(cur.wUp) {
					cand := f[cur.u] - int64(cur.wUp) + int64(val[cur.parent])
					if cand > f[cur.parent] {
						f[cur.parent] = cand
					}
				}
				if f[cur.parent] > g[cur.parent] {
					g[cur.parent] = f[cur.parent]
				}
			}
		}
	}
	var ans int64
	for i := 1; i <= n; i++ {
		if g[i] > ans {
			ans = g[i]
		}
	}
	return fmt.Sprintf("%d", ans)
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(42))
	var tests []test
	for len(tests) < 100 {
		n := rng.Intn(6) + 1
		vals := make([]int, n)
		for i := 0; i < n; i++ {
			vals[i] = rng.Intn(10)
		}
		edges := make([][3]int, n-1)
		for i := 1; i < n; i++ {
			u := rng.Intn(i) + 1
			v := i + 1
			w := rng.Intn(10)
			edges[i-1] = [3]int{u, v, w}
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", vals[i]))
		}
		sb.WriteByte('\n')
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", e[0], e[1], e[2]))
		}
		input := sb.String()
		expected := solveA(input)
		tests = append(tests, test{input, expected})
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		out, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(t.expected) {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i+1, t.input, t.expected, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
