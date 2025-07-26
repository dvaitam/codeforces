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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		return out.String() + errb.String(), err
	}
	return out.String(), nil
}

type Test struct {
	n, m  int
	grid  []string
	input string
}

func genTest(rng *rand.Rand) Test {
	n := rng.Intn(5) + 1
	m := rng.Intn(5) + 1
	grid := make([]string, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		row := make([]byte, m)
		for j := 0; j < m; j++ {
			if rng.Intn(2) == 0 {
				row[j] = '.'
			} else {
				row[j] = '#'
			}
		}
		grid[i] = string(row)
		sb.WriteString(grid[i])
		if i+1 < n {
			sb.WriteByte('\n')
		}
	}
	return Test{n: n, m: m, grid: grid, input: sb.String()}
}

func solve(t Test) string {
	rowMask := make([]uint64, t.n)
	colMask := make([]uint64, t.m)
	for i := 0; i < t.n; i++ {
		for j := 0; j < t.m; j++ {
			if t.grid[i][j] == '#' {
				rowMask[i] |= 1 << uint(j)
				colMask[j] |= 1 << uint(i)
			}
		}
	}
	for i := 0; i < t.n; i++ {
		for j := i + 1; j < t.n; j++ {
			if rowMask[i]&rowMask[j] != 0 && rowMask[i] != rowMask[j] {
				return "No"
			}
		}
	}
	for i := 0; i < t.m; i++ {
		for j := i + 1; j < t.m; j++ {
			if colMask[i]&colMask[j] != 0 && colMask[i] != colMask[j] {
				return "No"
			}
		}
	}
	return "Yes"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		t := genTest(rng)
		expected := solve(t)
		out, err := run(bin, t.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n%s", i+1, err, out)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		if out != expected {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected:%s got:%s\n", i+1, t.input, expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
