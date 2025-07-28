package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
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
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func solve(n, k int, grid [][]int) string {
	mismatches := 0
	for i := 0; i < n; i++ {
		ni := n - 1 - i
		for j := 0; j < n; j++ {
			nj := n - 1 - j
			if i < ni || (i == ni && j < nj) {
				if grid[i][j] != grid[ni][nj] {
					mismatches++
				}
			}
		}
	}
	if k < mismatches {
		return "NO"
	}
	if n%2 == 1 || (k-mismatches)%2 == 0 {
		return "YES"
	}
	return "NO"
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: verifierB.go path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	if b, err := filepath.Abs(bin); err == nil {
		bin = b
	}

	rand.Seed(2)
	const T = 100
	type test struct {
		n    int
		k    int
		grid [][]int
	}
	tests := make([]test, T)
	var input strings.Builder
	input.WriteString(strconv.Itoa(T) + "\n")
	for i := 0; i < T; i++ {
		n := rand.Intn(4) + 1 // 1..4
		k := rand.Intn(11)    // 0..10
		g := make([][]int, n)
		for r := 0; r < n; r++ {
			g[r] = make([]int, n)
			for c := 0; c < n; c++ {
				g[r][c] = rand.Intn(2)
			}
		}
		tests[i] = test{n, k, g}
		input.WriteString(fmt.Sprintf("%d %d\n", n, k))
		for r := 0; r < n; r++ {
			for c := 0; c < n; c++ {
				if c > 0 {
					input.WriteByte(' ')
				}
				input.WriteString(strconv.Itoa(g[r][c]))
			}
			input.WriteByte('\n')
		}
	}

	out, err := runBinary(bin, input.String())
	if err != nil {
		fmt.Println("binary error:", err)
		os.Exit(1)
	}
	fields := strings.Fields(out)
	if len(fields) != T {
		fmt.Printf("wrong number of outputs: got %d want %d\n", len(fields), T)
		os.Exit(1)
	}
	for i := 0; i < T; i++ {
		expect := solve(tests[i].n, tests[i].k, tests[i].grid)
		if fields[i] != expect {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, expect, fields[i])
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
