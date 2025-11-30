package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases, one per line in the same flattened format the solver expects.
const testcasesE = `1 3 1 2
2 2 1 1
1 4 1 1 1
3 5 1 2 1 1 3 1 1 3 1 1
1 3 1 2
1 5 1 2 3 2
1 2 1
2 1 5 1 1 2 3
3 2 1 4 1 2 2 2 1
2 3 1 1 1
2 3 1 2 2 1
1 4 1 2 2
1 3 1 2
3 5 1 2 1 1 2 1 5 1 1 2 2
2 4 1 1 2 1
2 3 1 2 3 1 1
3 4 1 2 3 2 1 3 1 2
2 5 1 2 3 1 3 1 2
1 3 1 2
3 1 4 1 2 1 1
1 1
1 2 1
3 1 5 1 1 2 1 1
3 3 1 1 4 1 1 2 4 1 1 1
2 4 1 2 1 4 1 1 1
2 3 1 1 2 1
2 2 1 1
2 5 1 2 3 1 4 1 1 2
1 5 1 2 2 4
3 3 1 1 3 1 2 4 1 2 3
3 2 1 4 1 2 3 1
1 4 1 2 2
1 1
1 2 1
3 3 1 2 4 1 1 1 2 1
3 4 1 1 1 1 5 1 2 3 3
3 2 1 3 1 1 3 1 1
1 5 1 2 1 3
2 2 1 4 1 1 1
3 3 1 2 4 1 1 3 3 1 2
1 5 1 2 1 2
3 2 1 4 1 2 3 1
3 2 1 3 1 1 4 1 1 3
1 4 1 2 1
1 1
3 1 5 1 2 3 3 5 1 2 2 2
2 5 1 2 1 4 1
3 2 1 3 1 1 5 1 1 1 4
3 2 1 3 1 1 3 1 1
3 1 3 1 1 2 1
3 5 1 1 1 4 2 1 4 1 2 3
1 3 1 1
3 3 1 2 5 1 2 2 4 4 1 2 1
3 4 1 2 2 5 1 1 3 1 4 1 2 1
2 4 1 2 3 1
1 2 1
2 2 1 2 1
3 5 1 1 3 2 2 1 2 1
2 1 1
2 1 3 1 1
3 4 1 1 3 3 1 1 5 1 1 1 3
2 5 1 2 1 2 3 1 2
2 3 1 2 1
2 1 3 1 2
1 1
1 1
3 1 2 1 5 1 2 3 2
2 4 1 1 3 5 1 1 2 3
2 4 1 1 2 5 1 2 3 3
1 2 1
1 1
3 3 1 2 1 2 1
2 5 1 1 1 1 3 1 2
1 5 1 1 1 3
3 1 1 1
2 2 1 3 1 2
2 4 1 1 1 4 1 1 3
3 3 1 1 4 1 1 1 1
2 4 1 2 2 4 1 1 2
2 5 1 1 2 3 5 1 1 3 3
1 2 1
3 2 1 2 1 1
3 3 1 2 5 1 1 1 2 5 1 1 2 3
3 5 1 1 3 2 1 2 1
2 1 5 1 2 3 3
1 2 1
1 1
1 1
3 3 1 1 2 1 4 1 2 3
2 3 1 2 1
1 2 1
2 2 1 2 1
3 2 1 4 1 1 1 2 1
1 5 1 1 1 1
1 3 1 2
1 3 1 1
1 4 1 2 2
1 4 1 2 1
3 4 1 1 2 3 1 1 4 1 1 1
3 2 1 5 1 1 2 4 2 1`

type IntMaxHeap []int

func (h IntMaxHeap) Len() int           { return len(h) }
func (h IntMaxHeap) Less(i, j int) bool { return h[i] > h[j] }
func (h IntMaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntMaxHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func (h *IntMaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func parseTestcase(line string) ([]int, error) {
	fields := strings.Fields(line)
	nums := make([]int, len(fields))
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q in testcase: %v", f, err)
		}
		nums[i] = v
	}
	return nums, nil
}

// solveCase reproduces the logic from 1994E.go for a single testcase.
func solveCase(data []int) (int, error) {
	if len(data) == 0 {
		return 0, fmt.Errorf("empty testcase")
	}
	idx := 0
	k := data[idx]
	idx++
	if k < 0 {
		return 0, fmt.Errorf("invalid k: %d", k)
	}

	caps := make([]int, 0, k)
	maxN := 0
	for i := 0; i < k; i++ {
		if idx >= len(data) {
			return 0, fmt.Errorf("testcase truncated reading n for tree %d", i+1)
		}
		n := data[idx]
		idx++
		if n > maxN {
			maxN = n
		}
		caps = append(caps, n)
		need := n - 1
		if idx+need > len(data) {
			return 0, fmt.Errorf("testcase truncated reading parents for tree %d", i+1)
		}
		idx += need
	}
	if idx != len(data) {
		return 0, fmt.Errorf("extra data at the end of testcase")
	}

	h := IntMaxHeap(caps)
	heap.Init(&h)
	res := 0
	for bit := 20; bit >= 0; bit-- {
		val := 1 << bit
		if val > maxN {
			continue
		}
		if h.Len() == 0 {
			break
		}
		if h[0] >= val {
			topCap := heap.Pop(&h).(int)
			topCap -= val
			if topCap > 0 {
				heap.Push(&h, topCap)
			}
			res |= val
		}
	}
	return res, nil
}

func runBinary(bin, input string) (string, string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	lines := strings.Split(testcasesE, "\n")
	passed := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		passed++
		nums, err := parseTestcase(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d parse error: %v\n", passed, err)
			os.Exit(1)
		}
		expectedVal, err := solveCase(nums)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d invalid: %v\n", passed, err)
			os.Exit(1)
		}
		input := "1\n" + line + "\n"
		out, stderr, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", passed, err, stderr)
			os.Exit(1)
		}
		got := strings.TrimSpace(out)
		expected := strconv.Itoa(expectedVal)
		if got != expected {
			fmt.Printf("test %d failed. Expected %s got %s\n", passed, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", passed)
}
