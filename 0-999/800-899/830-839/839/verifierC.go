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

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func expected(u, p int, adj [][]int) float64 {
	sum := 0.0
	cnt := 0
	for _, v := range adj[u] {
		if v == p {
			continue
		}
		sum += expected(v, u, adj) + 1
		cnt++
	}
	if cnt == 0 {
		return 0
	}
	return sum / float64(cnt)
}

func genCase(rng *rand.Rand) (int, [][2]int) {
	n := rng.Intn(20) + 1
	edges := make([][2]int, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges[i-2] = [2]int{p, i}
	}
	return n, edges
}

func solveCase(n int, edges [][2]int) float64 {
	adj := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	return expected(1, 0, adj)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, edges := genCase(rng)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for _, e := range edges {
			fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
		}
		input := sb.String()
		want := solveCase(n, edges)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\ninput:\n%s", i+1, err, input)
			return
		}
		got, err := strconv.ParseFloat(strings.TrimSpace(out), 64)
		if err != nil {
			fmt.Printf("case %d: bad float output %q\n", i+1, out)
			return
		}
		if diff := got - want; diff < -1e-6 || diff > 1e-6 {
			fmt.Printf("case %d failed\ninput:\n%sexpected: %.6f\ngot: %s\n", i+1, input, want, out)
			return
		}
	}
	fmt.Println("All tests passed")
}
