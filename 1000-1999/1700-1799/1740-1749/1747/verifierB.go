package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
)

// testCase represents one input and expected minimal operations for problem B
type testCase struct {
	n int
	m int
}

func generateTests() (string, []testCase) {
	rand.Seed(1)
	t := 100
	var buf bytes.Buffer
	var cases []testCase
	fmt.Fprintln(&buf, t)
	for i := 0; i < t; i++ {
		n := rand.Intn(100) + 1
		fmt.Fprintln(&buf, n)
		m := 1
		if n > 2 {
			m = n/2 + n%2
		}
		cases = append(cases, testCase{n: n, m: m})
	}
	return buf.String(), cases
}

func hasBAN(s []byte) bool {
	b := bytes.IndexByte(s, 'B')
	if b < 0 {
		return false
	}
	a := bytes.IndexByte(s[b+1:], 'A')
	if a < 0 {
		return false
	}
	n := bytes.IndexByte(s[b+a+2:], 'N')
	return n >= 0
}

func applyOps(s []byte, ops [][2]int) []byte {
	for _, op := range ops {
		i := op[0] - 1
		j := op[1] - 1
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	input, cases := generateTests()
	cmd := exec.Command(os.Args[1])
	cmd.Stdin = bytes.NewBufferString(input)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("failed to run binary:", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for idx, tc := range cases {
		if !scanner.Scan() {
			fmt.Printf("missing output for case %d\n", idx+1)
			os.Exit(1)
		}
		var m int
		fmt.Sscan(scanner.Text(), &m)
		if m != tc.m {
			fmt.Printf("case %d: expected %d operations, got %d\n", idx+1, tc.m, m)
			os.Exit(1)
		}
		var ops [][2]int
		for k := 0; k < m; k++ {
			if !scanner.Scan() {
				fmt.Printf("missing swap at case %d operation %d\n", idx+1, k+1)
				os.Exit(1)
			}
			var i, j int
			fmt.Sscan(scanner.Text(), &i, &j)
			if i < 1 || i > 3*tc.n || j < 1 || j > 3*tc.n || i == j {
				fmt.Printf("invalid indices in case %d: %d %d\n", idx+1, i, j)
				os.Exit(1)
			}
			ops = append(ops, [2]int{i, j})
		}
		s := bytes.Repeat([]byte("BAN"), tc.n)
		s = applyOps(s, ops)
		if hasBAN(s) {
			fmt.Printf("case %d: resulting string still contains BAN subsequence\n", idx+1)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
