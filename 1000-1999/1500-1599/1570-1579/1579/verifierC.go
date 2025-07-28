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

func solveC(n, m, k int, grid []string) string {
	covered := make([][]bool, n)
	for i := range covered {
		covered[i] = make([]bool, m)
	}
	for i := n - 1; i >= 0; i-- {
		for j := 0; j < m; j++ {
			if grid[i][j] != '*' {
				continue
			}
			size := 0
			for {
				x := i - size - 1
				y1 := j - size - 1
				y2 := j + size + 1
				if x < 0 || y1 < 0 || y2 >= m {
					break
				}
				if grid[x][y1] == '*' && grid[x][y2] == '*' {
					size++
				} else {
					break
				}
			}
			if size >= k {
				covered[i][j] = true
				for d := 1; d <= size; d++ {
					covered[i-d][j-d] = true
					covered[i-d][j+d] = true
				}
			}
		}
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == '*' && !covered[i][j] {
				return "NO\n"
			}
		}
	}
	return "YES\n"
}

func genCaseC(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	m := rng.Intn(19) + 1
	k := rng.Intn(n) + 1
	grid := make([]string, n)
	for i := 0; i < n; i++ {
		b := make([]byte, m)
		for j := 0; j < m; j++ {
			if rng.Intn(2) == 0 {
				b[j] = '.'
			} else {
				b[j] = '*'
			}
		}
		grid[i] = string(b)
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
	for _, row := range grid {
		sb.WriteString(row)
		sb.WriteByte('\n')
	}
	input := sb.String()
	expect := solveC(n, m, k, grid)
	return input, expect
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expect := genCaseC(rng)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\ninput:\n%soutput:\n%s", i+1, err, input, got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%sbut got:\n%s", i+1, input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
