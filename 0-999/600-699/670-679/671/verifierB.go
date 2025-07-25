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

type testCaseB struct {
	n   int
	k   int64
	arr []int64
}

func genTestsB() []testCaseB {
	rand.Seed(2)
	tests := make([]testCaseB, 100)
	for i := range tests {
		n := rand.Intn(6) + 1 //1..6
		k := int64(rand.Intn(20))
		arr := make([]int64, n)
		for j := range arr {
			arr[j] = int64(rand.Intn(20)) + 1
		}
		tests[i] = testCaseB{n, k, arr}
	}
	return tests
}

func solveB(tc testCaseB) int64 {
	n, k := tc.n, tc.k
	arr := append([]int64(nil), tc.arr...)
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	l, r := 0, n-1
	for l < r && k > 0 {
		if l+1 <= n-1-r {
			diff := arr[l+1] - arr[l]
			cost := diff * int64(l+1)
			if cost <= k {
				k -= cost
				l++
			} else {
				arr[l] += k / int64(l+1)
				k = 0
			}
		} else {
			diff := arr[r] - arr[r-1]
			cost := diff * int64(n-r)
			if cost <= k {
				k -= cost
				r--
			} else {
				arr[r] -= k / int64(n-r)
				k = 0
			}
		}
	}
	if arr[r] < arr[l] {
		return 0
	}
	return arr[r] - arr[l]
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsB()
	for i, tc := range tests {
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d %d\n", tc.n, tc.k)
		for j, v := range tc.arr {
			if j > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, v)
		}
		input.WriteByte('\n')
		expect := solveB(tc)
		out, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		val, err := strconv.ParseInt(strings.TrimSpace(out), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: invalid output %q\n", i+1, out)
			os.Exit(1)
		}
		if val != expect {
			fmt.Fprintf(os.Stderr, "test %d: expected %d got %d\n", i+1, expect, val)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
