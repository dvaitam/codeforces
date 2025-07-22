package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type door struct {
	a, b     int
	cab, cba int
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(5) + 1 // 1..5 nodes
	pairs := make([][2]int, 0)
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			pairs = append(pairs, [2]int{i, j})
		}
	}
	rng.Shuffle(len(pairs), func(i, j int) { pairs[i], pairs[j] = pairs[j], pairs[i] })
	m := rng.Intn(len(pairs) + 1)
	doors := make([]door, 0, m)
	for i := 0; i < m; i++ {
		p := pairs[i]
		doors = append(doors, door{
			a:   p[0],
			b:   p[1],
			cab: rng.Intn(21) - 10,
			cba: rng.Intn(21) - 10,
		})
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for _, d := range doors {
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", d.a, d.b, d.cab, d.cba))
	}
	input := sb.String()
	expect := solveCase(n, doors)
	return input, expect
}

func solveCase(n int, doors []door) int {
	const INF int64 = math.MaxInt64 / 4
	type edge struct {
		u, v int
		w    int64
	}
	edges := make([]edge, 0, len(doors)*2)
	for _, d := range doors {
		edges = append(edges, edge{d.a - 1, d.b - 1, int64(d.cab)})
		edges = append(edges, edge{d.b - 1, d.a - 1, int64(d.cba)})
	}
	dpPrev := make([][]int64, n)
	dpCur := make([][]int64, n)
	for i := 0; i < n; i++ {
		dpPrev[i] = make([]int64, n)
		dpCur[i] = make([]int64, n)
		for j := 0; j < n; j++ {
			dpPrev[i][j] = -INF
			dpCur[i][j] = -INF
		}
	}
	for _, e := range edges {
		if e.w > dpPrev[e.u][e.v] {
			dpPrev[e.u][e.v] = e.w
		}
	}
	for k := 2; k <= n; k++ {
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				dpCur[i][j] = -INF
			}
		}
		for _, e := range edges {
			u, v, w := e.u, e.v, e.w
			for i := 0; i < n; i++ {
				if dpPrev[i][u] > -INF {
					val := dpPrev[i][u] + w
					if val > dpCur[i][v] {
						dpCur[i][v] = val
					}
				}
			}
		}
		for i := 0; i < n; i++ {
			if dpCur[i][i] > 0 {
				return k
			}
		}
		dpPrev, dpCur = dpCur, dpPrev
	}
	return 0
}

func runCandidate(bin, input string) (string, error) {
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
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
		if out != fmt.Sprintf("%d", exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
