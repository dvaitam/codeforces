package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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
	return out.String(), err
}

func bfsFarthest(start int, adj [][]int) (int, int) {
	n := len(adj)
	dist := make([]int, n)
	for i := range dist {
		dist[i] = -1
	}
	q := []int{start}
	dist[start] = 0
	farNode, farDist := start, 0
	for len(q) > 0 {
		u := q[0]
		q = q[1:]
		for _, v := range adj[u] {
			if dist[v] == -1 {
				dist[v] = dist[u] + 1
				q = append(q, v)
			}
		}
		if dist[u] > farDist {
			farDist = dist[u]
			farNode = u
		}
	}
	return farNode, farDist
}

func expected(n int, spiders []string) string {
	idx := 0
	total := 0
	for i := 0; i < n; i++ {
		parts := strings.Fields(spiders[i])
		ni, _ := strconv.Atoi(parts[0])
		edges := make([][2]int, ni-1)
		for j := 0; j < ni-1; j++ {
			u, _ := strconv.Atoi(parts[1+2*j])
			v, _ := strconv.Atoi(parts[1+2*j+1])
			u--
			v--
			edges[j] = [2]int{u, v}
		}
		adj := make([][]int, ni)
		for _, e := range edges {
			u, v := e[0], e[1]
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}
		far, _ := bfsFarthest(0, adj)
		_, d := bfsFarthest(far, adj)
		total += d
		idx += ni - 1
	}
	return fmt.Sprintf("%d\n", total)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(3) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	spiders := make([]string, n)
	for i := 0; i < n; i++ {
		ni := rng.Intn(3) + 2
		edges := make([][2]int, ni-1)
		for j := 1; j < ni; j++ {
			p := rng.Intn(j)
			edges[j-1] = [2]int{p, j}
		}
		fmt.Fprintf(&sb, "%d", ni)
		var line strings.Builder
		fmt.Fprintf(&line, "%d", ni)
		for _, e := range edges {
			fmt.Fprintf(&sb, " %d %d", e[0]+1, e[1]+1)
			fmt.Fprintf(&line, " %d %d", e[0]+1, e[1]+1)
		}
		sb.WriteByte('\n')
		spiders[i] = line.String()
	}
	return sb.String(), expected(n, spiders)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: go run verifierF.go /path/to/binary\n")
		os.Exit(1)
	}
	candidate := os.Args[1]

	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	ref := filepath.Join(dir, "refF")
	cmd := exec.Command("go", "build", "-o", ref, filepath.Join(dir, "120F.go"))
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference solution: %v\n%s", err, out)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		input, expect := generateCase(rng)
		candOut, cErr := runBinary(candidate, input)
		refOut, rErr := runBinary(ref, input)
		if cErr != nil {
			fmt.Fprintf(os.Stderr, "test %d: candidate error: %v\n", t+1, cErr)
			os.Exit(1)
		}
		if rErr != nil {
			fmt.Fprintf(os.Stderr, "test %d: reference error: %v\n", t+1, rErr)
			os.Exit(1)
		}
		if strings.TrimSpace(candOut) != strings.TrimSpace(refOut) {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%sexpected:%sactual:%s\n", t+1, input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
