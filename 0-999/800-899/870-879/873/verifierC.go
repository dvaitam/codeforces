package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type TestC struct {
	n, m, k int
	mat     [][]int
}

func (t TestC) Input() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", t.n, t.m, t.k))
	for i := 0; i < t.n; i++ {
		for j := 0; j < t.m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(t.mat[i][j]))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func expectedC(t TestC) string {
	totalScore := 0
	totalRemove := 0
	n, m, k := t.n, t.m, t.k
	for col := 0; col < m; col++ {
		prefix := make([]int, n+1)
		for i := 1; i <= n; i++ {
			prefix[i] = prefix[i-1] + t.mat[i-1][col]
		}
		bestScore := 0
		bestRemove := 0
		for row := 1; row <= n; row++ {
			if t.mat[row-1][col] == 1 {
				onesAbove := prefix[row-1]
				end := row + k - 1
				if end > n {
					end = n
				}
				score := prefix[end] - prefix[row-1]
				if score > bestScore || (score == bestScore && onesAbove < bestRemove) {
					bestScore = score
					bestRemove = onesAbove
				}
			}
		}
		totalScore += bestScore
		totalRemove += bestRemove
	}
	return fmt.Sprintf("%d %d", totalScore, totalRemove)
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return errBuf.String(), fmt.Errorf("runtime error: %v", err)
	}
	return out.String(), nil
}

func genTests() []TestC {
	rand.Seed(3)
	tests := make([]TestC, 0, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(8) + 1
		m := rand.Intn(8) + 1
		k := rand.Intn(n) + 1
		mat := make([][]int, n)
		for r := 0; r < n; r++ {
			row := make([]int, m)
			for c := 0; c < m; c++ {
				row[c] = rand.Intn(2)
			}
			mat[r] = row
		}
		tests = append(tests, TestC{n: n, m: m, k: k, mat: mat})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		exp := strings.TrimSpace(expectedC(tc))
		gotRaw, err := run(bin, tc.Input())
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n%s", i+1, err, gotRaw)
			os.Exit(1)
		}
		got := strings.TrimSpace(gotRaw)
		if got != exp {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, tc.Input(), exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
