package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// 475C: Kamal-ol-molk's Painting
// Given n*m grid of 'X' and '.', determine if it can be painted by an a*b
// rectangular brush starting somewhere and moving only right or down (one cell
// at a time, staying inside the frame). Output minimum a*b, or -1.
//
// Standalone validator: we verify the candidate answer by checking all possible
// (a,b) brush sizes from smallest area upward.

func solve(n, m int, grid []string) int {
	// Precompute 2D prefix sums of X cells
	// sum[i][j] = number of X in grid[0..i-1][0..j-1]
	sum := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		sum[i] = make([]int, m+1)
	}
	totalX := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			val := 0
			if grid[i][j] == 'X' {
				val = 1
				totalX++
			}
			sum[i+1][j+1] = val + sum[i][j+1] + sum[i+1][j] - sum[i][j]
		}
	}
	if totalX == 0 {
		return -1
	}

	rectSum := func(r1, c1, r2, c2 int) int {
		// sum of grid[r1..r2-1][c1..c2-1]
		if r1 >= r2 || c1 >= c2 {
			return 0
		}
		return sum[r2][c2] - sum[r1][c2] - sum[r2][c1] + sum[r1][c1]
	}

	// Try all brush sizes (a, b) ordered by area
	type pair struct{ a, b int }
	var sizes []pair
	for a := 1; a <= n; a++ {
		for b := 1; b <= m; b++ {
			sizes = append(sizes, pair{a, b})
		}
	}
	// sort by area
	for i := 1; i < len(sizes); i++ {
		for j := i; j > 0 && sizes[j-1].a*sizes[j-1].b > sizes[j].a*sizes[j].b; j-- {
			sizes[j-1], sizes[j] = sizes[j], sizes[j-1]
		}
	}

	for _, sz := range sizes {
		a, b := sz.a, sz.b
		if canPaint(n, m, a, b, grid, totalX, rectSum) {
			return a * b
		}
	}
	return -1
}

func canPaint(n, m, a, b int, grid []string, totalX int, rectSum func(int, int, int, int) int) bool {
	// The brush starts at some position (r,c) covering rows [r..r+a-1], cols [c..c+b-1]
	// It moves right or down. The set of painted cells is the union of all brush positions along the path.
	//
	// Key observation: the path of top-left corners forms a staircase from (r0,c0) to (r1,c1)
	// where r1 >= r0 and c1 >= c0, moving only right or down.
	// The union of brush positions forms an "L-shaped" or rectangular region.
	//
	// Actually, the union of a monotone right-down path of an a*b brush is:
	// For each row i in the path of top-left corners, the brush covers columns [c_i, c_i+b-1]
	// where c_i is non-decreasing. Also rows covered are [r_start..r_end+a-1].
	//
	// Simpler approach for small grids: BFS/simulation.
	// Since grids are at most 4x4, we can do a direct check.

	// Find bounding box of X cells
	minR, maxR, minC, maxC := n, -1, m, -1
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == 'X' {
				if i < minR { minR = i }
				if i > maxR { maxR = i }
				if j < minC { minC = j }
				if j > maxC { maxC = j }
			}
		}
	}

	// The brush must be able to reach all X cells
	// Brush height a, width b. Path goes right/down.
	// Try all possible starting positions for the brush
	for r0 := 0; r0+a-1 < n; r0++ {
		for c0 := 0; c0+b-1 < m; c0++ {
			// Check if starting here, we can paint all X cells with right/down moves
			if tryPath(n, m, a, b, r0, c0, grid, totalX, rectSum) {
				return true
			}
		}
	}
	return false
}

