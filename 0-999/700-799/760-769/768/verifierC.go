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

type testCase struct {
	n         int
	k         int
	x         int
	arr       []int
	expectMax int
	expectMin int
}

func expectedVals(n, k, x int, arr []int) (int, int) {
	const maxVal = 1024
	cnt := make([]int, maxVal)
	for _, v := range arr {
		cnt[v]++
	}
	for iter := 0; iter < k; iter++ {
		next := make([]int, maxVal)
		odd := 0
		for v := 0; v < maxVal; v++ {
			c := cnt[v]
			if c == 0 {
				continue
			}
			var toXor int
			if odd == 0 {
				toXor = (c + 1) / 2
			} else {
				toXor = c / 2
			}
			next[v^x] += toXor
			next[v] += c - toXor
			if c%2 == 1 {
				odd ^= 1
			}
		}
		cnt = next
	}
	minVal, maxValIdx := -1, -1
	for i := 0; i < maxVal; i++ {
		if cnt[i] > 0 {
			if minVal == -1 {
				minVal = i
			}
			maxValIdx = i
		}
	}
	return maxValIdx, minVal
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 100)
	deterministic := []testCase{
		{n: 5, k: 1, x: 2, arr: []int{5, 7, 9, 11, 15}},
		{n: 3, k: 0, x: 5, arr: []int{1, 2, 3}},
		{n: 1, k: 10, x: 0, arr: []int{42}},
	}
	for _, tc := range deterministic {
		tc.expectMax, tc.expectMin = expectedVals(tc.n, tc.k, tc.x, append([]int(nil), tc.arr...))
		tests = append(tests, tc)
	}
	for len(tests) < 100 {
		n := rng.Intn(10) + 1
		k := rng.Intn(10)
		x := rng.Intn(1024)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = rng.Intn(1024)
		}
		maxV, minV := expectedVals(n, k, x, append([]int(nil), arr...))
		tests = append(tests, testCase{n: n, k: k, x: x, arr: arr, expectMax: maxV, expectMin: minV})
	}
	return tests
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.k, tc.x))
		for j, v := range tc.arr {
			if j > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')
		expect := fmt.Sprintf("%d %d", tc.expectMax, tc.expectMin)
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input.String())
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, input.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
