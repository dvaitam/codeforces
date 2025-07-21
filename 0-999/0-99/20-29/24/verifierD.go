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

func expected(N, M, si, sj int) float64 {
	si--
	sj--
	
	if si == N-1 {
		return 0.0
	}
	
	D := make([]float64, M)
	a := make([]float64, M)
	b := make([]float64, M)
	c := make([]float64, M)
	d := make([]float64, M)
	cp := make([]float64, M)
	dp := make([]float64, M)
	x := make([]float64, M)
	
	for r := N - 2; r >= 0; r-- {
		for j := 0; j < M; j++ {
			down := D[j]
			switch {
			case j == 0 && M > 1:
				a[j] = 0
				b[j] = 2
				c[j] = -1
				d[j] = 3 + down
			case j == M-1 && M > 1:
				a[j] = -1
				b[j] = 2
				c[j] = 0
				d[j] = 3 + down
			case M == 1:
				a[j] = 0
				b[j] = 1
				c[j] = 0
				d[j] = 2 + down
			default:
				a[j] = -1
				b[j] = 3
				c[j] = -1
				d[j] = 4 + down
			}
		}
		
		cp[0] = c[0] / b[0]
		dp[0] = d[0] / b[0]
		for j := 1; j < M; j++ {
			denom := b[j] - a[j]*cp[j-1]
			cp[j] = c[j] / denom
			dp[j] = (d[j] - a[j]*dp[j-1]) / denom
		}
		
		x[M-1] = dp[M-1]
		for j := M - 2; j >= 0; j-- {
			x[j] = dp[j] - cp[j]*x[j+1]
		}
		
		copy(D, x)
	}
	
	return D[sj]
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
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	for i := 0; i < 100; i++ {
		N := rng.Intn(20) + 1
		M := rng.Intn(20) + 1
		si := rng.Intn(N) + 1
		sj := rng.Intn(M) + 1
		
		input := fmt.Sprintf("%d %d\n%d %d\n", N, M, si, sj)
		
		expectedOut := expected(N, M, si, sj)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		
		gotFloat, parseErr := strconv.ParseFloat(got, 64)
		if parseErr != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: cannot parse output: %v\ninput:\n%s", i+1, parseErr, input)
			os.Exit(1)
		}
		
		if math.Abs(gotFloat-expectedOut) > 1e-3 {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %.4f got %.4f (diff %.6f)\ninput:\n%s", i+1, expectedOut, gotFloat, math.Abs(gotFloat-expectedOut), input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}