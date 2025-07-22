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

func solveE(n int, edges [][2]int, stVals, enVals []int64) string {
	adj := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	st := make([]int64, n+1)
	en := make([]int64, n+1)
	var sst, sen int64
	for i := 1; i <= n; i++ {
		st[i] = stVals[i-1]
		en[i] = enVals[i-1]
		sst += st[i]
		sen += en[i]
	}
	parent := make([]int, n+1)
	order := make([]int, 0, n)
	stack := []int{1}
	parent[1] = 0
	for len(stack) > 0 {
		x := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, x)
		for _, y := range adj[x] {
			if y != parent[x] {
				parent[y] = x
				stack = append(stack, y)
			}
		}
	}
	size := make([]int64, n+1)
	stSum := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		size[i] = 1
		stSum[i] = st[i]
	}
	var ans float64
	for i := len(order) - 1; i >= 0; i-- {
		x := order[i]
		p := parent[x]
		if p != 0 {
			size[p] += size[x]
			stSum[p] += stSum[x]
			ans += float64(stSum[x]) * float64(size[x]) * float64(en[p])
		}
		ans += float64(sst-stSum[x]) * float64(int64(n)-size[x]) * float64(en[x])
	}
	res := ans / float64(sst) / float64(sen)
	return fmt.Sprintf("%.11f", res)
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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) (string, int, [][2]int, []int64, []int64) {
	n := rng.Intn(5) + 1
	edges := make([][2]int, 0, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges = append(edges, [2]int{i, p})
	}
	stVals := make([]int64, n)
	enVals := make([]int64, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	for i := 0; i < n; i++ {
		stVals[i] = int64(rng.Intn(10) + 1)
		enVals[i] = int64(rng.Intn(10) + 1)
		sb.WriteString(fmt.Sprintf("%d %d\n", stVals[i], enVals[i]))
	}
	return sb.String(), n, edges, stVals, enVals
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, n, edges, stVals, enVals := genCase(rng)
		expect := solveE(n, edges, stVals, enVals)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
