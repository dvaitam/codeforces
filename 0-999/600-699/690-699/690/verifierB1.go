package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const refSource = "0-999/600-699/690-699/690/690B1.go"

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB1.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}
		expAns, err := parseOutput(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		gotAns, err := parseOutput(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		if gotAns != expAns {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: expected %s got %s\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
				idx+1, tc.name, expAns, gotAns, tc.input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-690B1-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref690B1.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return binPath, cleanup, nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutput(out string) (string, error) {
	ans := strings.TrimSpace(out)
	ans = strings.ToLower(ans)
	if ans == "yes" {
		return "yes", nil
	}
	if ans == "no" {
		return "no", nil
	}
	return "", fmt.Errorf("unexpected output %q", out)
}

func buildTests() []testCase {
	tests := []testCase{
		makeTestCase("single_rect", 5, true),
		makeTestCase("no_rect", 5, false),
		makeTestCase("edge_touching", 6, true),
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomCase(rng, i))
	}
	tests = append(tests, makeMaxCase(true))
	tests = append(tests, makeMaxCase(false))
	return tests
}

func makeTestCase(name string, n int, hasRect bool) testCase {
	grid := make([][]int, n)
	for i := range grid {
		grid[i] = make([]int, n)
	}
	if hasRect {
		x1, x2 := 1, n-1
		y1, y2 := 1, n-1
		fillRect(grid, x1, y1, x2, y2)
	} else {
		// set grid to invalid shape (two rectangles)
		fillRect(grid, 1, 1, n/2, n/2)
		fillRect(grid, n/2, n/2, n-1, n-1)
	}
	return testCase{name: name, input: formatInput(grid)}
}

func fillRect(grid [][]int, x1, y1, x2, y2 int) {
	for i := x1; i < x2; i++ {
		for j := y1; j < y2; j++ {
			grid[j][i] = 4
		}
	}
}

func randomCase(rng *rand.Rand, idx int) testCase {
	n := rng.Intn(46) + 5
	grid := make([][]int, n)
	for i := range grid {
		grid[i] = make([]int, n)
	}
	if rng.Intn(2) == 0 {
		// build valid rectangle
		x1 := rng.Intn(n - 1)
		x2 := rng.Intn(n-x1-1) + x1 + 1
		y1 := rng.Intn(n - 1)
		y2 := rng.Intn(n-y1-1) + y1 + 1
		grid = buildFromRect(n, x1, y1, x2, y2)
	} else {
		// random digits (likely invalid)
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				grid[i][j] = rng.Intn(5)
			}
		}
	}
	return testCase{
		name:  fmt.Sprintf("random_%d", idx+1),
		input: formatInput(grid),
	}
}

func buildFromRect(n, x1, y1, x2, y2 int) [][]int {
	grid := make([][]int, n)
	for i := range grid {
		grid[i] = make([]int, n)
	}
	for y := 0; y < n; y++ {
		for x := 0; x < n; x++ {
			c := 0
			if inside(x, y, x1, y1, x2, y2) {
				c++
			}
			if inside(x+1, y, x1, y1, x2, y2) {
				c++
			}
			if inside(x, y+1, x1, y1, x2, y2) {
				c++
			}
			if inside(x+1, y+1, x1, y1, x2, y2) {
				c++
			}
			grid[y][x] = c
		}
	}
	return grid
}

func inside(x, y, x1, y1, x2, y2 int) bool {
	return x >= x1 && x <= x2 && y >= y1 && y <= y2
}

func formatInput(grid [][]int) string {
	n := len(grid)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := n - 1; i >= 0; i-- {
		for j := 0; j < n; j++ {
			sb.WriteByte(byte('0' + grid[i][j]))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func makeMaxCase(valid bool) testCase {
	n := 50
	var grid [][]int
	if valid {
		grid = buildFromRect(n, 5, 5, 45, 40)
	} else {
		grid = buildFromRect(n, 5, 5, 25, 25)
		another := buildFromRect(n, 10, 30, 40, 45)
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				grid[i][j] = (grid[i][j] + another[i][j]) % 5
			}
		}
	}
	name := "max_valid"
	if !valid {
		name = "max_invalid"
	}
	return testCase{name: name, input: formatInput(grid)}
}
