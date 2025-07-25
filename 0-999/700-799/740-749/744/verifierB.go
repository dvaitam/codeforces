package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
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
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		return err
	}

	w := bufio.NewWriter(stdin)
	r := bufio.NewReader(stdout)
	fmt.Fprintln(w, n)
	w.Flush()
	queries := 0
	for {
		var token string
		if _, err := fmt.Fscan(r, &token); err != nil {
			return fmt.Errorf("read error: %v stderr:%s", err, stderr.String())
		}
		if token == "?" {
			queries++
			if queries > 20 {
				return fmt.Errorf("too many queries")
			}
			var k int
			fmt.Fscan(r, &k)
			idx := make([]int, k)
			for i := 0; i < k; i++ {
				fmt.Fscan(r, &idx[i])
			}
			res := make([]int, n)
			for i := 0; i < n; i++ {
				best := math.MaxInt32
				for _, id := range idx {
					id--
					if mat[i][id] < best {
						best = mat[i][id]
					}
				}
				res[i] = best
			}
			for i := 0; i < n; i++ {
				if i > 0 {
					fmt.Fprint(w, " ")
				}
				fmt.Fprint(w, res[i])
			}
			fmt.Fprint(w, "\n")
			w.Flush()
		} else if token == "-1" {
			ans := make([]int, n)
			for i := 0; i < n; i++ {
				fmt.Fscan(r, &ans[i])
			}
			expect := expectedRowMins(mat)
			for i := 0; i < n; i++ {
				if ans[i] != expect[i] {
					return fmt.Errorf("row %d expected %d got %d", i+1, expect[i], ans[i])
				}
			}
			stdin.Close()
			if err := cmd.Wait(); err != nil {
				return fmt.Errorf("runtime error: %v stderr:%s", err, stderr.String())
			}
			return nil
		} else {
			return fmt.Errorf("unexpected token %s", token)
		}
	}
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
