package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type testCase struct {
	input  string
	expect string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for idx, tc := range tests {
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\ninput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		if err := check(tc.expect, strings.TrimSpace(out)); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%s\nexpected:\n%s\nactual:\n%s\n", idx+1, err, tc.input, tc.expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

func check(expect, actual string) error {
	exp, _ := strconv.Atoi(expect)
	val, err := strconv.Atoi(actual)
	if err != nil {
		return fmt.Errorf("output is not integer: %v", err)
	}
	if val != exp {
		return fmt.Errorf("expected %d but got %d", exp, val)
	}
	return nil
}

func genTests() []testCase {
	rand.Seed(42)
	tests := []testCase{
		makeCase(3, 3, 120, [][]int{{20, 15, 110}, {90, 30, 80}, {100, 50, 90}}),
		makeCase(2, 2, 30, [][]int{{15, 15}, {15, 15}}),
		makeCase(1, 3, 60, [][]int{{20, 30, 40}}),
	}
	for i := 0; i < 200; i++ {
		n := rand.Intn(4) + 1
		m := rand.Intn(4) + 1
		h := rand.Intn(100) + 10
		participants := make([][]int, n)
		for j := 0; j < n; j++ {
			participants[j] = make([]int, m)
			for k := 0; k < m; k++ {
				participants[j][k] = rand.Intn(50) + 1
			}
		}
		tests = append(tests, makeCase(n, m, h, participants))
	}
	return tests
}

func makeCase(n, m, h int, participants [][]int) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d %d\n", n, m, h)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", participants[i][j])
		}
		sb.WriteByte('\n')
	}
	return testCase{
		input:  sb.String(),
		expect: fmt.Sprintf("%d", solveReference(n, m, h, participants)),
	}
}

func solveReference(n, m, h int, participants [][]int) int {
	type pair struct {
		solved int
		pen    int64
	}
	score := make([]pair, n)
	for i := 0; i < n; i++ {
		times := append([]int(nil), participants[i]...)
		sort.Ints(times)
		sum := int64(0)
		solved := 0
		for _, t := range times {
			if sum+int64(t) > int64(h) {
				break
			}
			sum += int64(t)
			solved++
		}
		penalty := int64(0)
		sum = 0
		for j := 0; j < solved; j++ {
			sum += int64(times[j])
			penalty += sum
		}
		score[i] = pair{solved, penalty}
	}
	rScore := score[0]
	rank := 1
	for i := 1; i < n; i++ {
		if score[i].solved > rScore.solved || (score[i].solved == rScore.solved && score[i].pen < rScore.pen) {
			rank++
		}
	}
	return rank
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}
