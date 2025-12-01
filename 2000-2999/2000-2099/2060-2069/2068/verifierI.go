package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type operation struct {
	t int
	r int
	c int
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runProgram(bin string, input []byte) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("timeout")
		}
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return out.String(), nil
}

func parseOutput(data string) (ok bool, dir byte, ops []operation, err error) {
	lines := strings.Fields(data)
	if len(lines) == 0 {
		return false, 0, nil, fmt.Errorf("empty output")
	}
	if lines[0] == "NO" {
		return false, 0, nil, nil
	}
	if lines[0] != "YES" {
		return false, 0, nil, fmt.Errorf("first token must be YES/NO")
	}
	if len(lines) < 3 {
		return false, 0, nil, fmt.Errorf("incomplete YES output")
	}
	if len(lines[1]) != 1 || !strings.Contains("UDLR", lines[1]) {
		return false, 0, nil, fmt.Errorf("invalid direction")
	}
	dir = lines[1][0]
	k, err := strconv.Atoi(lines[2])
	if err != nil || k < 0 {
		return false, 0, nil, fmt.Errorf("invalid k")
	}
	if len(lines) != 3+3*k {
		return false, 0, nil, fmt.Errorf("expected %d operation tokens, got %d", 3*k, len(lines)-3)
	}
	ops = make([]operation, k)
	for i := 0; i < k; i++ {
		t, err1 := strconv.Atoi(lines[3+3*i])
		r, err2 := strconv.Atoi(lines[4+3*i])
		c, err3 := strconv.Atoi(lines[5+3*i])
		if err1 != nil || err2 != nil || err3 != nil {
			return false, 0, nil, fmt.Errorf("invalid operation #%d", i+1)
		}
		ops[i] = operation{t: t, r: r - 1, c: c - 1}
	}
	return true, dir, ops, nil
}

func findStart(grid [][]byte) (int, int) {
	for i := range grid {
		for j, ch := range grid[i] {
			if ch == 'S' {
				return i, j
			}
		}
	}
	return -1, -1
}

func simulate(grid [][]byte, dir byte, ops []operation) error {
	h := len(grid)
	w := len(grid[0])
	sr, sc := findStart(grid)
	if sr == -1 {
		return fmt.Errorf("start not found")
	}
	// Validate operations
	seen := make(map[[2]int]bool)
	earliest := make(map[[2]int]int)
	prevT := -1
	for i, op := range ops {
		if op.t < 0 || op.t > 10_000_000 {
			return fmt.Errorf("operation %d time out of range", i+1)
		}
		if op.t < prevT {
			return fmt.Errorf("operations not nondecreasing")
		}
		prevT = op.t
		if op.r < 0 || op.r >= h || op.c < 0 || op.c >= w {
			return fmt.Errorf("operation %d position out of grid", i+1)
		}
		key := [2]int{op.r, op.c}
		if seen[key] {
			return fmt.Errorf("duplicate operation on cell %d %d", op.r+1, op.c+1)
		}
		seen[key] = true
		if grid[op.r][op.c] != '/' && grid[op.r][op.c] != '\\' {
			return fmt.Errorf("operation %d targets non-oblique cell", i+1)
		}
		earliest[key] = op.t
	}

	// Direction deltas
	dr := []int{-1, 0, 1, 0}
	dc := []int{0, 1, 0, -1}
	dirMap := map[byte]int{'U': 0, 'R': 1, 'D': 2, 'L': 3}
	d, ok := dirMap[dir]
	if !ok {
		return fmt.Errorf("invalid direction %c", dir)
	}

	timeNow := 0
	opIdx := 0

	applyOps := func(limit int) {
		for opIdx < len(ops) && ops[opIdx].t <= limit {
			r := ops[opIdx].r
			c := ops[opIdx].c
			grid[r][c] = '.'
			opIdx++
		}
	}

	type state struct {
		r   int
		c   int
		dir int
		op  int
	}
	visited := make(map[state]bool)

	for steps := 0; steps < 5_000_000; steps++ {
		applyOps(timeNow)
		st := state{sr, sc, d, opIdx}
		if visited[st] {
			return fmt.Errorf("ball is in a loop without escaping")
		}
		visited[st] = true

		nr := sr + dr[d]
		nc := sc + dc[d]
		newTime := timeNow + 1

		// escape
		if nr < 0 || nr >= h || nc < 0 || nc >= w {
			return nil
		}

		cell := grid[nr][nc]
		if cell == '/' || cell == '\\' {
			if t, ok := earliest[[2]int{nr, nc}]; ok && t <= newTime {
				cell = '.'
			}
		}

		switch cell {
		case '#':
			d ^= 2 // bounce, stay in current cell
		case '.':
			sr, sc = nr, nc
		case '/':
			sr, sc = nr, nc
			switch d {
			case 0:
				d = 1
			case 1:
				d = 0
			case 2:
				d = 3
			default:
				d = 2
			}
		case '\\':
			sr, sc = nr, nc
			switch d {
			case 0:
				d = 3
			case 3:
				d = 0
			case 2:
				d = 1
			default:
				d = 2
			}
		default:
			return fmt.Errorf("invalid cell character %c", cell)
		}
		timeNow = newTime
	}
	return fmt.Errorf("simulation step limit exceeded without escape")
}

func parseGrid(input []byte) ([][]byte, error) {
	in := bufio.NewReader(bytes.NewReader(input))
	var h, w int
	if _, err := fmt.Fscan(in, &h, &w); err != nil {
		return nil, fmt.Errorf("failed to read h w: %v", err)
	}
	grid := make([][]byte, h)
	for i := 0; i < h; i++ {
		var line string
		if _, err := fmt.Fscan(in, &line); err != nil {
			return nil, fmt.Errorf("failed to read row %d: %v", i+1, err)
		}
		if len(line) != w {
			return nil, fmt.Errorf("row %d length mismatch", i+1)
		}
		grid[i] = []byte(line)
	}
	return grid, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierI.go /path/to/binary")
		return
	}

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read input: %v\n", err)
		os.Exit(1)
	}

	grid, err := parseGrid(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	targetBin, cleanTarget, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanTarget()

	_, src, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(src)
	refPath := filepath.Join(baseDir, "2068I.go")
	refBin, cleanRef, err := buildIfGo(refPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer cleanRef()

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
		os.Exit(1)
	}
	tgtOut, err := runProgram(targetBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target failed: %v\n", err)
		os.Exit(1)
	}

	refYes, _, refOps, err := parseOutput(refOut)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\n", err)
		os.Exit(1)
	}
	tgtYes, tgtDir, tgtOps, err := parseOutput(tgtOut)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse target output: %v\n", err)
		os.Exit(1)
	}

	if refYes != tgtYes {
		fmt.Fprintf(os.Stderr, "feasibility mismatch: reference %v target %v\n", refYes, tgtYes)
		os.Exit(1)
	}
	if !refYes {
		fmt.Println("all tests passed")
		return
	}

	if len(tgtOps) != len(refOps) {
		fmt.Fprintf(os.Stderr, "operation count mismatch: expected %d got %d\n", len(refOps), len(tgtOps))
		os.Exit(1)
	}

	gridCopy := make([][]byte, len(grid))
	for i := range grid {
		gridCopy[i] = append([]byte(nil), grid[i]...)
	}
	if err := simulate(gridCopy, tgtDir, tgtOps); err != nil {
		fmt.Fprintf(os.Stderr, "invalid target output: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("all tests passed")
}
