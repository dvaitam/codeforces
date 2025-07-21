package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func solveF(n int, edges [][2]int) string {
	adj := make([][]int, n)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	w := (n + 63) >> 6
	bs := make([][]uint64, n)
	for i := 0; i < n; i++ {
		b := make([]uint64, w)
		b[i>>6] |= 1 << (uint(i) & 63)
		for _, v := range adj[i] {
			b[v>>6] |= 1 << (uint(v) & 63)
		}
		bs[i] = b
	}
	best := 0
	for i := 0; i < n; i++ {
		cnt := 0
		for j := 0; j < w; j++ {
			cnt += bits.OnesCount64(bs[i][j])
		}
		if cnt > best {
			best = cnt
		}
	}
	tmp := make([]uint64, w)
	for _, e := range edges {
		u, v := e[0], e[1]
		cnt := 0
		bu, bv := bs[u], bs[v]
		for j := 0; j < w; j++ {
			tmp[j] = bu[j] | bv[j]
			cnt += bits.OnesCount64(tmp[j])
		}
		if cnt > best {
			best = cnt
		}
	}
	res := n - best
	return fmt.Sprintf("%d", res)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 2
	edges := make([][2]int, 0)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if rng.Float32() < 0.3 {
				edges = append(edges, [2]int{i, j})
			}
		}
	}
	m := len(edges)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0]+1, e[1]+1))
	}
	ans := solveF(n, edges)
	return sb.String(), ans
}

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

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expect := generateCase(rng)
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
