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
	"strconv"
	"strings"
	"time"
)

const (
	refSource        = "1142C.go"
	tempOraclePrefix = "oracle-1142C-"
	maxRandomPoints  = 2000
	maxCoord         = 1_000_000
)

type point struct {
	x int
	y int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build oracle: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := make([][]point, 0)
	tests = append(tests, fixedTests()...)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomTests(rng, 50)...)

	for idx, pts := range tests {
		input := formatInput(pts)
		exp, err := runProgram(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle runtime error on test %d: %v\n", idx+1, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}
		got, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\n", idx+1, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %s got %s\n", idx+1, strings.TrimSpace(exp), strings.TrimSpace(got))
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("failed to determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", tempOraclePrefix)
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleC")
	cmd := exec.Command("go", "build", "-o", outPath, refSource)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func formatInput(points []point) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(points)))
	sb.WriteByte('\n')
	for _, p := range points {
		sb.WriteString(strconv.Itoa(p.x))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(p.y))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func fixedTests() [][]point {
	return [][]point{
		{{0, 0}},
		{{0, 0}, {1, 1}},
		{{0, 0}, {0, 1}, {1, 2}},
		{{-2, 4}, {-1, 1}, {0, 0}, {1, 1}, {2, 4}},
		{{0, 0}, {1, 1}, {2, 0}, {3, 1}, {4, 0}},
	}
}

func randomTests(rng *rand.Rand, count int) [][]point {
	tests := make([][]point, 0, count)
	for t := 0; t < count; t++ {
		n := rng.Intn(maxRandomPoints-1) + 1
		if n > maxRandomPoints {
			n = maxRandomPoints
		}
		points := make([]point, 0, n)
		used := make(map[point]struct{})
		for len(points) < n {
			p := point{
				x: rng.Intn(2*maxCoord+1) - maxCoord,
				y: rng.Intn(2*maxCoord+1) - maxCoord,
			}
			if _, ok := used[p]; ok {
				continue
			}
			used[p] = struct{}{}
			points = append(points, p)
		}
		tests = append(tests, points)
	}
	return tests
}
