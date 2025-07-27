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

func generateTests() [][][]byte {
	r := rand.New(rand.NewSource(45))
	tests := make([][][]byte, 100)
	for i := 0; i < 100; i++ {
		n := r.Intn(6) + 1
		m := r.Intn(6) + 1
		grid := make([][]byte, n)
		for a := 0; a < n; a++ {
			row := make([]byte, m)
			for b := 0; b < m; b++ {
				row[b] = byte('a' + r.Intn(3))
			}
			grid[a] = row
		}
		tests[i] = grid
	}
	return tests
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func expected(grid [][]byte) int64 {
	n := len(grid)
	m := len(grid[0])
	lft := make([][]int, n)
	rgt := make([][]int, n)
	for i := 0; i < n; i++ {
		lft[i] = make([]int, m)
		rgt[i] = make([]int, m)
		for j := 0; j < m; j++ {
			if j > 0 && grid[i][j] == grid[i][j-1] {
				lft[i][j] = lft[i][j-1] + 1
			} else {
				lft[i][j] = 1
			}
		}
		for j := m - 1; j >= 0; j-- {
			if j+1 < m && grid[i][j] == grid[i][j+1] {
				rgt[i][j] = rgt[i][j+1] + 1
			} else {
				rgt[i][j] = 1
			}
		}
	}
	d1 := make([][]int, n)
	for i := 0; i < n; i++ {
		d1[i] = make([]int, m)
		for j := 0; j < m; j++ {
			d1[i][j] = 1
			if i > 0 && grid[i][j] == grid[i-1][j] {
				ext := min(d1[i-1][j], min(lft[i][j], rgt[i][j])-1)
				if ext >= 0 {
					d1[i][j] = ext + 1
				}
			}
		}
	}
	d2 := make([][]int, n)
	for i := n - 1; i >= 0; i-- {
		d2[i] = make([]int, m)
		for j := 0; j < m; j++ {
			d2[i][j] = 1
			if i+1 < n && grid[i][j] == grid[i+1][j] {
				ext := min(d2[i+1][j], min(lft[i][j], rgt[i][j])-1)
				if ext >= 0 {
					d2[i][j] = ext + 1
				}
			}
		}
	}
	var ans int64
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			cnt := d1[i][j]
			if d2[i][j] < cnt {
				cnt = d2[i][j]
			}
			ans += int64(cnt)
		}
	}
	return ans
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	tests := generateTests()
	for idx, grid := range tests {
		var sb strings.Builder
		n := len(grid)
		m := len(grid[0])
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i := 0; i < n; i++ {
			sb.WriteString(string(grid[i]))
			sb.WriteByte('\n')
		}
		exp := expected(grid)
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		var ans int64
		if _, err := fmt.Sscan(got, &ans); err != nil {
			fmt.Printf("Test %d: cannot parse output %q\n", idx+1, got)
			os.Exit(1)
		}
		if ans != exp {
			fmt.Printf("Test %d failed. Input:\n%sExpected %d got %d\n", idx+1, sb.String(), exp, ans)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
