package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testD struct {
	arr []uint64
}

func genTestsD() []testD {
	rand.Seed(122004)
	tests := make([]testD, 100)
	for i := range tests {
		n := rand.Intn(10) + 1
		arr := make([]uint64, n)
		for j := range arr {
			arr[j] = uint64(rand.Int63n(1<<30) + 1)
		}
		tests[i] = testD{arr: arr}
	}
	return tests
}

func solveD(tc testD) int {
	cnt := make([]int, 64)
	for _, v := range tc.arr {
		tz := bits.TrailingZeros64(v)
		if tz < len(cnt) {
			cnt[tz]++
		}
	}
	maxCnt := 0
	for _, c := range cnt {
		if c > maxCnt {
			maxCnt = c
		}
	}
	return len(tc.arr) - maxCnt
}

func runCandidate(bin string, input string) (string, error) {
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
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsD()

	for i, tc := range tests {
		var input bytes.Buffer
		fmt.Fprintln(&input, len(tc.arr))
		for j, v := range tc.arr {
			if j > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, v)
		}
		input.WriteByte('\n')

		expectedK := solveD(tc)

		outStr, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}

		scanner := bufio.NewScanner(strings.NewReader(outStr))
		scanner.Split(bufio.ScanWords)
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "wrong output format on test %d\n", i+1)
			os.Exit(1)
		}
		k, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Fprintf(os.Stderr, "non-integer output on test %d\n", i+1)
			os.Exit(1)
		}
		if k != expectedK {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected length %d got %d\n", i+1, expectedK, k)
			os.Exit(1)
		}

		erased := make(map[uint64]int)
		for j := 0; j < k; j++ {
			if !scanner.Scan() {
				fmt.Fprintf(os.Stderr, "wrong output format on test %d (not enough erased elements)\n", i+1)
				os.Exit(1)
			}
			val, err := strconv.ParseUint(scanner.Text(), 10, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "non-integer output on test %d\n", i+1)
				os.Exit(1)
			}
			erased[val]++
		}

		if scanner.Scan() {
			fmt.Fprintln(os.Stderr, "extra output")
			os.Exit(1)
		}

		// Validate that erased elements are a subset of original array
		originalCounts := make(map[uint64]int)
		for _, v := range tc.arr {
			originalCounts[v]++
		}

		for v, count := range erased {
			if originalCounts[v] < count {
				fmt.Fprintf(os.Stderr, "erased element %d not in original array or erased too many times\n", v)
				os.Exit(1)
			}
			originalCounts[v] -= count
		}

		// and the remaining have same trailing zeros
		remainingTz := -1
		for v, count := range originalCounts {
			for c := 0; c < count; c++ {
				tz := bits.TrailingZeros64(v)
				if remainingTz == -1 {
					remainingTz = tz
				} else if remainingTz != tz {
					fmt.Fprintf(os.Stderr, "wrong answer on test %d: remaining elements have different trailing zeros\n", i+1)
					os.Exit(1)
				}
			}
		}
	}
	fmt.Println("Accepted")
}