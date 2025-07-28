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

type testCaseB struct {
	n   int
	arr []int
}

func genTestsB() []testCaseB {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCaseB, 100)
	for i := range cases {
		n := rng.Intn(6) + 1 // 1..6
		arr := make([]int, n)
		for j := range arr {
			arr[j] = rng.Intn(n) + 1
		}
		cases[i] = testCaseB{n: n, arr: arr}
	}
	return cases
}

func segments(a []int) int {
	if len(a) == 0 {
		return 0
	}
	cnt := 1
	for i := 1; i < len(a); i++ {
		if a[i] != a[i-1] {
			cnt++
		}
	}
	return cnt
}

func bruteMin(tc testCaseB) int {
	n := tc.n
	best := 1<<31 - 1
	for mask := 0; mask < 1<<n; mask++ {
		var w, b []int
		for i := 0; i < n; i++ {
			if mask>>i&1 == 1 {
				b = append(b, tc.arr[i])
			} else {
				w = append(w, tc.arr[i])
			}
		}
		val := segments(w) + segments(b)
		if val < best {
			best = val
		}
	}
	return best
}

func runCase(bin string, tc testCaseB, expected int) error {
	var input bytes.Buffer
	fmt.Fprintln(&input, 1)
	fmt.Fprintln(&input, tc.n)
	for i, v := range tc.arr {
		if i > 0 {
			input.WriteByte(' ')
		}
		fmt.Fprint(&input, v)
	}
	input.WriteByte('\n')

	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = bytes.NewReader(input.Bytes())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(out.String())
	if len(fields) == 0 {
		return fmt.Errorf("no output")
	}
	val, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("non-integer output")
	}
	if val != expected {
		return fmt.Errorf("expected %d got %d", expected, val)
	}
	if len(fields) > 1 {
		return fmt.Errorf("extra output")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genTestsB()
	for i, tc := range cases {
		exp := bruteMin(tc)
		if err := runCase(bin, tc, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
