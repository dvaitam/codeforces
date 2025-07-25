package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type query struct {
	v int
	h int
}

func runCandidate(bin string, input string) (string, error) {
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

func buildTree(n int, parents []int) [][]int {
	g := make([][]int, n+1)
	for i := 2; i <= n; i++ {
		p := parents[i-2]
		g[p] = append(g[p], i)
	}
	return g
}

func isAncestor(parents []int, v, u int) bool {
	for u != 0 {
		if u == v {
			return true
		}
		if u == 1 {
			break
		}
		u = parents[u-2]
	}
	return v == 1 && u == 1
}

func expected(n int, parents []int, letters string, qs []query) string {
	g := buildTree(n, parents)
	depth := make([]int, n+1)
	depth[1] = 1
	stack := []int{1}
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, u := range g[v] {
			depth[u] = depth[v] + 1
			stack = append(stack, u)
		}
	}
	var sb strings.Builder
	for idx, q := range qs {
		counts := make([]int, 26)
		for u := 1; u <= n; u++ {
			if depth[u] == q.h && isAncestor(parents, q.v, u) {
				counts[letters[u-1]-'a']++
			}
		}
		odd := 0
		for _, c := range counts {
			if c%2 == 1 {
				odd++
			}
		}
		if idx > 0 {
			sb.WriteByte('\n')
		}
		if odd <= 1 {
			sb.WriteString("Yes")
		} else {
			sb.WriteString("No")
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	lettersArr := []byte("abcde")
	for i := 0; i < 100; i++ {
		n := rng.Intn(7) + 1
		m := rng.Intn(7) + 1
		parents := make([]int, n-1)
		for j := 2; j <= n; j++ {
			parents[j-2] = rng.Intn(j-1) + 1
		}
		letters := make([]byte, n)
		for j := 0; j < n; j++ {
			letters[j] = lettersArr[rng.Intn(len(lettersArr))]
		}
		qs := make([]query, m)
		for j := 0; j < m; j++ {
			qs[j] = query{rng.Intn(n) + 1, rng.Intn(n) + 1}
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, m)
		for j := 0; j < n-1; j++ {
			fmt.Fprintf(&sb, "%d ", parents[j])
		}
		if n > 1 {
			sb.WriteByte('\n')
		}
		sb.WriteString(string(letters))
		sb.WriteByte('\n')
		for _, q := range qs {
			fmt.Fprintf(&sb, "%d %d\n", q.v, q.h)
		}
		input := sb.String()
		want := expected(n, parents, string(letters), qs)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\n", i+1, want, strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
