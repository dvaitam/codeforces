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

type Matrix [][]int

func generateMatrix(rng *rand.Rand) Matrix {
	n := rng.Intn(3) + 2 // 2..4
	m := make(Matrix, n)
	for i := 0; i < n; i++ {
		m[i] = make([]int, n)
		for j := 0; j < n; j++ {
			if i == j {
				m[i][j] = 0
			} else {
				m[i][j] = rng.Intn(10)
			}
		}
	}
	return m
}

func expectedRowMins(mat Matrix) []int {
	n := len(mat)
	res := make([]int, n)
	for i := 0; i < n; i++ {
		best := math.MaxInt32
		for j := 0; j < n; j++ {
			if i == j {
				continue
			}
			if mat[i][j] < best {
				best = mat[i][j]
			}
		}
		res[i] = best
	}
	return res
}

func runCase(bin string, mat Matrix) error {
	n := len(mat)
	var in bytes.Buffer
	fmt.Fprintln(&in, n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if j > 0 {
				fmt.Fprint(&in, " ")
			}
			fmt.Fprint(&in, mat[i][j])
		}
		fmt.Fprintln(&in)
	}

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(bin)
	cmd.Stdin = &in
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v stderr:%s", err, stderr.String())
	}

	tokens := strings.Fields(out.String())
	if len(tokens) != n+1 || tokens[0] != "-1" {
		return fmt.Errorf("unexpected output: %s stderr:%s", out.String(), stderr.String())
	}
	ans := make([]int, n)
	for i := 0; i < n; i++ {
		v, err := strconv.Atoi(tokens[i+1])
		if err != nil {
			return fmt.Errorf("parse error: %v", err)
		}
		ans[i] = v
	}
	expect := expectedRowMins(mat)
	for i := 0; i < n; i++ {
		if ans[i] != expect[i] {
			return fmt.Errorf("row %d expected %d got %d", i+1, expect[i], ans[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 1; t <= 100; t++ {
		mat := generateMatrix(rng)
		if err := runCase(bin, mat); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", t, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
