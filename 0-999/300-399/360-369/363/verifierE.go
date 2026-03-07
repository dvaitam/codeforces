package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

// circleSum returns the sum of grid cells within circle of radius r centered at (ci,cj) (0-indexed).
func circleSum(grid [][]int, r, ci, cj int) int {
	n, m := len(grid), len(grid[0])
	sum := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			di, dj := i-ci, j-cj
			if di*di+dj*dj <= r*r {
				sum += grid[i][j]
			}
		}
	}
	return sum
}

// intersects returns true if the two circles share at least one grid cell.
func intersects(r, i1, j1, i2, j2, n, m int) bool {
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			di1, dj1 := i-i1, j-j1
			di2, dj2 := i-i2, j-j2
			if di1*di1+dj1*dj1 <= r*r && di2*di2+dj2*dj2 <= r*r {
				return true
			}
		}
	}
	return false
}

// oracle computes the answer by brute force.
// Centers are 0-indexed; valid centers satisfy r ≤ i ≤ n-1-r, r ≤ j ≤ m-1-r.
func oracle(n, m, r int, grid [][]int) (int, int64) {
	type center struct{ i, j int }
	var centers []center
	for i := r; i < n-r; i++ {
		for j := r; j < m-r; j++ {
			centers = append(centers, center{i, j})
		}
	}
	if len(centers) == 0 {
		return 0, 0
	}
	sums := make([]int, len(centers))
	for k, c := range centers {
		sums[k] = circleSum(grid, r, c.i, c.j)
	}

	bestSum := 0
	var bestCnt int64
	for a := 0; a < len(centers); a++ {
		for b := a + 1; b < len(centers); b++ {
			if intersects(r, centers[a].i, centers[a].j, centers[b].i, centers[b].j, n, m) {
				continue
			}
			s := sums[a] + sums[b]
			if s > bestSum {
				bestSum = s
				bestCnt = 1
			} else if s == bestSum {
				bestCnt++
			}
		}
	}
	return bestSum, bestCnt
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCase(rng *rand.Rand) (int, int, int, [][]int) {
	n := rng.Intn(4) + 3
	m := rng.Intn(4) + 3
	maxR := (n - 1) / 2
	if (m-1)/2 < maxR {
		maxR = (m - 1) / 2
	}
	r := rng.Intn(maxR + 1)
	grid := make([][]int, n)
	for i := range grid {
		grid[i] = make([]int, m)
		for j := range grid[i] {
			grid[i][j] = rng.Intn(10) + 1
		}
	}
	return n, m, r, grid
}

func buildInput(n, m, r int, grid [][]int) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, m, r)
	for i := range grid {
		for j := range grid[i] {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", grid[i][j])
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 200; i++ {
		n, m, r, grid := generateCase(rng)
		input := buildInput(n, m, r, grid)

		wantSum, wantCnt := oracle(n, m, r, grid)
		want := fmt.Sprintf("%d %d", wantSum, wantCnt)

		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}

		// Parse candidate output — accept "sum count" on one or two lines,
		// and "%I64d"-style wide-padded output from Windows-compiled C++.
		var gotSum int
		var gotCnt int64
		n2, _ := fmt.Sscan(got, &gotSum, &gotCnt)
		if n2 < 2 || gotSum != wantSum || gotCnt != wantCnt {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, want, got)
			fmt.Fprintf(os.Stderr, "input:\n%s", input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
