package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type testCase struct {
	input  string
	output string
}

type pair struct {
	val int
	idx int
}

func solve(arr []int) string {
	n := len(arr)
	ps := make([]pair, n)
	for i := 0; i < n; i++ {
		ps[i] = pair{arr[i], i}
	}
	sort.Slice(ps, func(i, j int) bool { return ps[i].val < ps[j].val })

	active := make([]bool, n)
	L := make([]int, n)
	counts := make([]int, n+1)

	numSegments := 0
	maxLoc := 0
	bestK := 0

	for i := 0; i < n; i++ {
		idx := ps[i].idx
		leftLen := 0
		if idx > 0 && active[idx-1] {
			leftLen = L[idx-1]
		}
		rightLen := 0
		if idx < n-1 && active[idx+1] {
			rightLen = L[idx+1]
		}

		if leftLen > 0 {
			counts[leftLen]--
			numSegments--
		}
		if rightLen > 0 {
			counts[rightLen]--
			numSegments--
		}

		newLen := leftLen + rightLen + 1
		counts[newLen]++
		numSegments++

		active[idx] = true
		L[idx-leftLen] = newLen
		L[idx+rightLen] = newLen

		if counts[newLen] == numSegments {
			if numSegments > maxLoc {
				maxLoc = numSegments
				bestK = ps[i].val + 1
			}
		}
	}

	return fmt.Sprintf("%d", bestK)
}

func generateTests() []testCase {
	rand.Seed(4)
	var tests []testCase
	tests = append(tests, testCase{input: "1\n5\n", output: "6"})
	for len(tests) < 120 {
		n := rand.Intn(10) + 1
		perm := rand.Perm(n*3 + 5)[:n]
		arr := make([]int, n)
		for i, v := range perm {
			arr[i] = v + 1
		}
		var b strings.Builder
		b.WriteString(fmt.Sprintf("%d\n", n))
		for i, v := range arr {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", v)
		}
		b.WriteString("\n")
		tests = append(tests, testCase{input: b.String(), output: solve(arr)})
	}
	return tests
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		got, err := runBinary(binary, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != tc.output {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %q got %q\n", i+1, tc.output, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
