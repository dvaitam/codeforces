package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type TestCase struct {
	n   int
	arr []int
}

func genTests() []TestCase {
	rand.Seed(time.Now().UnixNano())
	const T = 100
	tests := make([]TestCase, T)
	for i := 0; i < T; i++ {
		n := rand.Intn(8) + 2
		arr := make([]int, n)
		for j := range arr {
			arr[j] = rand.Intn(n + 2)
		}
		tests[i] = TestCase{n: n, arr: arr}
	}
	return tests
}

func buildInput(tests []TestCase) []byte {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, len(tests))
	for _, tc := range tests {
		fmt.Fprint(&buf, tc.n)
		for _, v := range tc.arr {
			fmt.Fprint(&buf, " ", v)
		}
		fmt.Fprintln(&buf)
	}
	return buf.Bytes()
}

func runBinary(path string, input []byte) ([]byte, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = bytes.NewReader(input)
	return cmd.CombinedOutput()
}

func mex(arr []int) int {
	seen := make([]bool, len(arr)+2)
	for _, v := range arr {
		if v < len(seen) {
			seen[v] = true
		}
	}
	m := 0
	for seen[m] {
		m++
	}
	return m
}

func mexRange(arr []int, l, r int) int {
	seen := make([]bool, r-l+3)
	for i := l; i <= r; i++ {
		v := arr[i]
		if v < len(seen) {
			seen[v] = true
		}
	}
	m := 0
	for seen[m] {
		m++
	}
	return m
}

func canSplit(tc TestCase) (bool, int) {
	m := mex(tc.arr)
	if m == tc.n {
		return false, m
	}

	n := tc.n
	// try all possible partitions using a bitmask over the n-1 gaps
	for mask := 1; mask < (1 << (n - 1)); mask++ {
		l := 0
		ok := true
		for i := 0; i < n-1; i++ {
			if (mask>>i)&1 == 1 {
				if mexRange(tc.arr, l, i) != m {
					ok = false
					break
				}
				l = i + 1
			}
		}
		if !ok {
			continue
		}
		if mexRange(tc.arr, l, n-1) == m {
			return true, m
		}
	}
	return false, m
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]

	tests := genTests()
	input := buildInput(tests)

	out, err := runBinary(binary, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error executing binary: %v\n%s", err, out)
		os.Exit(1)
	}
	tokens := strings.Fields(string(out))
	idx := 0
	for i, tc := range tests {
		possible, m := canSplit(tc)
		if idx >= len(tokens) {
			fmt.Fprintf(os.Stderr, "output too short on test %d\n", i+1)
			os.Exit(1)
		}
		if !possible {
			if tokens[idx] != "-1" {
				fmt.Fprintf(os.Stderr, "mismatch on test %d\nexpected -1 got %s\n", i+1, tokens[idx])
				os.Exit(1)
			}
			idx++
			continue
		}
		k, err := strconv.Atoi(tokens[idx])
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid k on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		idx++
		if k < 2 {
			fmt.Fprintf(os.Stderr, "mismatch on test %d\nexpected k >= 2 got %d\n", i+1, k)
			os.Exit(1)
		}
		segs := make([][2]int, k)
		for j := 0; j < k; j++ {
			if idx+1 >= len(tokens) {
				fmt.Fprintf(os.Stderr, "output too short on test %d\n", i+1)
				os.Exit(1)
			}
			l, err1 := strconv.Atoi(tokens[idx])
			r, err2 := strconv.Atoi(tokens[idx+1])
			if err1 != nil || err2 != nil {
				fmt.Fprintf(os.Stderr, "invalid segment on test %d\n", i+1)
				os.Exit(1)
			}
			idx += 2
			segs[j] = [2]int{l, r}
		}
		// validate segments
		if segs[0][0] != 1 || segs[k-1][1] != tc.n {
			fmt.Fprintf(os.Stderr, "mismatch on test %d\nsegments must cover array\n", i+1)
			os.Exit(1)
		}
		for j := 0; j < k; j++ {
			l, r := segs[j][0], segs[j][1]
			if l < 1 || r > tc.n || l > r {
				fmt.Fprintf(os.Stderr, "invalid segment boundaries on test %d\n", i+1)
				os.Exit(1)
			}
			if j > 0 && l != segs[j-1][1]+1 {
				fmt.Fprintf(os.Stderr, "segments not contiguous on test %d\n", i+1)
				os.Exit(1)
			}
			if mexRange(tc.arr, l-1, r-1) != m {
				fmt.Fprintf(os.Stderr, "wrong mex on test %d\n", i+1)
				os.Exit(1)
			}
		}
	}
	if idx != len(tokens) {
		fmt.Fprintln(os.Stderr, "extra output detected")
		os.Exit(1)
	}
	fmt.Println("all tests passed")
}
