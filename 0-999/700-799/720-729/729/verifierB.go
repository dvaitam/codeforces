package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(n, m int, grid [][]int) string {
	prefRow := make([][]int, n)
	for i := 0; i < n; i++ {
		prefRow[i] = make([]int, m)
		sum := 0
		for j := 0; j < m; j++ {
			if grid[i][j] == 1 {
				sum++
			}
			prefRow[i][j] = sum
		}
	}
	prefCol := make([][]int, n)
	for i := 0; i < n; i++ {
		prefCol[i] = make([]int, m)
	}
	for j := 0; j < m; j++ {
		sum := 0
		for i := 0; i < n; i++ {
			if grid[i][j] == 1 {
				sum++
			}
			prefCol[i][j] = sum
		}
	}
	ans := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == 1 {
				continue
			}
			if j > 0 && prefRow[i][j-1] > 0 {
				ans++
			}
			if j+1 < m && prefRow[i][m-1]-prefRow[i][j] > 0 {
				ans++
			}
			if i > 0 && prefCol[i-1][j] > 0 {
				ans++
			}
			if i+1 < n && prefCol[n-1][j]-prefCol[i][j] > 0 {
				ans++
			}
		}
	}
	return fmt.Sprintf("%d", ans)
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	r := rand.New(rand.NewSource(2))
	for tc := 1; tc <= 100; tc++ {
		n := r.Intn(8) + 1
		m := r.Intn(8) + 1
		grid := make([][]int, n)
		ones := 0
		for i := 0; i < n; i++ {
			grid[i] = make([]int, m)
			for j := 0; j < m; j++ {
				v := r.Intn(2)
				grid[i][j] = v
				if v == 1 {
					ones++
				}
			}
		}
		if ones == 0 {
			grid[0][0] = 1
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, m)
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				fmt.Fprintf(&sb, "%d", grid[i][j])
				if j+1 < m {
					sb.WriteByte(' ')
				}
			}
			sb.WriteByte('\n')
		}
		input := sb.String()
		expect := expected(n, m, grid)
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\n", tc, err)
			os.Exit(1)
		}
		if out != expect {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", tc, input, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}