func tryPath(n, m, a, b, r0, c0 int, grid []string, totalX int, rectSum func(int, int, int, int) int) bool {
	// BFS: state is (r, c) = top-left of brush. Moves: right (r, c+1) or down (r+1, c).
	// We collect all cells painted. Check if it matches the X pattern exactly.

	type state struct{ r, c int }
	// Use DFS with memoization to find if there's a path that covers all X and only X
	// Actually we need to find ANY path whose union of brush positions = set of X cells.

	// For small grids (up to 4x4), we can enumerate all possible paths.
	// A path is a sequence of R and D moves. Max path length is (n-a) + (m-b).
	// Max moves = (4-1)+(4-1) = 6, so 2^6 = 64 paths max per starting position.

	maxR := n - a  // max row for top-left
	maxC := m - b  // max col for top-left
	if r0 > maxR || c0 > maxC {
		return false
	}

	// Generate all monotone paths from (r0, c0) by going right or down
	// Use BFS to enumerate all reachable endpoints and track painted cells
	// Actually, the union of all brush positions along ANY monotone path from (r0,c0) to (r1,c1)
	// is the same regardless of the specific path taken (it's the union of all positions
	// with r0<=r<=r1, c0<=c<=c1 where the brush fits... no, that's not right).

	// For correctness, enumerate all possible paths using BFS over subsets of painted cells.
	// Since grid is small (max 4x4=16 cells), we can use bitmask.

	totalCells := n * m
	if totalCells > 20 {
		// Shouldn't happen with our test gen (max 4x4)
		return false
	}

	targetMask := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == 'X' {
				targetMask |= 1 << (i*m + j)
			}
		}
	}

	// BFS over (r, c, paintedMask)
	brushMask := func(r, c int) int {
		mask := 0
		for dr := 0; dr < a; dr++ {
			for dc := 0; dc < b; dc++ {
				mask |= 1 << ((r+dr)*m + (c+dc))
			}
		}
		return mask
	}

	initMask := brushMask(r0, c0)
	// If initial brush covers any non-X cell... no, the brush paints cells.
	// The painted cells must equal exactly the X cells. So the union of all brush
	// positions must be a subset of X cells? No - the brush CREATES the X pattern.
	// So we need: union of brush positions = X cells.
	// That means we need to find a path where the union equals targetMask.

	type bfsState struct {
		r, c int
		mask int
	}

	visited := make(map[bfsState]bool)
	queue := []bfsState{{r0, c0, initMask}}
	visited[bfsState{r0, c0, initMask}] = true

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		// Check if current mask matches target
		if cur.mask == targetMask {
			return true
		}
		// If we've already painted cells outside target, prune
		if cur.mask & ^targetMask != 0 {
			continue
		}

		// Try move right
		nc := cur.c + 1
		if nc+b-1 < m {
			nm := cur.mask | brushMask(cur.r, nc)
			s := bfsState{cur.r, nc, nm}
			if !visited[s] {
				visited[s] = true
				queue = append(queue, s)
			}
		}
		// Try move down
		nr := cur.r + 1
		if nr+a-1 < n {
			nm := cur.mask | brushMask(nr, cur.c)
			s := bfsState{nr, cur.c, nm}
			if !visited[s] {
				visited[s] = true
				queue = append(queue, s)
			}
		}
	}
	return false
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(4) + 1
	m := rng.Intn(4) + 1
	grid := make([][]byte, n)
	hasX := false
	for i := 0; i < n; i++ {
		grid[i] = make([]byte, m)
		for j := 0; j < m; j++ {
			if rng.Intn(3) == 0 {
				grid[i][j] = 'X'
				hasX = true
			} else {
				grid[i][j] = '.'
			}
		}
	}
	if !hasX {
		grid[0][0] = 'X'
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 0; i < n; i++ {
		sb.Write(grid[i])
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runProg(exe, input string) (string, error) {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveFromInput(input string) int {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	parts := strings.Fields(lines[0])
	n, _ := strconv.Atoi(parts[0])
	m, _ := strconv.Atoi(parts[1])
	grid := make([]string, n)
	for i := 0; i < n; i++ {
		grid[i] = strings.TrimSpace(lines[i+1])
	}
	return solve(n, m, grid)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	// Fixed test cases including the known failing one
	fixedTests := []string{
		"2 3\nXXX\nX..\n",
		"1 1\nX\n",
		"2 2\nXX\nXX\n",
		"1 3\nXXX\n",
		"3 1\nX\nX\nX\n",
		"2 2\nX.\n.X\n",
		"3 3\nXXX\nXXX\nXXX\n",
		"2 3\nXXX\nXXX\n",
		"3 3\nXX.\nXXX\n.XX\n",
		"2 2\n.X\nX.\n",
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 80; i++ {
		fixedTests = append(fixedTests, generateCase(rng))
	}

	for i, input := range fixedTests {
		exp := solveFromInput(input)
		got, err := runProg(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		expStr := strconv.Itoa(exp)
		if got != expStr {
			fmt.Fprintf(os.Stderr, "case %d mismatch\nexpected:%s\n got:%s\ninput:\n%s", i+1, expStr, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
