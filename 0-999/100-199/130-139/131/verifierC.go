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

func combTable(max int) [][]int64 {
	C := make([][]int64, max+1)
	for i := 0; i <= max; i++ {
		C[i] = make([]int64, i+1)
		C[i][0] = 1
		for j := 1; j <= i; j++ {
			if j == i {
				C[i][j] = 1
			} else {
				C[i][j] = C[i-1][j-1] + C[i-1][j]
			}
		}
	}
	return C
}

func solveCase(n, m, t int) string {
	maxNM := n
	if m > maxNM {
		maxNM = m
	}
	C := combTable(maxNM)
	var result int64
	for b := 4; b <= n; b++ {
		g := t - b
		if g < 1 || g > m {
			continue
		}
		result += C[n][b] * C[m][g]
	}
	return fmt.Sprint(result)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(27) + 4    // 4..30
	m := rng.Intn(30) + 1    // 1..30
	t := rng.Intn(n+m-4) + 5 // ensure between 5 and n+m
	if t > n+m {
		t = n + m
	}
	input := fmt.Sprintf("%d %d %d\n", n, m, t)
	expected := solveCase(n, m, t)
	return input, expected
}

func runCase(bin, input, expected string) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	if outStr != expected {
		return fmt.Errorf("expected %s got %s", expected, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
