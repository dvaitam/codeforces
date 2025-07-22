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

type testCaseA struct {
	n   int
	k   int
	arr []int
}

func genTestsA() []testCaseA {
	rand.Seed(1)
	tests := make([]testCaseA, 100)
	for i := range tests {
		n := rand.Intn(8) + 2 // 2..9
		k := rand.Intn(4)
		if k > n {
			k = n
		}
		arr := make([]int, n)
		for j := range arr {
			arr[j] = rand.Intn(11) - 5
		}
		tests[i] = testCaseA{n, k, arr}
	}
	return tests
}

func solveA(tc testCaseA) int64 {
	n, k := tc.n, tc.k
	a := tc.arr
	prefix := make([]int, n+1)
	for i := 0; i < n; i++ {
		prefix[i+1] = prefix[i] + a[i]
	}
	maxSum := int64(-1 << 60)
	for l := 0; l < n; l++ {
		for r := l; r < n; r++ {
			base := prefix[r+1] - prefix[l]
			inside := make([]int, 0, r-l+1)
			outside := make([]int, 0, n-(r-l+1))
			for i := 0; i < n; i++ {
				if i >= l && i <= r {
					inside = append(inside, a[i])
				} else {
					outside = append(outside, a[i])
				}
			}
			sort.Ints(inside)
			sort.Sort(sort.Reverse(sort.IntSlice(outside)))
			cur := int64(base)
			if cur > maxSum {
				maxSum = cur
			}
			maxSwap := k
			if len(inside) < maxSwap {
				maxSwap = len(inside)
			}
			if len(outside) < maxSwap {
				maxSwap = len(outside)
			}
			for t := 0; t < maxSwap; t++ {
				if outside[t] > inside[t] {
					cur += int64(outside[t] - inside[t])
					if cur > maxSum {
						maxSum = cur
					}
				} else {
					break
				}
			}
		}
	}
	return maxSum
}

func run(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsA()
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
		expect := solveA(tc)
		out, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		val, err := strconv.ParseInt(strings.TrimSpace(out), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: non-integer output %q\n", i+1, out)
			os.Exit(1)
		}
		if val != expect {
			fmt.Fprintf(os.Stderr, "test %d: expected %d got %d\n", i+1, expect, val)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
