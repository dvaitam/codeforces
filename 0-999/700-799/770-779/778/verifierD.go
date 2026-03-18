package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func genGrid(rng *rand.Rand, n, m int) []string {
	g := make([][]byte, n)
	for i := range g {
		g[i] = make([]byte, m)
	}
	for i := 0; i < n; i += 2 {
		for j := 0; j < m; j += 2 {
			if rng.Intn(2) == 0 {
				g[i][j], g[i+1][j] = 'U', 'D'
				g[i][j+1], g[i+1][j+1] = 'U', 'D'
			} else {
				g[i][j], g[i][j+1] = 'L', 'R'
				g[i+1][j], g[i+1][j+1] = 'L', 'R'
			}
		}
	}
	res := make([]string, n)
	for i := range g {
		res[i] = string(g[i])
	}
	return res
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(3)*2 + 2
	m := rng.Intn(3)*2 + 2
	g1 := genGrid(rng, n, m)
	g2 := genGrid(rng, n, m)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		sb.WriteString(g1[i])
		sb.WriteByte('\n')
	}
	for i := 0; i < n; i++ {
		sb.WriteString(g2[i])
		sb.WriteByte('\n')
	}
	return sb.String()
}

// checkAnswer verifies that the candidate output is a valid transformation
// from grid1 to grid2 via the given sequence of operations.
func checkAnswer(n, m int, grid1 [][]byte, grid2 [][]byte, output string) error {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) == 0 {
		return fmt.Errorf("empty output")
	}
	cnt, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return fmt.Errorf("cannot parse count: %v", err)
	}
	if cnt == -1 {
		// Candidate says impossible. Check if it's actually possible by
		// seeing if the reference can solve it. Since genCase always produces
		// valid domino tilings with same dimensions, a solution always exists.
		return fmt.Errorf("candidate says -1 but solution should exist")
	}
	if len(lines) != cnt+1 {
		return fmt.Errorf("expected %d operation lines, got %d", cnt, len(lines)-1)
	}

	// Apply operations to grid1
	g := make([][]byte, n)
	for i := range g {
		g[i] = make([]byte, m)
		copy(g[i], grid1[i])
	}

	for idx := 1; idx <= cnt; idx++ {
		parts := strings.Fields(strings.TrimSpace(lines[idx]))
		if len(parts) != 2 {
			return fmt.Errorf("op %d: expected 2 values, got %d", idx, len(parts))
		}
		r, err1 := strconv.Atoi(parts[0])
		c, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil {
			return fmt.Errorf("op %d: parse error", idx)
		}
		r-- // 1-indexed to 0-indexed
		c--
		if r < 0 || r+1 >= n || c < 0 || c+1 >= m {
			return fmt.Errorf("op %d: (%d,%d) out of bounds", idx, r+1, c+1)
		}
		// Check that the 2x2 block is a valid state to flip
		a, b := g[r][c], g[r][c+1]
		d, e := g[r+1][c], g[r+1][c+1]
		if a == 'L' && b == 'R' && d == 'L' && e == 'R' {
			// horizontal pair -> vertical pair
			g[r][c], g[r][c+1] = 'U', 'U'
			g[r+1][c], g[r+1][c+1] = 'D', 'D'
		} else if a == 'U' && b == 'U' && d == 'D' && e == 'D' {
			// vertical pair -> horizontal pair
			g[r][c], g[r][c+1] = 'L', 'R'
			g[r+1][c], g[r+1][c+1] = 'L', 'R'
		} else {
			return fmt.Errorf("op %d: invalid 2x2 block at (%d,%d): %c%c/%c%c", idx, r+1, c+1, a, b, d, e)
		}
	}

	// Check final grid matches grid2
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if g[i][j] != grid2[i][j] {
				return fmt.Errorf("after ops, grid[%d][%d]=%c but target=%c", i, j, g[i][j], grid2[i][j])
			}
		}
	}
	return nil
}

func parseInput(input string) (int, int, [][]byte, [][]byte) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	parts := strings.Fields(lines[0])
	n, _ := strconv.Atoi(parts[0])
	m, _ := strconv.Atoi(parts[1])
	g1 := make([][]byte, n)
	for i := 0; i < n; i++ {
		g1[i] = []byte(strings.TrimSpace(lines[1+i]))
	}
	g2 := make([][]byte, n)
	for i := 0; i < n; i++ {
		g2[i] = []byte(strings.TrimSpace(lines[1+n+i]))
	}
	return n, m, g1, g2
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierD /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		input := genCase(rng)
		candOut, cErr := runBinary(candidate, input)
		if cErr != nil {
			fmt.Fprintf(os.Stderr, "case %d: candidate error: %v\ninput:\n%s", t+1, cErr, input)
			os.Exit(1)
		}
		n, m, g1, g2 := parseInput(input)
		if err := checkAnswer(n, m, g1, g2, candOut); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\noutput:\n%s\n", t+1, err, input, candOut)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
