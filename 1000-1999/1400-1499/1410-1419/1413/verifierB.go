package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// testCase structure
type testCase struct {
	input  string
	expect string
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func solveB(input string) string {
	in := bufio.NewScanner(strings.NewReader(input))
	in.Split(bufio.ScanWords)
	var sb strings.Builder
	nextInt := func() int {
		if !in.Scan() {
			return 0
		}
		v, _ := strconv.Atoi(in.Text())
		return v
	}
	t := nextInt()
	for ; t > 0; t-- {
		n, m := nextInt(), nextInt()
		rows := make([][]int, n)
		elemToRow := make(map[int]int)
		for i := 0; i < n; i++ {
			row := make([]int, m)
			for j := 0; j < m; j++ {
				row[j] = nextInt()
				elemToRow[row[j]] = i
			}
			rows[i] = row
		}
		cols := make([][]int, m)
		for j := 0; j < m; j++ {
			col := make([]int, n)
			for i := 0; i < n; i++ {
				col[i] = nextInt()
			}
			cols[j] = col
		}
		order := make([]int, n)
		for i := 0; i < n; i++ {
			order[i] = elemToRow[cols[0][i]]
		}
		for _, idx := range order {
			row := rows[idx]
			for j, v := range row {
				if j > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(strconv.Itoa(v))
			}
			if t > 1 || idx != order[len(order)-1] {
				sb.WriteByte('\n')
			}
		}
	}
	return strings.TrimSpace(sb.String())
}

func genTests() []testCase {
	rng := rand.New(rand.NewSource(2))
	tests := make([]testCase, 100)
	for i := 0; i < 100; i++ {
		n := rng.Intn(3) + 1
		m := rng.Intn(3) + 1
		matrix := make([][]int, n)
		val := 1
		for r := 0; r < n; r++ {
			row := make([]int, m)
			for c := 0; c < m; c++ {
				row[c] = val
				val++
			}
			matrix[r] = row
		}
		rows := make([][]int, n)
		permRows := rand.Perm(n)
		for i2, pr := range permRows {
			rows[i2] = matrix[pr]
		}
		cols := make([][]int, m)
		// build columns
		baseCols := make([][]int, m)
		for c := 0; c < m; c++ {
			col := make([]int, n)
			for r := 0; r < n; r++ {
				col[r] = matrix[r][c]
			}
			baseCols[c] = col
		}
		permCols := rand.Perm(m)
		for i2, pc := range permCols {
			cols[i2] = baseCols[pc]
		}
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for r := 0; r < n; r++ {
			for c := 0; c < m; c++ {
				if c > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(strconv.Itoa(rows[r][c]))
			}
			sb.WriteByte('\n')
		}
		for c := 0; c < m; c++ {
			for r := 0; r < n; r++ {
				if r > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(strconv.Itoa(cols[c][r]))
			}
			sb.WriteByte('\n')
		}
		input := sb.String()
		expect := solveB(input)
		tests[i] = testCase{input: input, expect: expect}
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		out, err := runBinary(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			return
		}
		if out != tc.expect {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected:\n%s\nactual:\n%s\n", i+1, tc.input, tc.expect, out)
			return
		}
	}
	fmt.Println("All tests passed")
}
