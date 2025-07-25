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

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveB(grid []string) int {
	n := len(grid)
	m := len(grid[0])
	minRow, maxRow := n, -1
	minCol, maxCol := m, -1
	for i := 0; i < n; i++ {
		for j, c := range grid[i] {
			if c == '*' {
				if i < minRow {
					minRow = i
				}
				if i > maxRow {
					maxRow = i
				}
				if j < minCol {
					minCol = j
				}
				if j > maxCol {
					maxCol = j
				}
			}
		}
	}
	height := maxRow - minRow + 1
	width := maxCol - minCol + 1
	if height > width {
		return height
	}
	return width
}

func buildCaseB(n, m int, grid []string) (string, string) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		sb.WriteString(grid[i])
		sb.WriteByte('\n')
	}
	expect := fmt.Sprintf("%d", solveB(grid))
	return sb.String(), expect
}

func randomCaseB(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	m := rng.Intn(20) + 1
	grid := make([]string, n)
	for i := 0; i < n; i++ {
		row := make([]byte, m)
		for j := 0; j < m; j++ {
			if rng.Intn(3) == 0 {
				row[j] = '*'
			} else {
				row[j] = '.'
			}
		}
		grid[i] = string(row)
	}
	// ensure at least one star
	r := rng.Intn(n)
	c := rng.Intn(m)
	rowBytes := []byte(grid[r])
	rowBytes[c] = '*'
	grid[r] = string(rowBytes)
	return buildCaseB(n, m, grid)
}

func runCase(bin string, input, expect string) error {
	out, err := runBinary(bin, input)
	if err != nil {
		return err
	}
	if out != expect {
		return fmt.Errorf("expected %q got %q", expect, out)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var tests [][2]string
	// basic cases
	in, exp := buildCaseB(1, 1, []string{"*"})
	tests = append(tests, [2]string{in, exp})
	in, exp = buildCaseB(2, 2, []string{"*.", ".."})
	tests = append(tests, [2]string{in, exp})
	in, exp = buildCaseB(3, 3, []string{"*..", ".*.", "..*"})
	tests = append(tests, [2]string{in, exp})
	in, exp = buildCaseB(2, 3, []string{"***", "***"})
	tests = append(tests, [2]string{in, exp})
	in, exp = buildCaseB(3, 4, []string{"....", ".*..", "...."})
	tests = append(tests, [2]string{in, exp})

	for len(tests) < 100 {
		in, exp := randomCaseB(rng)
		tests = append(tests, [2]string{in, exp})
	}

	for i, tc := range tests {
		if err := runCase(bin, tc[0], tc[1]); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc[0])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
