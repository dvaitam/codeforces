package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"
)

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func expected(n, m int, grid [][]int, A, B, C int) string {
	target := []int{A, B, C}
	sort.Ints(target)
	rowSum := make([]int, n)
	for i := 0; i < n; i++ {
		sum := 0
		for j := 0; j < m; j++ {
			sum += grid[i][j]
		}
		rowSum[i] = sum
	}
	prefRow := make([]int, n+1)
	for i := 0; i < n; i++ {
		prefRow[i+1] = prefRow[i] + rowSum[i]
	}
	colSum := make([]int, m)
	for j := 0; j < m; j++ {
		sum := 0
		for i := 0; i < n; i++ {
			sum += grid[i][j]
		}
		colSum[j] = sum
	}
	prefCol := make([]int, m+1)
	for j := 0; j < m; j++ {
		prefCol[j+1] = prefCol[j] + colSum[j]
	}
	var ways int
	if n >= 3 {
		for c1 := 1; c1 <= n-2; c1++ {
			for c2 := c1 + 1; c2 <= n-1; c2++ {
				s1 := prefRow[c1]
				s2 := prefRow[c2] - prefRow[c1]
				s3 := prefRow[n] - prefRow[c2]
				parts := []int{s1, s2, s3}
				sort.Ints(parts)
				if parts[0] == target[0] && parts[1] == target[1] && parts[2] == target[2] {
					ways++
				}
			}
		}
	}
	if m >= 3 {
		for c1 := 1; c1 <= m-2; c1++ {
			for c2 := c1 + 1; c2 <= m-1; c2++ {
				s1 := prefCol[c1]
				s2 := prefCol[c2] - prefCol[c1]
				s3 := prefCol[m] - prefCol[c2]
				parts := []int{s1, s2, s3}
				sort.Ints(parts)
				if parts[0] == target[0] && parts[1] == target[1] && parts[2] == target[2] {
					ways++
				}
			}
		}
	}
	return fmt.Sprintf("%d\n", ways)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(3) + 3
	m := rng.Intn(3) + 3
	grid := make([][]int, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]int, m)
		for j := 0; j < m; j++ {
			grid[i][j] = rng.Intn(5)
		}
	}
	A := rng.Intn(20)
	B := rng.Intn(20)
	C := rng.Intn(20)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", grid[i][j])
		}
		sb.WriteByte('\n')
	}
	fmt.Fprintf(&sb, "%d %d %d\n", A, B, C)
	return sb.String(), expected(n, m, grid, A, B, C)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: go run verifierD.go /path/to/binary\n")
		os.Exit(1)
	}
	candidate := os.Args[1]

	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	ref := filepath.Join(dir, "refD")
	cmd := exec.Command("go", "build", "-o", ref, filepath.Join(dir, "120D.go"))
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference solution: %v\n%s", err, out)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		input, expect := generateCase(rng)
		candOut, cErr := runBinary(candidate, input)
		refOut, rErr := runBinary(ref, input)
		if cErr != nil {
			fmt.Fprintf(os.Stderr, "test %d: candidate error: %v\n", t+1, cErr)
			os.Exit(1)
		}
		if rErr != nil {
			fmt.Fprintf(os.Stderr, "test %d: reference error: %v\n", t+1, rErr)
			os.Exit(1)
		}
		if strings.TrimSpace(candOut) != strings.TrimSpace(refOut) {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%sexpected:%sactual:%s\n", t+1, input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
