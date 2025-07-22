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

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
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

type grid [3][3]int

func expected(g grid) [3][3]int {
	dirs := [][2]int{{0, 0}, {1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	var res [3][3]int
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			sum := 0
			for _, d := range dirs {
				ni, nj := i+d[0], j+d[1]
				if ni >= 0 && ni < 3 && nj >= 0 && nj < 3 {
					sum += g[ni][nj]
				}
			}
			if sum%2 == 0 {
				res[i][j] = 1
			} else {
				res[i][j] = 0
			}
		}
	}
	return res
}

func gridToInput(g grid) string {
	var sb strings.Builder
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			sb.WriteString(fmt.Sprintf("%d", g[i][j]))
			if j+1 < 3 {
				sb.WriteByte(' ')
			}
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func outputToGrid(out string) ([3][3]int, error) {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 {
		return [3][3]int{}, fmt.Errorf("expected 3 lines, got %d", len(lines))
	}
	var res [3][3]int
	for i := 0; i < 3; i++ {
		if len(lines[i]) < 3 {
			return res, fmt.Errorf("line %d too short", i+1)
		}
		for j := 0; j < 3; j++ {
			if lines[i][j] == '0' {
				res[i][j] = 0
			} else if lines[i][j] == '1' {
				res[i][j] = 1
			} else {
				return res, fmt.Errorf("invalid char %q", lines[i][j])
			}
		}
	}
	return res, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 100; i++ {
		var g grid
		for r := 0; r < 3; r++ {
			for c := 0; c < 3; c++ {
				g[r][c] = rng.Intn(101) // 0..100
			}
		}
		input := gridToInput(g)
		expectedGrid := expected(g)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := outputToGrid(out)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", i+1, err, input, out)
			os.Exit(1)
		}
		if got != expectedGrid {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %v got %v\ninput:\n%s", i+1, expectedGrid, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
