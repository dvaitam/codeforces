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

func compress(a []int) []int {
    if len(a) == 0 {
        return a
    }
    out := make([]int, 0, len(a))
    prev := a[0]
    out = append(out, prev)
    for i := 1; i < len(a); i++ {
        if a[i] != prev {
            out = append(out, a[i])
            prev = a[i]
        }
    }
    return out
}

func bruteMin(tc testCaseB) int {
    arr := compress(tc.arr)
    n := len(arr)
    best := 1<<31 - 1
    for mask := 0; mask < 1<<n; mask++ {
        var w, b []int
        for i := 0; i < n; i++ {
            if mask>>i&1 == 1 {
                b = append(b, arr[i])
            } else {
                w = append(w, arr[i])
            }
        }
        val := segments(w) + segments(b)
        if val < best {
            best = val
        }
    }
    return best
}

func buildInputs(tc testCaseB) (withT, withoutT []byte) {
    var a, b bytes.Buffer
    // with T header
    fmt.Fprintln(&a, 1)
    fmt.Fprintln(&a, tc.n)
    for i, v := range tc.arr {
        if i > 0 { a.WriteByte(' ') }
        fmt.Fprint(&a, v)
    }
    a.WriteByte('\n')
    // without T header
    fmt.Fprintln(&b, tc.n)
    for i, v := range tc.arr {
        if i > 0 { b.WriteByte(' ') }
        fmt.Fprint(&b, v)
    }
    b.WriteByte('\n')
    return a.Bytes(), b.Bytes()
}

func runOnce(bin string, stdin []byte) (string, string, error) {
    var cmd *exec.Cmd
    if strings.HasSuffix(bin, ".go") {
        cmd = exec.Command("go", "run", bin)
    } else {
        cmd = exec.Command(bin)
    }
    cmd.Stdin = bytes.NewReader(stdin)
    var out bytes.Buffer
    var errBuf bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &errBuf
    if err := cmd.Run(); err != nil {
        return "", errBuf.String(), err
    }
    return out.String(), errBuf.String(), nil
}

func runCase(bin string, tc testCaseB, expected int) error {
    inWithT, inNoT := buildInputs(tc)
    // try with T header first
    outStr, errStr, err := runOnce(bin, inWithT)
    if err != nil {
        // try without T
        outStr, errStr, err = runOnce(bin, inNoT)
        if err != nil {
            return fmt.Errorf("runtime error: %v\n%s", err, errStr)
        }
    }
    fields := strings.Fields(outStr)
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
