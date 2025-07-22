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

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
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

type query struct {
	op byte
	x  int
	y  int
}

func buildChildren(n int, parent []int) [][]int {
	ch := make([][]int, n+1)
	for i := 2; i <= n; i++ {
		p := parent[i]
		ch[p] = append(ch[p], i)
	}
	return ch
}

func computeDepth(n int, parent []int) []int {
	ch := buildChildren(n, parent)
	depth := make([]int, n+1)
	var dfs func(int, int)
	dfs = func(v, d int) {
		depth[v] = d
		for _, u := range ch[v] {
			dfs(u, d+1)
		}
	}
	dfs(1, 0)
	return depth
}

func lca(a, b int, parent []int, depth []int) int {
	for depth[a] > depth[b] {
		a = parent[a]
	}
	for depth[b] > depth[a] {
		b = parent[b]
	}
	for a != b {
		a = parent[a]
		b = parent[b]
	}
	return a
}

func expected(n int, parent []int, val []int) float64 {
	depth := computeDepth(n, parent)
	sum := 0.0
	for i := 1; i <= n; i++ {
		for j := 1; j <= n; j++ {
			l := lca(i, j, parent, depth)
			sum += float64(val[l])
		}
	}
	return sum / float64(n*n)
}

func inSubtree(n int, parent []int, v, u int) bool {
	ch := buildChildren(n, parent)
	stack := []int{v}
	for len(stack) > 0 {
		x := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if x == u {
			return true
		}
		for _, w := range ch[x] {
			stack = append(stack, w)
		}
	}
	return false
}

func solve(n int, parent []int, val []int, qs []query) []float64 {
	res := make([]float64, 0, len(qs)+1)
	res = append(res, expected(n, parent, val))
	for _, q := range qs {
		if q.op == 'V' {
			val[q.x] = q.y
		} else {
			if inSubtree(n, parent, q.x, q.y) {
				parent[q.y] = q.x
			} else {
				parent[q.x] = q.y
			}
		}
		res = append(res, expected(n, parent, val))
	}
	return res
}

func generateCase(rng *rand.Rand) (int, []int, []int, []query) {
	n := rng.Intn(5) + 2
	parent := make([]int, n+1)
	for i := 2; i <= n; i++ {
		parent[i] = rng.Intn(i-1) + 1
	}
	val := make([]int, n+1)
	for i := 1; i <= n; i++ {
		val[i] = rng.Intn(11)
	}
	q := rng.Intn(5) + 1
	qs := make([]query, q)
	for i := 0; i < q; i++ {
		if rng.Intn(2) == 0 {
			v := rng.Intn(n) + 1
			u := rng.Intn(n) + 1
			for u == v {
				u = rng.Intn(n) + 1
			}
			qs[i] = query{'P', v, u}
		} else {
			v := rng.Intn(n) + 1
			t := rng.Intn(11)
			qs[i] = query{'V', v, t}
		}
	}
	return n, parent, val, qs
}

func formatInput(n int, parent []int, val []int, qs []query) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 2; i <= n; i++ {
		if i > 2 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(parent[i]))
	}
	sb.WriteByte('\n')
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(val[i]))
	}
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("%d\n", len(qs)))
	for _, qu := range qs {
		if qu.op == 'P' {
			sb.WriteString(fmt.Sprintf("P %d %d\n", qu.x, qu.y))
		} else {
			sb.WriteString(fmt.Sprintf("V %d %d\n", qu.x, qu.y))
		}
	}
	return sb.String()
}

func check(exp []float64, out string) error {
	lines := strings.Fields(out)
	if len(lines) != len(exp) {
		return fmt.Errorf("expected %d numbers got %d", len(exp), len(lines))
	}
	for i, s := range lines {
		val, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return fmt.Errorf("invalid float")
		}
		diff := val - exp[i]
		if diff < -1e-6 || diff > 1e-6 {
			return fmt.Errorf("line %d expected %.6f got %.6f", i+1, exp[i], val)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, parent, val, qs := generateCase(rng)
		input := formatInput(n, append([]int(nil), parent...), append([]int(nil), val...), qs)
		exp := solve(n, append([]int(nil), parent...), append([]int(nil), val...), append([]query(nil), qs...))
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := check(exp, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
