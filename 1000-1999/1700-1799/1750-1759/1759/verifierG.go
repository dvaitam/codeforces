package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
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

func solveCase(n int, b []int) (string, bool) {
	m := n / 2
	used := make([]bool, n+1)
	for _, v := range b {
		if v < 1 || v > n || used[v] {
			return "", false
		}
		used[v] = true
	}
	avail := make([]int, 0, m)
	for i := 1; i <= n; i++ {
		if !used[i] {
			avail = append(avail, i)
		}
	}
	parent := make([]int, len(avail))
	for i := range parent {
		parent[i] = i
	}
	var find func(int) int
	find = func(x int) int {
		if x < 0 {
			return -1
		}
		if parent[x] == x {
			return x
		}
		parent[x] = find(parent[x])
		return parent[x]
	}
	remove := func(x int) { parent[x] = find(x - 1) }
	res := make([]int, n)
	for i := m - 1; i >= 0; i-- {
		idx := sort.SearchInts(avail, b[i]) - 1
		idx = find(idx)
		if idx < 0 {
			return "", false
		}
		res[2*i] = avail[idx]
		res[2*i+1] = b[i]
		remove(idx)
	}
	var sb strings.Builder
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(res[i]))
	}
	return sb.String(), true
}

func generateCase(rng *rand.Rand) (string, string) {
	n := (rng.Intn(5) + 1) * 2
	m := n / 2
	used := map[int]bool{}
	b := make([]int, m)
	for i := 0; i < m; i++ {
		for {
			v := rng.Intn(n) + 1
			if !used[v] {
				used[v] = true
				b[i] = v
				break
			}
		}
	}
	out, ok := solveCase(n, b)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteString("\n")
	if !ok {
		return sb.String(), "-1"
	}
	return sb.String(), out
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
