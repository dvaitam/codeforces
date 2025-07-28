package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func runBinary(bin string, input string) (string, error) {
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

func solve(input string) string {
	in := bufio.NewScanner(strings.NewReader(input))
	grid := make([]string, 0, 3)
	for i := 0; i < 3; i++ {
		in.Scan()
		grid = append(grid, strings.TrimSpace(in.Text()))
	}
	dirs := [][2]int{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1}}
	best := "~~~"
	for x1 := 0; x1 < 3; x1++ {
		for y1 := 0; y1 < 3; y1++ {
			for _, d1 := range dirs {
				x2, y2 := x1+d1[0], y1+d1[1]
				if x2 < 0 || x2 >= 3 || y2 < 0 || y2 >= 3 {
					continue
				}
				for _, d2 := range dirs {
					x3, y3 := x2+d2[0], y2+d2[1]
					if x3 < 0 || x3 >= 3 || y3 < 0 || y3 >= 3 {
						continue
					}
					if (x3 == x1 && y3 == y1) || (x3 == x2 && y3 == y2) {
						continue
					}
					word := string([]byte{grid[x1][y1], grid[x2][y2], grid[x3][y3]})
					if word < best {
						best = word
					}
				}
			}
		}
	}
	return best
}

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]

	file, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not open testcasesA.txt:", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	idx := 0
	for {
		lines := make([]string, 0, 3)
		for i := 0; i < 3; i++ {
			if !scanner.Scan() {
				if i == 0 {
					// end of file
					fmt.Printf("All %d tests passed\n", idx)
					return
				}
				fmt.Fprintln(os.Stderr, "unexpected EOF")
				os.Exit(1)
			}
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				i--
				continue
			}
			if len(line) != 3 {
				fmt.Fprintf(os.Stderr, "invalid test case at %d\n", idx+1)
				os.Exit(1)
			}
			lines = append(lines, line)
		}
		// read optional blank line between testcases
		scanner.Scan()
		input := strings.Join(lines, "\n") + "\n"
		exp := solve(input)
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", idx+1, exp, got, input)
			os.Exit(1)
		}
		idx++
	}
}
