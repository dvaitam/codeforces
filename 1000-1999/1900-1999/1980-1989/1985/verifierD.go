package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func runExe(path string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildRef() (string, error) {
	refSrc := os.Getenv("REFERENCE_SOURCE_PATH")
	if refSrc == "" {
		refSrc = "1985D.go"
	}
	ref := filepath.Join(os.TempDir(), "refD.bin")
	cmd := exec.Command("go", "build", "-o", ref, refSrc)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func buildCase(grid []string) []byte {
	n := len(grid)
	m := len(grid[0])
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for _, row := range grid {
		sb.WriteString(row)
		sb.WriteByte('\n')
	}
	return []byte(sb.String())
}

func genManhattanCircle(rng *rand.Rand) []byte {
	// Generate a valid Manhattan circle on a grid
	// Choose center (h,k) and radius r, then build the grid
	r := rng.Intn(4) + 1
	// Grid must be large enough to contain the circle
	// The circle spans from h-r+1 to h+r-1 in rows and k-r+1 to k+r-1 in cols
	// Add some padding
	minN := 2*r - 1
	minM := 2*r - 1
	n := minN + rng.Intn(3)
	m := minM + rng.Intn(3)
	// Choose center so the circle fits
	h := rng.Intn(n-minN+1) + (r - 1) // 0-indexed row
	k := rng.Intn(m-minM+1) + (r - 1) // 0-indexed col

	grid := make([]string, n)
	for i := 0; i < n; i++ {
		row := make([]byte, m)
		for j := 0; j < m; j++ {
			dist := abs(i-h) + abs(j-k)
			if dist < r {
				row[j] = '#'
			} else {
				row[j] = '.'
			}
		}
		grid[i] = string(row)
	}
	return buildCase(grid)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func genTests() [][]byte {
	rng := rand.New(rand.NewSource(4))
	tests := [][]byte{
		buildCase([]string{"#"}),
	}
	for len(tests) < 100 {
		tests = append(tests, genManhattanCircle(rng))
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := genTests()
	for i, tc := range tests {
		exp, err := runExe(ref, tc)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n%s", i+1, err, exp)
			os.Exit(1)
		}
		got, err := runExe(bin, tc)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n%s", i+1, err, got)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, string(tc), exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
