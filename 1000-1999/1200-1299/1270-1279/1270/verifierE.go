package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type point struct {
	x int
	y int
}

type testCase struct {
	pts []point
}

func buildReferenceBinary() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("unable to resolve verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "ref-1270E-")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "ref_1270E")
	cmd := exec.Command("go", "build", "-o", binPath, "1270E.go")
	cmd.Dir = dir
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, out.String())
	}
	cleanup := func() { os.RemoveAll(tmpDir) }
	return binPath, cleanup, nil
}

func inputString(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tc.pts)))
	for _, p := range tc.pts {
		sb.WriteString(fmt.Sprintf("%d %d\n", p.x, p.y))
	}
	return sb.String()
}

func runProgram(path, input string) (string, error) {
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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

func parseOutput(output string, n int) ([]int, error) {
	fields := strings.Fields(output)
	if len(fields) < 1 {
		return nil, fmt.Errorf("output is empty")
	}
	a := 0
	if _, err := fmt.Sscan(fields[0], &a); err != nil {
		return nil, fmt.Errorf("failed to read group size: %v", err)
	}
	if a <= 0 || a >= n {
		return nil, fmt.Errorf("invalid group size %d", a)
	}
	if len(fields) < 1+a {
		return nil, fmt.Errorf("expected %d indices, got %d", a, len(fields)-1)
	}
	group := make([]int, a)
	seen := make(map[int]bool)
	for i := 0; i < a; i++ {
		var idx int
		if _, err := fmt.Sscan(fields[1+i], &idx); err != nil {
			return nil, fmt.Errorf("failed to parse index: %v", err)
		}
		if idx < 1 || idx > n {
			return nil, fmt.Errorf("index %d out of range", idx)
		}
		if seen[idx] {
			return nil, fmt.Errorf("duplicate index %d", idx)
		}
		seen[idx] = true
		group[i] = idx - 1
	}
	return group, nil
}

func checkPartition(tc testCase, group []int) error {
	n := len(tc.pts)
	mark := make([]bool, n)
	for _, idx := range group {
		mark[idx] = true
	}
	type key struct{ a, b int64 }
	distYellow := make(map[int64]bool)
	distBlue := make(map[int64]bool)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			dx := int64(tc.pts[i].x - tc.pts[j].x)
			dy := int64(tc.pts[i].y - tc.pts[j].y)
			dist := dx*dx + dy*dy
			if mark[i] == mark[j] {
				if distBlue[dist] {
					return fmt.Errorf("distance %d appears both as yellow and blue", dist)
				}
				distYellow[dist] = true
			} else {
				if distYellow[dist] {
					return fmt.Errorf("distance %d appears both as yellow and blue", dist)
				}
				distBlue[dist] = true
			}
		}
	}
	return nil
}

func deterministicTests() []testCase {
	return []testCase{
		{pts: []point{{0, 0}, {1, 0}}},
		{pts: []point{{0, 0}, {1, 0}, {0, 1}}},
		{pts: []point{{0, 0}, {2, 0}, {0, 2}, {2, 2}}},
		{pts: []point{{0, 0}, {1, 1}, {2, 0}, {3, 1}}},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	for i := 0; i < 120; i++ {
		n := rng.Intn(8) + 2
		pts := make([]point, n)
		used := make(map[[2]int]bool)
		for j := 0; j < n; j++ {
			for {
				x := rng.Intn(2000) - 1000
				y := rng.Intn(2000) - 1000
				if !used[[2]int{x, y}] {
					pts[j] = point{x: x, y: y}
					used[[2]int{x, y}] = true
					break
				}
			}
		}
		tests = append(tests, testCase{pts: pts})
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[len(os.Args)-1]
	if target == "--" {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	refBin, cleanup, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := append(deterministicTests(), randomTests()...)
	for idx, tc := range tests {
		input := inputString(tc)
		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		refGroup, err := parseOutput(refOut, len(tc.pts))
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference parse error on test %d: %v\ninput:\n%soutput:\n%s", idx+1, err, input, refOut)
			os.Exit(1)
		}
		if err := checkPartition(tc, refGroup); err != nil {
			fmt.Fprintf(os.Stderr, "reference validation failed on test %d: %v\ninput:\n%soutput:\n%s", idx+1, err, input, refOut)
			os.Exit(1)
		}

		out, err := runProgram(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d runtime error: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		group, err := parseOutput(out, len(tc.pts))
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d parse error: %v\ninput:\n%soutput:\n%s", idx+1, err, input, out)
			os.Exit(1)
		}
		if err := checkPartition(tc, group); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%soutput:\n%s", idx+1, err, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
