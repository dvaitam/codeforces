package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func generateFrame(n, m int, rng *rand.Rand) (int, int, int, int) {
	x1 := rng.Intn(n-2) + 1
	x2 := rng.Intn(n-x1-1) + x1 + 2
	y1 := rng.Intn(m-2) + 1
	y2 := rng.Intn(m-y1-1) + y1 + 2
	return x1, y1, x2, y2
}

func applyFrame(g [][]byte, x1, y1, x2, y2 int) {
	for y := y1; y <= y2; y++ {
		g[x1][y] = '#'
		g[x2][y] = '#'
	}
	for x := x1; x <= x2; x++ {
		g[x][y1] = '#'
		g[x][y2] = '#'
	}
}

func generateYesCase(rng *rand.Rand) (string, [][]byte) {
	n := rng.Intn(4) + 3 // 3..6
	m := rng.Intn(4) + 3
	grid := make([][]byte, n+1)
	for i := range grid {
		grid[i] = make([]byte, m+1)
		for j := range grid[i] {
			grid[i][j] = '.'
		}
	}
	x1, y1, x2, y2 := generateFrame(n, m, rng)
	x3, y3, x4, y4 := generateFrame(n, m, rng)
	applyFrame(grid, x1, y1, x2, y2)
	applyFrame(grid, x3, y3, x4, y4)

	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			sb.WriteByte(grid[i][j])
		}
		sb.WriteByte('\n')
	}
	return sb.String(), grid
}

func generateNoCase(rng *rand.Rand) (string, [][]byte) {
	s, g := generateYesCase(rng)
	nFields := strings.Fields(s)
	n := toInt(nFields[0])
	m := toInt(nFields[1])
	// flip random cell
	x := rng.Intn(n) + 1
	y := rng.Intn(m) + 1
	if g[x][y] == '#' {
		g[x][y] = '.'
	} else {
		g[x][y] = '#'
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			sb.WriteByte(g[i][j])
		}
		sb.WriteByte('\n')
	}
	return sb.String(), g
}

func framesMatch(grid [][]byte, rects [][4]int) bool {
	n := len(grid) - 1
	m := len(grid[0]) - 1
	other := make([][]byte, n+1)
	for i := range other {
		other[i] = make([]byte, m+1)
		for j := range other[i] {
			other[i][j] = '.'
		}
	}
	for _, r := range rects {
		x1, y1, x2, y2 := r[0], r[1], r[2], r[3]
		if x1 < 1 || y1 < 1 || x2 > n || y2 > m || x1+2 > x2 || y1+2 > y2 {
			return false
		}
		for y := y1; y <= y2; y++ {
			other[x1][y] = '#'
			other[x2][y] = '#'
		}
		for x := x1; x <= x2; x++ {
			other[x][y1] = '#'
			other[x][y2] = '#'
		}
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if other[i][j] != grid[i][j] {
				return false
			}
		}
	}
	return true
}

func findFrames(grid [][]byte) (bool, [4]int, [4]int) {
	n := len(grid) - 1
	m := len(grid[0]) - 1
	for x1 := 1; x1 <= n-2; x1++ {
		for x2 := x1 + 2; x2 <= n; x2++ {
			for y1 := 1; y1 <= m-2; y1++ {
				for y2 := y1 + 2; y2 <= m; y2++ {
					for a1 := 1; a1 <= n-2; a1++ {
						for a2 := a1 + 2; a2 <= n; a2++ {
							for b1 := 1; b1 <= m-2; b1++ {
								for b2 := b1 + 2; b2 <= m; b2++ {
									rects := [][4]int{{x1, y1, x2, y2}, {a1, b1, a2, b2}}
									if framesMatch(grid, rects) {
										return true, [4]int{x1, y1, x2, y2}, [4]int{a1, b1, a2, b2}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return false, [4]int{}, [4]int{}
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

func verifyOutput(grid [][]byte, output string, expectYes bool) error {
	outR := bufio.NewReader(strings.NewReader(output))
	var ans string
	if _, err := fmt.Fscan(outR, &ans); err != nil {
		return fmt.Errorf("bad output")
	}
	ans = strings.ToUpper(ans)
	if ans == "NO" {
		if expectYes {
			return fmt.Errorf("expected YES got NO")
		}
		return nil
	}
	if ans != "YES" {
		return fmt.Errorf("invalid first token")
	}
	var r1 [4]int
	var r2 [4]int
	if _, err := fmt.Fscan(outR, &r1[0], &r1[1], &r1[2], &r1[3]); err != nil {
		return fmt.Errorf("failed to read first rectangle")
	}
	if _, err := fmt.Fscan(outR, &r2[0], &r2[1], &r2[2], &r2[3]); err != nil {
		return fmt.Errorf("failed to read second rectangle")
	}
	if !framesMatch(grid, [][4]int{r1, r2}) {
		return fmt.Errorf("rectangles do not match grid")
	}
	if !expectYes {
		return fmt.Errorf("expected NO but got YES")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		var tc string
		var grid [][]byte
		var expectYes bool
		if i%2 == 0 {
			tc, grid = generateYesCase(rng)
			expectYes = true
		} else {
			tc, grid = generateNoCase(rng)
			ok, _, _ := findFrames(grid)
			expectYes = ok
		}
		got, err := run(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if err := verifyOutput(grid, got, expectYes); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func toInt(s string) int {
	var v int
	fmt.Sscan(s, &v)
	return v
}
