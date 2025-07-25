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

type testCaseC struct {
	n    int
	arr  []int64
	perm []int
}

func genTestsC() []testCaseC {
	rand.Seed(44)
	tests := make([]testCaseC, 100)
	for i := range tests {
		n := rand.Intn(8) + 2
		arr := make([]int64, n)
		for j := range arr {
			arr[j] = int64(rand.Intn(20))
		}
		p := rand.Perm(n)
		for j := range p {
			p[j]++
		}
		tests[i] = testCaseC{n, arr, p}
	}
	return tests
}

var parent []int
var segSum []int64
var used []bool

func find(x int) int {
	if parent[x] != x {
		parent[x] = find(parent[x])
	}
	return parent[x]
}

func union(x, y int) {
	rx := find(x)
	ry := find(y)
	if rx == ry {
		return
	}
	parent[ry] = rx
	segSum[rx] += segSum[ry]
}

func solveC(tc testCaseC) []int64 {
	n := tc.n
	parent = make([]int, n)
	segSum = make([]int64, n)
	used = make([]bool, n)
	ans := make([]int64, n)
	curMax := int64(0)
	for i := n - 1; i >= 0; i-- {
		ans[i] = curMax
		pos := tc.perm[i] - 1
		used[pos] = true
		parent[pos] = pos
		segSum[pos] = tc.arr[pos]
		if pos > 0 && used[pos-1] {
			union(pos, pos-1)
		}
		if pos+1 < n && used[pos+1] {
			union(pos, pos+1)
		}
		root := find(pos)
		if segSum[root] > curMax {
			curMax = segSum[root]
		}
	}
	return ans
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("runtime error: %v", err)
	}
	return out.String(), nil
}

func parseOutputC(out string, n int) ([]int64, error) {
	scanner := bufio.NewScanner(strings.NewReader(out))
	scanner.Split(bufio.ScanWords)
	res := make([]int64, n)
	for i := 0; i < n; i++ {
		if !scanner.Scan() {
			return nil, fmt.Errorf("not enough output")
		}
		v, err := strconv.ParseInt(scanner.Text(), 10, 64)
		if err != nil {
			return nil, err
		}
		res[i] = v
	}
	if scanner.Scan() {
		return nil, fmt.Errorf("extra output")
	}
	return res, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsC()
	for i, tc := range tests {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", tc.n)
		for j, v := range tc.arr {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		sb.WriteByte('\n')
		for j, v := range tc.perm {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		sb.WriteByte('\n')
		input := sb.String()
		expArr := solveC(tc)
		gotStr, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\noutput:\n%s", i+1, err, gotStr)
			os.Exit(1)
		}
		gotArr, err := parseOutputC(gotStr, tc.n)
		if err != nil {
			fmt.Printf("test %d: bad output: %v\ninput:\n%soutput:\n%s", i+1, err, input, gotStr)
			os.Exit(1)
		}
		for j := 0; j < tc.n; j++ {
			if gotArr[j] != expArr[j] {
				fmt.Printf("test %d failed on line %d: expected %d got %d\ninput:\n%s", i+1, j+1, expArr[j], gotArr[j], input)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
