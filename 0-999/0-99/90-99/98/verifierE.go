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

const maxN = 1005

var memo [maxN][maxN]float64
var seen [maxN][maxN]bool

func dfs(x, y int) float64 {
	if x == 0 || y == 0 {
		return 1.0 / float64(y+1)
	}
	if seen[x][y] {
		return memo[x][y]
	}
	seen[x][y] = true
	a := 1 - dfs(y, x-1)
	b := float64(y) / float64(y+1) * (1 - dfs(y-1, x))
	c := b + 1.0/float64(y+1)
	p := (c - b) / (1 - a - b + c)
	memo[x][y] = p*(a-c) + c
	return memo[x][y]
}

func solve(n, m int) string {
	res := dfs(n, m)
	return fmt.Sprintf("%.10f %.10f", res, 1-res)
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(1001)
	m := rng.Intn(1001)
	return fmt.Sprintf("%d %d\n", n, m)
}

func runCase(bin string, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(strings.TrimSpace(out.String()))
	expFields := strings.Fields(expected)
	if len(fields) != 2 {
		return fmt.Errorf("expected two numbers got %v", fields)
	}
	var got1, got2 float64
	if _, err := fmt.Sscan(fields[0], &got1); err != nil {
		return fmt.Errorf("bad out: %v", err)
	}
	if _, err := fmt.Sscan(fields[1], &got2); err != nil {
		return fmt.Errorf("bad out: %v", err)
	}
	var exp1, exp2 float64
	fmt.Sscan(expFields[0], &exp1)
	fmt.Sscan(expFields[1], &exp2)
	if math.Abs(got1-exp1) > 1e-9 || math.Abs(got2-exp2) > 1e-9 {
		return fmt.Errorf("expected %s got %s", expected, strings.Join(fields, " "))
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		var n, m int
		fmt.Sscan(input, &n, &m)
		// reset memo arrays
		for a := 0; a <= n; a++ {
			for b := 0; b <= m; b++ {
				seen[a][b] = false
			}
		}
		expected := solve(n, m)
		if err := runCase(bin, input, expected); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
