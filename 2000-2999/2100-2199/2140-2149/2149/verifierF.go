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

type testCase struct {
	input  string
	expect string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for idx, tc := range tests {
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\ninput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		if err := check(tc.expect, strings.TrimSpace(out)); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%s\nexpected:\n%s\nactual:\n%s\n", idx+1, err, tc.input, tc.expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

func check(expect, actual string) error {
	exp, _ := strconv.ParseInt(expect, 10, 64)
	val, err := strconv.ParseInt(actual, 10, 64)
	if err != nil {
		return fmt.Errorf("output is not integer: %v", err)
	}
	if val != exp {
		return fmt.Errorf("expected %d but got %d", exp, val)
	}
	return nil
}

func genTests() []testCase {
	rand.Seed(42)
	tests := []testCase{
		makeCase(3, 2),
		makeCase(1, 1),
		makeCase(5, 3),
		makeCase(2, 4),
		makeCase(10, 7),
	}
	for i := 0; i < 200; i++ {
		h := rand.Int63n(1000) + 1
		d := rand.Int63n(1000) + 1
		tests = append(tests, makeCase(h, d))
	}
	return tests
}

func makeCase(h, d int64) testCase {
	return testCase{
		input:  fmt.Sprintf("1\n%d %d\n", h, d),
		expect: fmt.Sprintf("%d", solveReference(h, d)),
	}
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func solveReference(h, d int64) int64 {
	lo, hi := int64(0), d
	for lo < hi {
		mid := (lo + hi) / 2
		k := mid + 1
		if k > d {
			k = d
		}
		if minTriSum(d, k) <= h+mid-1 {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	return d + lo
}

func minTriSum(d, k int64) int64 {
	q := d / k
	r := d % k
	triQ := triangle(q)
	triQ1 := triangle(q + 1)
	return (k-r)*triQ + r*triQ1
}

func triangle(n int64) int64 {
	return n * (n + 1) / 2
}
