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

func runProg(exe, input string) (string, error) {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return out.String() + errBuf.String(), fmt.Errorf("%v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(6) + 1
	a := rng.Perm(n)
	b := rng.Perm(n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range a {
		sb.WriteString(fmt.Sprintf("%d", v+1))
		if i+1 < n {
			sb.WriteByte(' ')
		}
	}
	sb.WriteByte('\n')
	for i, v := range b {
		sb.WriteString(fmt.Sprintf("%d", v+1))
		if i+1 < n {
			sb.WriteByte(' ')
		}
	}
	sb.WriteByte('\n')
	return sb.String()
}

// validate checks if the candidate output is a valid solution for the given input.
// Returns "" if valid, or an error description.
func validate(input, output string) string {
	// Parse input
	fields := strings.Fields(input)
	if len(fields) < 1 {
		return "empty input"
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil || n < 1 {
		return "bad n"
	}
	if len(fields) < 1+2*n {
		return "input too short"
	}
	r := make([]int, n+1) // r[i]: entering row i from left exits at row r[i] on right
	c := make([]int, n+1) // c[i]: entering col i from top exits at col c[i] on bottom
	for i := 1; i <= n; i++ {
		r[i], _ = strconv.Atoi(fields[i])
	}
	for i := 1; i <= n; i++ {
		c[i], _ = strconv.Atoi(fields[n+i])
	}

	// Parse output
	outFields := strings.Fields(output)
	if len(outFields) < 1 {
		return "empty output"
	}
	// Check for -1
	if outFields[0] == "-1" {
		return "candidate says -1 but a solution should exist"
	}
	m, err := strconv.Atoi(outFields[0])
	if err != nil {
		return fmt.Sprintf("bad m: %v", err)
	}
	if m < 0 || m > n*n/2 {
		return fmt.Sprintf("m=%d out of range", m)
	}
	if len(outFields) < 1+4*m {
		return fmt.Sprintf("output too short: need %d fields for %d portals, got %d", 1+4*m, m, len(outFields))
	}

	// Build portal map: (row,col) -> (row,col)
	type cell struct{ r, c int }
	portal := make(map[cell]cell)
	used := make(map[cell]bool)
	idx := 1
	for i := 0; i < m; i++ {
		x1, _ := strconv.Atoi(outFields[idx])
		y1, _ := strconv.Atoi(outFields[idx+1])
		x2, _ := strconv.Atoi(outFields[idx+2])
		y2, _ := strconv.Atoi(outFields[idx+3])
		idx += 4

		c1 := cell{x1, y1}
		c2 := cell{x2, y2}
		if x1 < 1 || x1 > n || y1 < 1 || y1 > n || x2 < 1 || x2 > n || y2 < 1 || y2 > n {
			return fmt.Sprintf("portal %d out of bounds: (%d,%d)-(%d,%d)", i+1, x1, y1, x2, y2)
		}
		if c1 == c2 {
			return fmt.Sprintf("portal %d: same cell (%d,%d)", i+1, x1, y1)
		}
		if used[c1] {
			return fmt.Sprintf("portal %d: cell (%d,%d) already used", i+1, x1, y1)
		}
		if used[c2] {
			return fmt.Sprintf("portal %d: cell (%d,%d) already used", i+1, x2, y2)
		}
		used[c1] = true
		used[c2] = true
		portal[c1] = c2
		portal[c2] = c1
	}

	// Simulate walking right from (i, 1) for each row i
	for i := 1; i <= n; i++ {
		row, col := i, 1
		// direction: right (dr=0, dc=1)
		dr, dc := 0, 1
		steps := 0
		maxSteps := n * n * 4
		for {
			steps++
			if steps > maxSteps {
				return fmt.Sprintf("row %d: infinite loop", i)
			}
			// Check if current cell has a portal door
			cur := cell{row, col}
			if dest, ok := portal[cur]; ok {
				row, col = dest.r, dest.c
			}
			// Move to next cell
			row += dr
			col += dc
			// Check if exited grid
			if row < 1 || row > n || col < 1 || col > n {
				break
			}
		}
		// For walking right, we should exit from column n+1 (col > n)
		// The exit row should be r[i]
		if col != n+1 {
			return fmt.Sprintf("row %d: exited from wrong side (%d,%d)", i, row, col)
		}
		if row != r[i] {
			return fmt.Sprintf("row %d: expected exit row %d, got %d", i, r[i], row)
		}
	}

	// Simulate walking down from (1, i) for each column i
	for i := 1; i <= n; i++ {
		row, col := 1, i
		dr, dc := 1, 0
		steps := 0
		maxSteps := n * n * 4
		for {
			steps++
			if steps > maxSteps {
				return fmt.Sprintf("col %d: infinite loop", i)
			}
			cur := cell{row, col}
			if dest, ok := portal[cur]; ok {
				row, col = dest.r, dest.c
			}
			row += dr
			col += dc
			if row < 1 || row > n || col < 1 || col > n {
				break
			}
		}
		if row != n+1 {
			return fmt.Sprintf("col %d: exited from wrong side (%d,%d)", i, row, col)
		}
		if col != c[i] {
			return fmt.Sprintf("col %d: expected exit col %d, got %d", i, c[i], col)
		}
	}

	return ""
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := generateCase(rng)
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s", i+1, err, got)
			os.Exit(1)
		}
		errMsg := validate(input, got)
		if errMsg != "" {
			fmt.Fprintf(os.Stderr, "case %d failed: %s\ninput:\n%s\noutput:\n%s\n", i+1, errMsg, input, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
