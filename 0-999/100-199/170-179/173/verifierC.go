package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

func solveC(n, m int, a [][]int64) int64 {
	P := make([][]int64, n+1)
	for i := 0; i <= n; i++ {
		P[i] = make([]int64, m+1)
	}
	for i := 1; i <= n; i++ {
		rowSum := int64(0)
		for j := 1; j <= m; j++ {
			rowSum += a[i-1][j-1]
			P[i][j] = P[i-1][j] + rowSum
		}
	}
	maxSum := int64(-9e18)
	maxK := n
	if m < maxK {
		maxK = m
	}
	for k := 3; k <= maxK; k += 2 {
		kk := k
		for i := 0; i+kk <= n; i++ {
			i2 := i + kk
			for j := 0; j+kk <= m; j++ {
				j2 := j + kk
				sum := P[i2][j2] - P[i][j2] - P[i2][j] + P[i][j]
				if sum > maxSum {
					maxSum = sum
				}
			}
		}
	}
	return maxSum
}

func runCase(bin string, n, m int, a [][]int64) error {
	var input bytes.Buffer
	fmt.Fprintf(&input, "%d %d\n", n, m)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprintf(&input, "%d", a[i][j])
		}
		input.WriteByte('\n')
	}
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int64
	if _, err := fmt.Fscan(&out, &got); err != nil {
		return fmt.Errorf("parse error: %v", err)
	}
	want := solveC(n, m, a)
	if got != want {
		return fmt.Errorf("expected %d got %d", want, got)
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

	const tests = 100
	for t := 0; t < tests; t++ {
		n := rng.Intn(8) + 3
		m := rng.Intn(8) + 3
		a := make([][]int64, n)
		for i := 0; i < n; i++ {
			row := make([]int64, m)
			for j := 0; j < m; j++ {
				row[j] = int64(rng.Intn(2001) - 1000)
			}
			a[i] = row
		}
		if err := runCase(bin, n, m, a); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", t+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
