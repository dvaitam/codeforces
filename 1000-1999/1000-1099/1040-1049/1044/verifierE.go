package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func genGrid(r *rand.Rand) (int, int, [][]int) {
	n := r.Intn(3) + 3 // 3..5
	m := r.Intn(3) + 3
	vals := r.Perm(n * m)
	grid := make([][]int, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]int, m)
		for j := 0; j < m; j++ {
			grid[i][j] = vals[i*m+j] + 1
		}
	}
	return n, m, grid
}

func gridInput(n, m int, g [][]int) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", g[i][j])
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func verifyOutput(n, m int, grid [][]int, out string) error {
	scanner := bufio.NewScanner(strings.NewReader(out))
	if !scanner.Scan() {
		return fmt.Errorf("no output")
	}
	k, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil || k < 0 {
		return fmt.Errorf("invalid k")
	}
	pos := make(map[int][2]int)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			pos[grid[i][j]] = [2]int{i, j}
		}
	}
	total := 0
	for step := 0; step < k; step++ {
		if !scanner.Scan() {
			return fmt.Errorf("missing move %d", step+1)
		}
		fields := strings.Fields(scanner.Text())
		if len(fields) < 1 {
			return fmt.Errorf("bad line for move %d", step+1)
		}
		s, err := strconv.Atoi(fields[0])
		if err != nil || s < 4 || len(fields) != s+1 {
			return fmt.Errorf("bad cycle length at move %d", step+1)
		}
		total += s
		if total > 100000 {
			return fmt.Errorf("total cycle length exceeds limit")
		}
		seq := make([]int, s)
		seen := make(map[int]bool)
		for i := 0; i < s; i++ {
			v, err := strconv.Atoi(fields[i+1])
			if err != nil {
				return fmt.Errorf("bad number in move %d", step+1)
			}
			if v < 1 || v > n*m || seen[v] {
				return fmt.Errorf("invalid or duplicate number in move %d", step+1)
			}
			seen[v] = true
			seq[i] = v
		}
		// check adjacency
		for i := 0; i < s; i++ {
			a := seq[i]
			b := seq[(i+1)%s]
			pa := pos[a]
			pb := pos[b]
			if abs(pa[0]-pb[0])+abs(pa[1]-pb[1]) != 1 {
				return fmt.Errorf("non-adjacent cells in move %d", step+1)
			}
		}
		// apply rotation
		newPos := make(map[int][2]int)
		for v, p := range pos {
			newPos[v] = p
		}
		for i := 0; i < s; i++ {
			from := seq[i]
			to := pos[seq[(i+1)%s]]
			newPos[from] = to
		}
		// update grid and positions
		for v, p := range newPos {
			grid[p[0]][p[1]] = v
		}
		pos = newPos
	}
	// grid should be sorted
	val := 1
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] != val {
				return fmt.Errorf("grid not sorted")
			}
			val++
		}
	}
	if scanner.Scan() {
		return fmt.Errorf("extra output")
	}
	return nil
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	r := rand.New(rand.NewSource(1))
	for i := 0; i < 100; i++ {
		n, m, g := genGrid(r)
		input := gridInput(n, m, g)
		output, err := run(exe, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if err := verifyOutput(n, m, g, output); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%soutput:\n%s", i+1, err, input, output)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", 100)
}
