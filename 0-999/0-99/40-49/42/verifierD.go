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

func validateMatrix(mat [][]int) error {
	n := len(mat)
	seen := map[int]bool{}
	base := make([]int, n)
	for i := 0; i < n; i++ {
		if len(mat[i]) != n {
			return fmt.Errorf("row %d length mismatch", i)
		}
		if mat[i][i] != 0 {
			return fmt.Errorf("diagonal not zero")
		}
		if i > 0 {
			base[i] = mat[0][i]
			if base[i] <= 0 || base[i] > 1000 {
				return fmt.Errorf("edge weight out of range")
			}
		}
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if mat[i][j] != mat[j][i] {
				return fmt.Errorf("matrix not symmetric")
			}
			if mat[i][j] <= 0 || mat[i][j] > 1000 {
				return fmt.Errorf("edge weight out of range")
			}
			if seen[mat[i][j]] {
				return fmt.Errorf("weights not distinct")
			}
			seen[mat[i][j]] = true
			if base[i]+base[j] != mat[i][j] {
				return fmt.Errorf("does not satisfy sum property")
			}
		}
	}
	return nil
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(4) + 3 // 3..6
	return fmt.Sprintf("%d\n", n)
}

func runCase(bin, input string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	n := 0
	fmt.Sscan(strings.TrimSpace(input), &n)
	if len(lines) != n {
		return fmt.Errorf("expected %d lines got %d", n, len(lines))
	}
	mat := make([][]int, n)
	for i := 0; i < n; i++ {
		parts := strings.Fields(lines[i])
		if len(parts) != n {
			return fmt.Errorf("line %d expected %d numbers", i+1, n)
		}
		mat[i] = make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Sscan(parts[j], &mat[i][j])
		}
	}
	return validateMatrix(mat)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in := generateCase(rng)
		if err := runCase(bin, in); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
