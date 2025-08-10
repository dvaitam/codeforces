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
	arr := tc.arr

	// Compute basic statistics
	var total, minVal, maxVal int64
	minVal, maxVal = arr[0], arr[0]
	for _, v := range arr {
		total += v
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
	}

	maxMinPossible := total / int64(n)
	minMaxPossible := (total + int64(n) - 1) / int64(n)

	// Binary search for the highest achievable minimum wealth
	checkMin := func(target int64) bool {
		var need int64
		for _, v := range arr {
			if v < target {
				need += target - v
				if need > k {
					return false
				}
			}
		}
		return need <= k
	}

	low, high := minVal, maxMinPossible
	finalMin := minVal
	for low <= high {
		mid := low + (high-low)/2
		if checkMin(mid) {
			finalMin = mid
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	// Binary search for the lowest achievable maximum wealth
	checkMax := func(target int64) bool {
		var removed int64
		for _, v := range arr {
			if v > target {
				removed += v - target
				if removed > k {
					return false
				}
			}
		}
		return removed <= k
	}

	low, high = minMaxPossible, maxVal
	finalMax := maxVal
	for low <= high {
		mid := low + (high-low)/2
		if checkMax(mid) {
			finalMax = mid
			high = mid - 1
		} else {
			low = mid + 1
		}
	}

	if finalMax < finalMin {
		return 0
	}
	return finalMax - finalMin
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
