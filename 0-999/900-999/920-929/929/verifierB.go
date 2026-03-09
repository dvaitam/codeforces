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

type plane [][]byte

func solveB(n, k int, rows []string) (int, []string) {
	s := make([][]byte, n)
	for i := 0; i < n; i++ {
		s[i] = []byte(rows[i])
	}
	kLeft := k
	for i := 0; i < n; i++ {
		if kLeft == 0 {
			break
		}
		if s[i][0] == '.' && s[i][1] != 'S' {
			s[i][0] = 'x'
			kLeft--
		}
		if kLeft == 0 {
			break
		}
		if s[i][11] == '.' && s[i][10] != 'S' {
			s[i][11] = 'x'
			kLeft--
		}
		if kLeft == 0 {
			break
		}
		for j := 1; j < 11 && kLeft > 0; j++ {
			if s[i][j] == '.' && s[i][j-1] != 'S' && s[i][j+1] != 'S' {
				s[i][j] = 'x'
				kLeft--
			}
		}
	}

	for i := 0; i < n; i++ {
		if kLeft == 0 {
			break
		}
		if s[i][0] == '.' {
			s[i][0] = 'x'
			kLeft--
		}
		if kLeft == 0 {
			break
		}
		if s[i][11] == '.' {
			s[i][11] = 'x'
			kLeft--
		}
		if kLeft == 0 {
			break
		}
		for j := 1; j < 11 && kLeft > 0; j++ {
			if s[i][j] == '.' && !(s[i][j-1] == 'S' && s[i][j+1] == 'S') {
				s[i][j] = 'x'
				kLeft--
			}
		}
	}

	for i := 0; i < n && kLeft > 0; i++ {
		for j := 0; j < 12 && kLeft > 0; j++ {
			if s[i][j] == '.' {
				s[i][j] = 'x'
				kLeft--
			}
		}
	}

	ans := 0
	for i := 0; i < n; i++ {
		if s[i][0] == 'S' {
			if s[i][1] != '.' && s[i][1] != '-' {
				ans++
			}
		}
		if s[i][11] == 'S' {
			if s[i][10] != '.' && s[i][10] != '-' {
				ans++
			}
		}
		for j := 1; j < 11; j++ {
			if s[i][j] == 'S' {
				if s[i][j-1] != '.' && s[i][j-1] != '-' {
					ans++
				}
				if s[i][j+1] != '.' && s[i][j+1] != '-' {
					ans++
				}
			}
		}
	}

	out := make([]string, n)
	for i := 0; i < n; i++ {
		out[i] = string(s[i])
	}
	return ans, out
}

// countNeighbors computes the total number of neighbors of 'S' passengers.
// Two seats are neighbors if they are adjacent in the same row with no '-' between them.
func countNeighbors(grid [][]byte) int {
	ans := 0
	for _, row := range grid {
		for j, ch := range row {
			if ch != 'S' {
				continue
			}
			if j > 0 && row[j-1] != '-' && row[j-1] != '.' {
				ans++
			}
			if j < len(row)-1 && row[j+1] != '-' && row[j+1] != '.' {
				ans++
			}
		}
	}
	return ans
}

type caseInfo struct {
	input      string
	n, k       int
	rows       []string
	optimalAns int
}

func generateCase(rng *rand.Rand) caseInfo {
	n := rng.Intn(20) + 1
	rows := make([]string, n)
	freeSeats := 0
	for i := 0; i < n; i++ {
		b := make([]byte, 12)
		for j := 0; j < 12; j++ {
			if j == 3 || j == 8 {
				b[j] = '-'
			} else {
				if rng.Intn(5) == 0 {
					b[j] = 'S'
				} else {
					b[j] = '.'
					freeSeats++
				}
			}
		}
		rows[i] = string(b)
	}
	// k must not exceed the number of free seats (problem guarantee).
	maxK := freeSeats
	if maxK == 0 {
		maxK = 1
	}
	k := rng.Intn(maxK + 1)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i := 0; i < n; i++ {
		sb.WriteString(rows[i])
		sb.WriteByte('\n')
	}
	optimalAns, _ := solveB(n, k, rows)
	return caseInfo{input: sb.String(), n: n, k: k, rows: rows, optimalAns: optimalAns}
}

func runCase(exe string, ci caseInfo) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(ci.input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\nstderr:\n%s", err, stderr.String())
	}
	return checkOutput(ci, strings.TrimSpace(stdout.String()))
}

func checkOutput(ci caseInfo, output string) error {
	lines := strings.Split(output, "\n")
	if len(lines) < 1+ci.n {
		return fmt.Errorf("expected %d output lines, got %d", 1+ci.n, len(lines))
	}

	statedCount, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return fmt.Errorf("first line is not an integer: %q", lines[0])
	}

	gotGrid := make([][]byte, ci.n)
	for i := 0; i < ci.n; i++ {
		gotGrid[i] = []byte(strings.TrimSpace(lines[1+i]))
		if len(gotGrid[i]) != len(ci.rows[i]) {
			return fmt.Errorf("row %d: expected length %d, got %d", i+1, len(ci.rows[i]), len(gotGrid[i]))
		}
	}

	// Verify exactly k free seats were filled with 'x', nothing else changed.
	filled := 0
	for i := 0; i < ci.n; i++ {
		orig := ci.rows[i]
		got := gotGrid[i]
		for j := range orig {
			if orig[j] == '.' && got[j] == 'x' {
				filled++
			} else if orig[j] != got[j] {
				return fmt.Errorf("row %d col %d: original %q illegally changed to %q", i+1, j+1, orig[j], got[j])
			}
		}
	}
	if filled != ci.k {
		return fmt.Errorf("expected exactly %d seats filled, got %d", ci.k, filled)
	}

	// Verify stated count matches actual grid.
	actualCount := countNeighbors(gotGrid)
	if actualCount != statedCount {
		return fmt.Errorf("stated neighbor count %d does not match computed %d", statedCount, actualCount)
	}

	// Verify optimality.
	if actualCount != ci.optimalAns {
		return fmt.Errorf("neighbor count %d is not optimal (expected %d)", actualCount, ci.optimalAns)
	}

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		ci := generateCase(rng)
		if err := runCase(exe, ci); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, ci.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
