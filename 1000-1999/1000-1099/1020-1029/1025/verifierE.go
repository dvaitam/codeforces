package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func buildOracle() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "1025E.go")
	bin := filepath.Join(os.TempDir(), "oracle1025E.bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return bin, nil
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(r *rand.Rand) string {
	n := r.Intn(5) + 4 // Generates n in [4, 8]
	m := r.Intn(n) + 1
	type cell struct{ x, y int }
	used := make(map[[2]int]bool)
	cells := make([]cell, 0, n*n)
	for x := 1; x <= n; x++ {
		for y := 1; y <= n; y++ {
			cells = append(cells, cell{x, y})
		}
	}
	r.Shuffle(len(cells), func(i, j int) { cells[i], cells[j] = cells[j], cells[i] })
	start := make([]cell, m)
	target := make([]cell, m)
	idx := 0
	for i := 0; i < m; i++ {
		start[i] = cells[idx]
		used[[2]int{cells[idx].x, cells[idx].y}] = true
		idx++
	}
	for i := 0; i < m; i++ {
		// choose unused cell for target
		for used[[2]int{cells[idx].x, cells[idx].y}] {
			idx++
		}
		target[i] = cells[idx]
		used[[2]int{cells[idx].x, cells[idx].y}] = true
		idx++
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for _, c := range start {
		fmt.Fprintf(&sb, "%d %d\n", c.x, c.y)
	}
	for _, c := range target {
		fmt.Fprintf(&sb, "%d %d\n", c.x, c.y)
	}
	return sb.String()
}

func checkSolution(input, output string) error {
	// Parse input
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanWords)
	nextInt := func() int {
		scanner.Scan()
		v, _ := strconv.Atoi(scanner.Text())
		return v
	}
	n := nextInt()
	m := nextInt()
	type Point struct{ x, y int }
	start := make([]Point, m)
	target := make([]Point, m)
	for i := 0; i < m; i++ {
		start[i] = Point{nextInt(), nextInt()}
	}
	for i := 0; i < m; i++ {
		target[i] = Point{nextInt(), nextInt()}
	}

	// Parse output
	outScanner := bufio.NewScanner(strings.NewReader(output))
	outScanner.Split(bufio.ScanWords)
	if !outScanner.Scan() {
		return fmt.Errorf("empty output")
	}
	k, err := strconv.Atoi(outScanner.Text())
	if err != nil {
		return fmt.Errorf("invalid number of moves: %v", err)
	}
	if k > 10800 {
		return fmt.Errorf("too many moves: %d > 10800", k)
	}

	// Simulate
	grid := make([][]int, n+1)
	for i := range grid {
		grid[i] = make([]int, n+1)
	}
	for i, p := range start {
		grid[p.x][p.y] = i + 1
	}

	for i := 0; i < k; i++ {
		if !outScanner.Scan() {
			return fmt.Errorf("expected more moves")
		}
		x1, _ := strconv.Atoi(outScanner.Text())
		if !outScanner.Scan() { return fmt.Errorf("incomplete move") }
		y1, _ := strconv.Atoi(outScanner.Text())
		if !outScanner.Scan() { return fmt.Errorf("incomplete move") }
		x2, _ := strconv.Atoi(outScanner.Text())
		if !outScanner.Scan() { return fmt.Errorf("incomplete move") }
		y2, _ := strconv.Atoi(outScanner.Text())

		if x1 < 1 || x1 > n || y1 < 1 || y1 > n {
			return fmt.Errorf("move %d: start (%d,%d) out of bounds", i+1, x1, y1)
		}
		if x2 < 1 || x2 > n || y2 < 1 || y2 > n {
			return fmt.Errorf("move %d: end (%d,%d) out of bounds", i+1, x2, y2)
		}
		if abs(x1-x2)+abs(y1-y2) != 1 {
			return fmt.Errorf("move %d: not adjacent (%d,%d)->(%d,%d)", i+1, x1, y1, x2, y2)
		}
		if grid[x1][y1] == 0 {
			return fmt.Errorf("move %d: no cube at (%d,%d)", i+1, x1, y1)
		}
		if grid[x2][y2] != 0 {
			return fmt.Errorf("move %d: destination (%d,%d) occupied", i+1, x2, y2)
		}
		grid[x2][y2] = grid[x1][y1]
		grid[x1][y1] = 0
	}

	// Check final state
	for i, p := range target {
		if grid[p.x][p.y] == 0 {
			return fmt.Errorf("target (%d,%d) empty", p.x, p.y)
		}
		// Note: The problem states "place each cube on its place", implying identity of cubes matters.
		// However, the input says "all cubes have distinct colors", but doesn't explicitly link start[i] to target[i].
		// Wait, "Each of the next m lines contains... initial positions... The next m lines describe the designated places... in the same format and order."
		// This implies start[i] must move to target[i].
		// BUT, usually in such problems if cubes are "colored" and "distinct", then yes, cube i starts at start[i] and must end at target[i].
		// Let's verify the oracle actually respects this.
		// The oracle logic is hidden, but the problem statement says "all cubes have distinct colors".
		// Let's assume cube at start[i] must end at target[i].
		if grid[p.x][p.y] != i+1 {
			return fmt.Errorf("target (%d,%d) has cube %d, expected %d", p.x, p.y, grid[p.x][p.y], i+1)
		}
	}
	return nil
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	userBin := os.Args[1]
	
	// We don't need the oracle if we verify validity directly.
	// The oracle was used to check string equality, which is wrong for problems with multiple solutions.
	// The verifier itself was flawed because it compared exact output string with oracle.
	
	r := rand.New(rand.NewSource(1))
	for i := 0; i < 100; i++ {
		input := genCase(r)
		got, err := run(userBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if err := checkSolution(input, got); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%sgot:\n%s\n", i+1, err, input, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", 100)
}