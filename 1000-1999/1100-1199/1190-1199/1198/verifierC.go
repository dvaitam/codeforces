package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type edge struct{ u, v int }

func run(bin string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveC(n, m int, edges []edge) (string, []int) {
	used := make([]bool, 3*n+1)
	matching := make([]int, 0, n)
	for i := 0; i < m; i++ {
		u := edges[i].u
		v := edges[i].v
		if !used[u] && !used[v] {
			used[u] = true
			used[v] = true
			matching = append(matching, i+1)
			if len(matching) == n {
				break
			}
		}
	}
	if len(matching) >= n {
		return "Matching", matching[:n]
	}
	ind := make([]int, 0, n)
	for v := 1; v <= 3*n && len(ind) < n; v++ {
		if !used[v] {
			ind = append(ind, v)
		}
	}
	return "IndSet", ind
}

func generateCase(rng *rand.Rand) ([]byte, string) {
	n := rng.Intn(3) + 1
	m := rng.Intn(3*n) + n
	edges := make([]edge, m)
	var b bytes.Buffer
	fmt.Fprintf(&b, "1\n")
	fmt.Fprintf(&b, "%d %d\n", n, m)
	for i := 0; i < m; i++ {
		u := rng.Intn(3*n) + 1
		v := rng.Intn(3*n) + 1
		edges[i] = edge{u, v}
		fmt.Fprintf(&b, "%d %d\n", u, v)
	}
	kind, ans := solveC(n, m, edges)
	var expect bytes.Buffer
	fmt.Fprintln(&expect, kind)
	for i, v := range ans {
		if i > 0 {
			expect.WriteByte(' ')
		}
		fmt.Fprint(&expect, v)
	}
	return b.Bytes(), strings.TrimSpace(expect.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input, expect := generateCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, string(input))
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i, expect, strings.TrimSpace(out), string(input))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
