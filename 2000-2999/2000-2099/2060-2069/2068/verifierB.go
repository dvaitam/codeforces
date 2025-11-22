package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const maxHW = 2025

type testCase struct {
	name string
	k    int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary-or-source")
		os.Exit(1)
	}

	candidate, cleanup, err := prepareCandidateBinary(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()
	for idx, tc := range tests {
		input := fmt.Sprintf("%d\n", tc.k)
		output, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d (%s): %v\n", idx+1, tc.name, err)
			continueCleanupAndExit(cleanup, 1)
		}
		if err := validateOutput(tc.k, output); err != nil {
			fmt.Fprintf(os.Stderr, "invalid output on test %d (%s): %v\n", idx+1, tc.name, err)
			fmt.Fprintln(os.Stderr, previewInput(input))
			continueCleanupAndExit(cleanup, 1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func prepareCandidateBinary(path string) (string, func(), error) {
	if !strings.HasSuffix(path, ".go") {
		return path, func() {}, nil
	}
	tmpDir, err := os.MkdirTemp("", "verifier2068B-cand")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	bin := filepath.Join(tmpDir, "candidate2068B")
	cmd := exec.Command("go", "build", "-o", bin, path)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build candidate: %v\n%s", err, out.String())
	}
	cleanup := func() { _ = os.RemoveAll(tmpDir) }
	return bin, cleanup, nil
}

func runProgram(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func validateOutput(k int64, out string) error {
	h, w, grid, err := parseGrid(out)
	if err != nil {
		return err
	}
	bitsRows := toBitRows(grid, w)
	rects := countRectangularWalks(bitsRows, h, w)
	if rects != k {
		return fmt.Errorf("expected %d rectangular walks, got %d", k, rects)
	}
	return nil
}

func parseGrid(out string) (int, int, []string, error) {
	reader := bufio.NewReader(strings.NewReader(out))
	var h, w int
	if _, err := fmt.Fscan(reader, &h, &w); err != nil {
		return 0, 0, nil, fmt.Errorf("failed to read grid size: %v", err)
	}
	if h < 1 || w < 1 || h > maxHW || w > maxHW {
		return 0, 0, nil, fmt.Errorf("grid size out of bounds h=%d w=%d", h, w)
	}
	// consume rest of line
	if _, err := reader.ReadString('\n'); err != nil && err != io.EOF {
		return 0, 0, nil, fmt.Errorf("failed to consume first line: %v", err)
	}

	grid := make([]string, h)
	for i := 0; i < h; i++ {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return 0, 0, nil, fmt.Errorf("failed to read row %d: %v", i+1, err)
		}
		line = strings.TrimRight(line, "\r\n")
		if len(line) != w {
			return 0, 0, nil, fmt.Errorf("row %d has length %d, expected %d", i+1, len(line), w)
		}
		for j, ch := range line {
			if ch != '#' && ch != '.' {
				return 0, 0, nil, fmt.Errorf("invalid character at row %d col %d: %q", i+1, j+1, ch)
			}
		}
		grid[i] = line
	}
	return h, w, grid, nil
}

func toBitRows(grid []string, w int) [][]uint64 {
	words := (w + 63) / 64
	rows := make([][]uint64, len(grid))
	for i, line := range grid {
		row := make([]uint64, words)
		for c, ch := range line {
			if ch == '#' {
				row[c>>6] |= 1 << uint(c&63)
			}
		}
		rows[i] = row
	}
	return rows
}

func countRectangularWalks(rows [][]uint64, h, w int) int64 {
	if h < 2 || w < 2 {
		return 0
	}
	words := (w + 63) / 64
	colAll := make([]uint64, words)
	both := make([]uint64, words)
	var total int64
	for top := 0; top < h-1; top++ {
		copy(colAll, rows[top])
		for bottom := top + 1; bottom < h; bottom++ {
			for i := 0; i < words; i++ {
				colAll[i] &= rows[bottom][i]
				both[i] = rows[top][i] & rows[bottom][i]
			}
			edgeCount := popcountMasked(colAll, w)
			if edgeCount < 2 {
				continue
			}
			if popcountMasked(both, w) == w {
				total += choose2(int64(edgeCount))
				continue
			}
			total += countRuns(both, colAll, w)
		}
	}
	return total
}

func countRuns(both, edges []uint64, w int) int64 {
	var res int64
	for c := 0; c < w; {
		if !bitAt(both, c) {
			c++
			continue
		}
		edgeCnt := 0
		for c < w && bitAt(both, c) {
			if bitAt(edges, c) {
				edgeCnt++
			}
			c++
		}
		if edgeCnt >= 2 {
			res += choose2(int64(edgeCnt))
		}
	}
	return res
}

func popcountMasked(bitsRow []uint64, w int) int {
	if len(bitsRow) == 0 {
		return 0
	}
	total := 0
	lastIdx := len(bitsRow) - 1
	for i := 0; i < lastIdx; i++ {
		total += bits.OnesCount64(bitsRow[i])
	}
	last := bitsRow[lastIdx]
	if rem := w & 63; rem != 0 {
		last &= (uint64(1) << uint(rem)) - 1
	}
	total += bits.OnesCount64(last)
	return total
}

func bitAt(arr []uint64, idx int) bool {
	return (arr[idx>>6]>>uint(idx&63))&1 == 1
}

func choose2(v int64) int64 {
	return v * (v - 1) / 2
}

func tri(x int64) int64 {
	return x * (x - 1) / 2
}

func buildTests() []testCase {
	var tests []testCase
	tests = append(tests,
		testCase{name: "zero", k: 0},
		testCase{name: "one", k: 1},
		testCase{name: "two", k: 2},
		testCase{name: "five", k: 5},
		testCase{name: "small_randomish", k: 37},
		testCase{name: "mid", k: 1_000_000},
	)

	maxK := tri(maxHW) * tri(maxHW)
	tests = append(tests,
		testCase{name: "near_limit_minus", k: maxK - 1},
		testCase{name: "limit", k: maxK},
	)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 40; i++ {
		val := rng.Int63n(maxK + 1)
		tests = append(tests, testCase{
			name: fmt.Sprintf("rand_%d", i+1),
			k:    val,
		})
	}
	return tests
}

func previewInput(in string) string {
	if len(in) <= 200 {
		return "input:\n" + in
	}
	return "input (truncated):\n" + in[:200] + "..."
}

func continueCleanupAndExit(cleanup func(), code int) {
	cleanup()
	os.Exit(code)
}
