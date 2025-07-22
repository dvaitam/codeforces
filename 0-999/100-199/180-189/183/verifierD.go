package main

import (
	"bytes"
	"fmt"
	"math"
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func computeExpected(n, m int, prob [][]int) float64 {
	a := make([][]float64, n+1)
	for i := range a {
		a[i] = make([]float64, m+1)
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			a[i][j] = float64(prob[i-1][j-1]) / 1000.0
		}
	}
	g := make([][][2]float64, m+1)
	for j := 0; j <= m; j++ {
		g[j] = make([][2]float64, n+1)
	}
	G := make([]bool, m+1)
	var ans float64
	for j := 1; j <= m; j++ {
		idx := 0
		if G[j] {
			idx = 1
		}
		for i := 1; i <= n; i++ {
			other := idx ^ 1
			g[j][i][idx] = (g[j][i-1][other]+1)*a[i][j] + g[j][i-1][idx]*(1-a[i][j])
		}
	}
	for t := 1; t <= n; t++ {
		mx := 0.0
		p := 1
		for j := 1; j <= m; j++ {
			idx := 0
			if G[j] {
				idx = 1
			}
			diff := g[j][n][idx] - g[j][n][idx^1]
			if diff > mx {
				mx = diff
				p = j
			}
		}
		ans += mx
		G[p] = !G[p]
		idx := 0
		if G[p] {
			idx = 1
		}
		other := idx ^ 1
		for i := 1; i <= n; i++ {
			g[p][i][idx] = (g[p][i-1][other]+1)*a[i][p] + g[p][i-1][idx]*(1-a[i][p])
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(3) + 1
	m := rng.Intn(4) + 1
	prob := make([][]int, n)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 0; i < n; i++ {
		prob[i] = make([]int, m)
		for j := 0; j < m; j++ {
			prob[i][j] = rng.Intn(1001)
			fmt.Fprintf(&sb, "%d ", prob[i][j])
		}
		sb.WriteByte('\n')
	}
	ans := computeExpected(n, m, prob)
	return sb.String(), fmt.Sprintf("%.10f", ans)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		var got float64
		if _, err := fmt.Sscan(out, &got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d: bad output: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		expVal, _ := strconv.ParseFloat(exp, 64)
		if math.Abs(got-expVal) > 1e-6 {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %.10f got %.10f\ninput:\n%s", i+1, expVal, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
