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

const (
	refSourceA = "193A.go"
	refBinaryA = "ref193A.bin"
	totalTests = 150
)

type testCase struct {
	n, m int
	grid []string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refPath, err := buildReference()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refPath)

	tests := generateTests()
	for i, tc := range tests {
		input := formatInput(tc)

		refOut, err := runProgram(refPath, input)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", i+1, err)
			printInput(input)
			os.Exit(1)
		}
		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			printInput(input)
			os.Exit(1)
		}

		refVal, err := parseOutput(refOut)
		if err != nil {
			fmt.Printf("failed to parse reference output on test %d: %v\noutput:\n%s\n", i+1, err, refOut)
			printInput(input)
			os.Exit(1)
		}
		candVal, err := parseOutput(candOut)
		if err != nil {
			fmt.Printf("failed to parse candidate output on test %d: %v\noutput:\n%s\n", i+1, err, candOut)
			printInput(input)
			os.Exit(1)
		}

		if refVal != candVal {
			fmt.Printf("test %d failed: expected %d, got %d\n", i+1, refVal, candVal)
			printInput(input)
			fmt.Println("Reference output:")
			fmt.Println(refOut)
			fmt.Println("Candidate output:")
			fmt.Println(candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinaryA, refSourceA)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinaryA), nil
}

func runProgram(path string, input []byte) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func formatInput(tc testCase) []byte {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
	for _, line := range tc.grid {
		sb.WriteString(line)
		sb.WriteByte('\n')
	}
	return []byte(sb.String())
}

func parseOutput(out string) (int, error) {
	fields := strings.Fields(out)
	if len(fields) != 1 {
		return 0, fmt.Errorf("expected 1 integer, got %d", len(fields))
	}
	val, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q: %v", fields[0], err)
	}
	return val, nil
}

func generateTests() []testCase {
	tests := []testCase{
		{n: 1, m: 1, grid: []string{"#"}},
		{n: 1, m: 2, grid: []string{"##"}},
		{n: 2, m: 2, grid: []string{"#.", "##"}},
		{n: 3, m: 3, grid: []string{"###", ".#.", "###"}},
		{n: 3, m: 5, grid: []string{"..#..", "..#..", "..#.."}},
		{n: 4, m: 4, grid: []string{"....", ".##.", ".##.", "...."}},
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < totalTests-1 {
		n := rnd.Intn(50) + 1
		m := rnd.Intn(50) + 1
		grid := randomConnectedGrid(rnd, n, m)
		tests = append(tests, testCase{n: n, m: m, grid: grid})
	}

	// Stress test near the limits with a dense grid.
	n, m := 50, 50
	grid := make([]string, n)
	for i := 0; i < n; i++ {
		grid[i] = strings.Repeat("#", m)
	}
	tests = append(tests, testCase{n: n, m: m, grid: grid})
	return tests
}

func randomConnectedGrid(rnd *rand.Rand, n, m int) []string {
	grid := make([][]byte, n)
	visited := make([][]bool, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]byte, m)
		visited[i] = make([]bool, m)
		for j := 0; j < m; j++ {
			grid[i][j] = '.'
		}
	}

	startX := rnd.Intn(n)
	startY := rnd.Intn(m)
	grid[startX][startY] = '#'
	visited[startX][startY] = true
	cells := 1
	target := rnd.Intn(n*m) + 1

	type cell struct{ x, y int }
	active := []cell{{startX, startY}}

	dirs := []cell{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for cells < target {
		base := active[rnd.Intn(len(active))]
		var options []cell
		for _, d := range dirs {
			nx, ny := base.x+d.x, base.y+d.y
			if nx >= 0 && nx < n && ny >= 0 && ny < m && !visited[nx][ny] {
				options = append(options, cell{nx, ny})
			}
		}
		if len(options) == 0 {
			full := true
			for _, c := range active {
				for _, d := range dirs {
					nx, ny := c.x+d.x, c.y+d.y
					if nx >= 0 && nx < n && ny >= 0 && ny < m && !visited[nx][ny] {
						full = false
						break
					}
				}
				if !full {
					break
				}
			}
			if full {
				break
			}
			continue
		}
		next := options[rnd.Intn(len(options))]
		grid[next.x][next.y] = '#'
		visited[next.x][next.y] = true
		active = append(active, next)
		cells++
	}

	res := make([]string, n)
	for i := 0; i < n; i++ {
		res[i] = string(grid[i])
	}
	return res
}

func printInput(in []byte) {
	fmt.Println("Input used:")
	fmt.Println(string(in))
}
