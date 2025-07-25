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

func calcMaxRow(row []int) int {
	maxCnt := 0
	cnt := 0
	for _, v := range row {
		if v == 1 {
			cnt++
			if cnt > maxCnt {
				maxCnt = cnt
			}
		} else {
			cnt = 0
		}
	}
	return maxCnt
}

func solveCase(n, m, q int, grid [][]int, ops [][2]int) string {
	rowMax := make([]int, n)
	for i := 0; i < n; i++ {
		rowMax[i] = calcMaxRow(grid[i])
	}
	var sb strings.Builder
	for _, op := range ops {
		r := op[0] - 1
		c := op[1] - 1
		if grid[r][c] == 1 {
			grid[r][c] = 0
		} else {
			grid[r][c] = 1
		}
		rowMax[r] = calcMaxRow(grid[r])
		ans := 0
		for _, v := range rowMax {
			if v > ans {
				ans = v
			}
		}
		sb.WriteString(fmt.Sprintf("%d\n", ans))
	}
	return strings.TrimSpace(sb.String())
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	m := rng.Intn(5) + 1
	q := rng.Intn(20) + 1
	grid := make([][]int, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]int, m)
		for j := 0; j < m; j++ {
			grid[i][j] = rng.Intn(2)
		}
	}
	ops := make([][2]int, q)
	for i := 0; i < q; i++ {
		ops[i] = [2]int{rng.Intn(n) + 1, rng.Intn(m) + 1}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, q))
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", grid[i][j]))
		}
		sb.WriteByte('\n')
	}
	for i := 0; i < q; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", ops[i][0], ops[i][1]))
	}
	input := sb.String()
	expect := solveCase(n, m, q, copyGrid(grid), ops)
	return input, expect
}

func copyGrid(src [][]int) [][]int {
	n := len(src)
	dest := make([][]int, n)
	for i := 0; i < n; i++ {
		dest[i] = make([]int, len(src[i]))
		copy(dest[i], src[i])
	}
	return dest
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expect := genCase(rng)
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
