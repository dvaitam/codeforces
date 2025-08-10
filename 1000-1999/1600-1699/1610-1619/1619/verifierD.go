package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("time limit")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

func check(matrix [][]int, m, n int, x int) bool {
	friendOk := make([]bool, n)
	hasDouble := false
	for i := 0; i < m; i++ {
		cnt := 0
		row := matrix[i]
		for j := 0; j < n; j++ {
			if row[j] >= x {
				if !friendOk[j] {
					friendOk[j] = true
				}
				cnt++
			}
		}
		if cnt >= 2 {
			hasDouble = true
		}
	}
	if !hasDouble {
		return false
	}
	for j := 0; j < n; j++ {
		if !friendOk[j] {
			return false
		}
	}
	return true
}

func solveCase(matrix [][]int, m, n int) int {
	low := 1
	high := 1000000000
	for low < high {
		mid := (low + high + 1) / 2
		if check(matrix, m, n, mid) {
			low = mid
		} else {
			high = mid - 1
		}
	}
	return low
}

func generateCase(rng *rand.Rand) (string, string) {
	m := rng.Intn(5) + 2
	n := rng.Intn(5) + 2
	matrix := make([][]int, m)
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n\n%d %d\n", m, n)
	for i := 0; i < m; i++ {
		matrix[i] = make([]int, n)
		for j := 0; j < n; j++ {
			val := rng.Intn(1000000000) + 1
			matrix[i][j] = val
			fmt.Fprintf(&sb, "%d ", val)
		}
		sb.WriteByte('\n')
	}
	ans := solveCase(matrix, m, n)
	return sb.String(), fmt.Sprint(ans)
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
		out, err := runBinary(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
