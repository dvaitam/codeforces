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

func solveC(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var n, m int
	fmt.Fscan(in, &n, &m)
	s := make([]int, n+1)
	a := make([][]int, n+1)
	cnt := 0
	root := 0
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &s[i])
		if s[i] > cnt {
			cnt = s[i]
			root = i
		}
		if s[i] > 0 {
			a[i] = make([]int, s[i])
			for j := 0; j < s[i]; j++ {
				fmt.Fscan(in, &a[i][j])
			}
		}
	}
	g := make([][]int, n+1)
	for i := 1; i < n; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}
	used := make([]int, m+2)
	col := make([]int, m+2)
	parent := make([]int, n+1)
	stack := []int{root}
	for len(stack) > 0 {
		x := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, y := range a[x] {
			if col[y] != 0 {
				used[col[y]] = x
			}
		}
		j := 1
		for _, y := range a[x] {
			if col[y] != 0 {
				continue
			}
			for used[j] == x {
				j++
			}
			col[y] = j
			j++
		}
		for _, y := range g[x] {
			if y != parent[x] {
				parent[y] = x
				stack = append(stack, y)
			}
		}
	}
	if cnt == 0 {
		cnt = 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", cnt))
	for i := 1; i <= m; i++ {
		if col[i] == 0 {
			col[i] = 1
		}
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", col[i]))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genTest() string {
	r := rand.New(rand.NewSource(rand.Int63()))
	n := r.Intn(3) + 1
	m := n
	lines := []string{fmt.Sprintf("%d %d", n, m)}
	for i := 1; i <= n; i++ {
		lines = append(lines, fmt.Sprintf("1 %d", i))
	}
	for i := 1; i < n; i++ {
		lines = append(lines, fmt.Sprintf("%d %d", i, i+1))
	}
	return strings.Join(lines, "\n") + "\n"
}

func generateTests() []string {
	tests := make([]string, 100)
	rand.Seed(3)
	for i := 0; i < 100; i++ {
		tests[i] = genTest()
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		expected := strings.TrimSpace(solveC(t))
		got, err := runBinary(bin, t)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("test %d failed. expected %s got %s\ninput:\n%s", i+1, expected, got, t)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
