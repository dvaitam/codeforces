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

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		if out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, string(out))
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runProg(exe string, input []byte) (string, error) {
	cmd := exec.Command(exe)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return out.String() + errBuf.String(), err
	}
	return strings.TrimSpace(out.String()), nil
}

func genTest() []byte {
	n := rand.Intn(4) + 2 // 2..5, at least 2 to have off-diagonal entries
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	chars := []byte{'F', 'S', '?'}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i == j {
				sb.WriteByte('.')
			} else if i < j {
				sb.WriteByte(chars[rand.Intn(len(chars))])
			} else {
				// Will be filled to match (j,i) - but for now we need to be symmetric in input
				// Actually the problem input is a symmetric matrix where grid[i][j] == grid[j][i]
				// We already wrote grid[j][i] for j < i, so we need to look it up
				// Simpler: build the grid first, then write it
				sb.WriteByte('X') // placeholder
			}
		}
		sb.WriteByte('\n')
	}
	// Rebuild properly with symmetry
	sb.Reset()
	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]byte, n)
		for j := 0; j < n; j++ {
			if i == j {
				grid[i][j] = '.'
			}
		}
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			c := chars[rand.Intn(len(chars))]
			grid[i][j] = c
			grid[j][i] = c
		}
	}
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		sb.Write(grid[i])
		sb.WriteByte('\n')
	}
	return []byte(sb.String())
}

func validate(input, output []byte) error {
	inLines := strings.Split(strings.TrimSpace(string(input)), "\n")
	if len(inLines) < 1 {
		return fmt.Errorf("empty input")
	}
	var n int
	fmt.Sscan(inLines[0], &n)
	if len(inLines) < n+1 {
		return fmt.Errorf("input has %d lines, expected %d", len(inLines), n+1)
	}
	grid := make([]string, n)
	for i := 0; i < n; i++ {
		grid[i] = strings.TrimSpace(inLines[i+1])
	}

	outLines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(outLines) < n {
		return fmt.Errorf("output has %d lines, expected %d", len(outLines), n)
	}
	out := make([]string, n)
	for i := 0; i < n; i++ {
		out[i] = strings.TrimSpace(outLines[i])
		if len(out[i]) != n {
			return fmt.Errorf("output row %d has length %d, expected %d", i, len(out[i]), n)
		}
	}

	// Check consistency with input and symmetry
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i == j {
				// diagonal should be '.'
				if out[i][j] != '.' {
					return fmt.Errorf("position (%d,%d): expected '.' on diagonal, got %c", i, j, out[i][j])
				}
				continue
			}
			if grid[i][j] == 'F' || grid[i][j] == 'S' {
				if out[i][j] != grid[i][j] {
					return fmt.Errorf("position (%d,%d): expected %c from input, got %c", i, j, grid[i][j], out[i][j])
				}
			} else if grid[i][j] == '?' {
				if out[i][j] != 'F' && out[i][j] != 'S' {
					return fmt.Errorf("position (%d,%d): expected F or S, got %c", i, j, out[i][j])
				}
			}
			if out[i][j] != out[j][i] {
				return fmt.Errorf("not symmetric at (%d,%d): %c vs %c", i, j, out[i][j], out[j][i])
			}
		}
	}

	// Check the constraint: for each guest g, the sequence of show types
	// (considering other guests in order) should not have a consecutive run
	// of F's or S's longer than ceil(3n/4).
	lim := (3*n + 3) / 4

	for g := 0; g < n; g++ {
		runLen := 0
		var lastChar byte
		for j := 0; j < n; j++ {
			if j == g {
				continue
			}
			ch := out[g][j]
			if ch == lastChar {
				runLen++
			} else {
				runLen = 1
				lastChar = ch
			}
			if runLen > lim {
				return fmt.Errorf("guest %d has %d consecutive %c shows, limit is %d", g, runLen, ch, lim)
			}
		}
	}

	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierD /path/to/binary")
		return
	}
	exe, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer cleanup()

	for i := 1; i <= 100; i++ {
		in := genTest()
		got, err := runProg(exe, in)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n%s", i, err, got)
			os.Exit(1)
		}
		if err := validate(in, []byte(got)); err != nil {
			fmt.Printf("wrong answer on test %d: %v\ninput:\n%soutput:\n%s\n", i, err, string(in), got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
