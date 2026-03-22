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

func compileBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "verifierA-bin-*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("go build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, nil, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

// embeddedSolve produces a valid placement sequence for the given tile string.
// Strategy: place vertical tiles in columns 0,1 and horizontal tiles in columns 2,3
// using two dedicated column-pairs, ensuring clearing happens regularly.
func embeddedSolve(s string) string {
	var grid [4][4]bool
	var sb strings.Builder

	clearFull := func() {
		var fullRow, fullCol [4]bool
		for i := 0; i < 4; i++ {
			fullRow[i] = true
			for j := 0; j < 4; j++ {
				if !grid[i][j] {
					fullRow[i] = false
					break
				}
			}
		}
		for j := 0; j < 4; j++ {
			fullCol[j] = true
			for i := 0; i < 4; i++ {
				if !grid[i][j] {
					fullCol[j] = false
					break
				}
			}
		}
		for i := 0; i < 4; i++ {
			for j := 0; j < 4; j++ {
				if fullRow[i] || fullCol[j] {
					grid[i][j] = false
				}
			}
		}
	}

	canPlaceVert := func(r, c int) bool {
		if r < 0 || r > 2 || c < 0 || c > 3 {
			return false
		}
		return !grid[r][c] && !grid[r+1][c]
	}

	canPlaceHori := func(r, c int) bool {
		if r < 0 || r > 3 || c < 0 || c > 2 {
			return false
		}
		return !grid[r][c] && !grid[r][c+1]
	}

	// Vertical tiles: prefer columns 0 and 1 (rows 0-1, then 2-3)
	// Horizontal tiles: prefer columns 2 and 3 (rows 0, then 1, 2, 3)
	// This separation ensures they don't interfere with each other.
	// When a column fills up, clearing happens automatically.

	for _, ch := range s {
		if ch == '0' {
			// Vertical 2x1: try col 0 rows 0,2; then col 1 rows 0,2; then all
			placed := false
			for _, pos := range [][2]int{{0, 0}, {2, 0}, {0, 1}, {2, 1}, {0, 2}, {2, 2}, {0, 3}, {2, 3}} {
				if canPlaceVert(pos[0], pos[1]) {
					grid[pos[0]][pos[1]] = true
					grid[pos[0]+1][pos[1]] = true
					clearFull()
					fmt.Fprintf(&sb, "%d %d\n", pos[0]+1, pos[1]+1)
					placed = true
					break
				}
			}
			if !placed {
				// Should never happen if strategy is correct
				fmt.Fprintf(&sb, "1 1\n")
			}
		} else {
			// Horizontal 1x2: try col 2 rows 0,1,2,3; then col 0 rows 0,1,2,3; then all
			placed := false
			for _, pos := range [][2]int{{0, 2}, {1, 2}, {2, 2}, {3, 2}, {0, 0}, {1, 0}, {2, 0}, {3, 0}} {
				if canPlaceHori(pos[0], pos[1]) {
					grid[pos[0]][pos[1]] = true
					grid[pos[0]][pos[1]+1] = true
					clearFull()
					fmt.Fprintf(&sb, "%d %d\n", pos[0]+1, pos[1]+1)
					placed = true
					break
				}
			}
			if !placed {
				fmt.Fprintf(&sb, "1 1\n")
			}
		}
	}
	return strings.TrimSpace(sb.String())
}

func verifyCase(s, out string) error {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != len(s) {
		return fmt.Errorf("expected %d lines, got %d", len(s), len(lines))
	}
	grid := [4][4]int{}
	for idx, ch := range s {
		var r, c int
		if _, err := fmt.Sscanf(lines[idx], "%d %d", &r, &c); err != nil {
			return fmt.Errorf("line %d parse error: %v", idx+1, err)
		}
		if r < 1 || r > 4 || c < 1 || c > 4 {
			return fmt.Errorf("line %d out of bounds", idx+1)
		}
		r--
		c--
		if ch == '0' {
			if r+1 >= 4 {
				return fmt.Errorf("tile %d out of bounds", idx+1)
			}
			if grid[r][c] != 0 || grid[r+1][c] != 0 {
				return fmt.Errorf("tile %d overlaps", idx+1)
			}
			grid[r][c] = 1
			grid[r+1][c] = 1
		} else {
			if c+1 >= 4 {
				return fmt.Errorf("tile %d out of bounds", idx+1)
			}
			if grid[r][c] != 0 || grid[r][c+1] != 0 {
				return fmt.Errorf("tile %d overlaps", idx+1)
			}
			grid[r][c] = 1
			grid[r][c+1] = 1
		}
		// clear full rows and cols simultaneously
		var fullRow, fullCol [4]bool
		for i := 0; i < 4; i++ {
			fullRow[i] = true
			for j := 0; j < 4; j++ {
				if grid[i][j] == 0 {
					fullRow[i] = false
					break
				}
			}
		}
		for j := 0; j < 4; j++ {
			fullCol[j] = true
			for i := 0; i < 4; i++ {
				if grid[i][j] == 0 {
					fullCol[j] = false
					break
				}
			}
		}
		for i := 0; i < 4; i++ {
			for j := 0; j < 4; j++ {
				if fullRow[i] || fullCol[j] {
					grid[i][j] = 0
				}
			}
		}
	}
	return nil
}

func genCase(r *rand.Rand) string {
	n := r.Intn(1000) + 1
	var sb strings.Builder
	for i := 0; i < n; i++ {
		if r.Intn(2) == 0 {
			sb.WriteByte('0')
		} else {
			sb.WriteByte('1')
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin, cleanup, err := compileBinary(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "compile error: %v\n", err)
		os.Exit(1)
	}
	if cleanup != nil {
		defer cleanup()
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []string{"0", "1"}
	cases = append(cases, strings.Repeat("0", 1000))
	cases = append(cases, strings.Repeat("1", 1000))
	for i := len(cases); i < 100; i++ {
		cases = append(cases, genCase(rng))
	}
	for idx, s := range cases {
		// First verify our embedded solver produces valid output
		refOut := embeddedSolve(s)
		if err := verifyCase(s, refOut); err != nil {
			fmt.Fprintf(os.Stderr, "INTERNAL: embedded solver failed case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		// Now test the candidate
		out, err := runCandidate(bin, s+"\n")
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		if err := verifyCase(s, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: %s\noutput:\n%s\n", idx+1, err, s, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
